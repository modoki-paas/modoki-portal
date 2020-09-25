<template>
  <v-tabs-items v-model="tab">
    <v-tab-item value="info">
      <v-container v-if="app !== null">
        <v-row>
          <v-col cols="12">
            <span class="text-h4">Spec</span>
          </v-col>
          <v-col cols="12">
            <v-simple-table>
              <template v-slot:default width="100%">
                <thead>
                  <tr>
                    <th class="text-center">Key</th>
                    <th class="text-center">Value</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(domain, i) in app.spec.domains" :key="domain">
                    <td class="text-center">{{i === 0 ? "Domains" : ""}}</td>
                    <td class="text-center"><div style="white-space: nowrap">{{domain}}</div></td>
                  </tr>
                  <tr>
                    <td class="text-center">Image</td>
                    <td class="text-center"><div style="white-space: nowrap">{{app.spec.image}}</div></td>
                  </tr>
                  <tr>
                    <td class="text-center">Commands</td>
                    <td class="text-center">{{app.spec.commands}}</td>
                  </tr>
                  <tr>
                    <td class="text-center">Args</td>
                    <td class="text-center">{{app.spec.args}}</td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-col>
          <v-col cols="12">
            <span class="text-h6">Attributes</span>
          </v-col>
          <v-col cols="12">
            <v-simple-table>
              <template v-slot:default width="100%">
                <thead>
                  <tr>
                    <th class="text-center">Key</th>
                    <th class="text-center">Value</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="entry in attributes" :key="entry.key">
                    <td class="text-center">{{entry.key}}</td>
                    <td class="text-center"><div style="white-space: nowrap">{{entry.value}}</div></td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12">
            <span class="text-h4">Status</span>
          </v-col>
          <v-col cols="12">
            <v-simple-table>
              <template v-slot:default width="100%">
                <thead>
                  <tr>
                    <th class="text-center">Key</th>
                    <th class="text-center">Value</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td class="text-center">Status</td>
                    <td class="text-center"><div style="white-space: nowrap">{{app.status.status}}</div></td>
                  </tr>
                  <tr v-for="(domain, i) in app.status.domains" :key="domain">
                    <td class="text-center">{{i === 0 ? "Available Domains" : ""}}</td>
                    <td class="text-center"><div style="white-space: nowrap">{{domain}}</div></td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-col>
          <v-col cols="12">
            <span class="text-h6">Managed Resources</span>
          </v-col>
          <v-col cols="12">
            <v-simple-table>
              <template v-slot:default width="100%">
                <thead>
                  <tr>
                    <th class="text-center">Key</th>
                    <th class="text-center">Value</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="entry in app.status.resources" :key="`${entry.apiVersion}/${entry.kind}/${entry.name}/${entry.namespace}`">
                    <td class="text-center">{{entry.kind}}({{entry.apiVersion}})</td>
                    <td class="text-center"><div style="white-space: nowrap">{{entry.name}}{{entry.namespace !== app.metadata.namespace ? entry.namespace : ""}}</div></td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-col>
        </v-row>
      </v-container>
    </v-tab-item>
    <v-tab-item value="remote">
      <v-container v-if="app !== null">
        <v-row>
          <v-col cols="12">
            <span class="text-h4">Spec</span>
          </v-col>
        </v-row>
      </v-container>
    </v-tab-item>
  </v-tabs-items>
</template>

<script lang="ts">
import Vue from 'vue'
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'
import { fetch } from "~/util/proxy";
import { CoreV1Api, AppsV1Api, Configuration, ConfigurationParameters, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1Application } from "@modoki-paas/kubernetes-fetch-client";

const namespace = "default";

export default Vue.extend({
  components: {
    Logo,
    VuetifyLogo,
  },
  props: {

  },
  data() {
    return {
      tab: "info",
      app: null,
      attributes: {},
    } as {
      tab: string;
      app: DevTsuzuModokiV1alpha1Application | null,
      attributes: {key: string; value: string}[];
    }
  },
  beforeUpdate() {
    console.log(this.$route.hash);
  },
  async mounted() {
    if(this.$route.hash.length === 0) {
      this.$router.push(`${this.$route.path}#info`)
    }

    this.$nuxt.$emit(
      "headerInfo", {
        title: this.$route.params.name ?? "Loading...",
        prev: ()=> {
          this.$router.push("/")
        },
        path: this.$route.path,
        tabs: [{
          name: "Info",
          key: "info",
        }, {
          name: "Remote",
          key: "remote",
        }, {
          name: "Log",
          key: "log",
        }],
      }
    )
    console.log(this.$route.path);

    this.$nuxt.$on(
      "selectTab", this.selectTab
    );

    const conf = new Configuration({
      fetchApi: fetch,
    })

    const modokiApi = new ModokiTsuzuDevV1alpha1Api(conf);
    const appsApi = new AppsV1Api(conf);
    const coreApi = new CoreV1Api(conf);

    this.app = await modokiApi.readNamespacedApplication({
      name: this.$route.params.name,
      namespace: namespace,
    })
    if(this.app.spec) {
      this.attributes = Object.entries(this.app.spec.attributes ?? {}).map((arr: string[]) => ({
        key: arr[0],
        value: arr[1],
      })) ?? [];
      console.log(this.attributes);

      if(this.app.status) {
        const deployments = this.app.status.resources.filter(x => x.apiVersion === "apps/v1" && x.kind === "Deployment");

        if (deployments.length) {
          const dply = await appsApi.readNamespacedDeployment({
            name: deployments[0].name,
            namespace: deployments[0].namespace ?? this.app.metadata?.namespace ?? "",
          })



          // dply.spec.selector.matchLabels
        }
      }
    }
  },
  methods: {
    openApp(domain: string) {
      window.open(domain, '_blank');
    },
    selectTab(name: string) {
      // this.$data.tab = name
    }
  },
  watch: {
    $route () {
      if(this.$route.hash.length === 0) {
        this.$router.push(`${this.$route.path}#info`)
      }

      this.tab = this.$route.hash.substring(1)
    }
  }
})
</script>

<style>
.text-end {
  text-align: end;
}
</style>
