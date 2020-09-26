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
    <v-tab-item value="pods">
      <pods :pods="pods"></pods>
    </v-tab-item>
    <v-tab-item value="remote">
      <remote-sync :rs="remoteSync"></remote-sync>
    </v-tab-item>
  </v-tabs-items>
</template>

<script lang="ts">
import Vue from 'vue'
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'
import Pods from '~/components/Pods.vue'
import RemoteSync from '~/components/RemoteSync.vue'
import { fetch } from "~/util/proxy";
import { CoreV1Api, AppsV1Api, Configuration, ConfigurationParameters, V1Pod, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1Application, DevTsuzuModokiV1alpha1RemoteSync } from "@modoki-paas/kubernetes-fetch-client";

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
      app: null as DevTsuzuModokiV1alpha1Application | null,
      attributes: [] as {key: string; value: string}[],
      pods: [] as V1Pod[],
      basePath: `/app/${this.$route.params.name}`,
      modokiApi:       null as (ModokiTsuzuDevV1alpha1Api | null),
      appsApi: null as (AppsV1Api | null),
      coreApi: null as (CoreV1Api | null),
      remoteSync: null as (DevTsuzuModokiV1alpha1RemoteSync | null),
    }
  },
  beforeUpdate() {
    console.log(this.$route.hash);
  },
  async mounted() {
    const tab = this.$route.hash;
    this.tab = tab.length ? tab.substr(1) : "info";
    console.log(this.$route.params)

    this.$nuxt.$emit(
      "headerInfo", {
        title: this.$route.params.name ?? "Loading...",
        prev: ()=> {
          this.$router.back()
        },
        path: this.$route.path,
        tabs: [{
          name: "Info",
          path: `${this.basePath}`,
          hash: ""
        }, {
          name: "Pods",
          path: `${this.basePath}`,
          hash: "#pods"
        }, {
          name: "Remote",
          path: `${this.basePath}`,
          hash: "#remote"
        }, {
          name: "Log",
          path: `${this.basePath}`,
          hash: "#log"
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

    this.modokiApi = new ModokiTsuzuDevV1alpha1Api(conf);
    this.appsApi = new AppsV1Api(conf);
    this.coreApi = new CoreV1Api(conf);

    this.app = await this.modokiApi.readNamespacedApplication({
      name: this.$route.params.name,
      namespace: namespace,
    })

    await this.listPods()
    await this.getRemoteSync()
  },
  methods: {
    openApp(domain: string) {
      window.open(domain, '_blank');
    },
    selectTab(name: string) {
      // this.$data.tab = name
    },
    async listPods() {
      if(this.app?.spec) {
        this.pods = (await this.coreApi!.listNamespacedPod({
          labelSelector: `modoki-app=${this.app.metadata?.name}`,
          namespace: this.app.metadata?.namespace ?? namespace,
        })).items
      }
    },
    async getRemoteSync() {
      const rss = (await this.modokiApi!.listNamespacedRemoteSync({
        namespace,
      })).items.filter(rs => rs.spec?.applicationRef.name === this.app!.metadata!.name);
      if(rss.length) {
        this.remoteSync = rss[0];
      }
    }
  },
  watch: {
    $route () {
      const tab = this.$route.hash;
      this.tab = tab.length ? tab.substr(1) : "info";
      console.log(tab)
    }
  }
})
</script>

<style>
.text-end {
  text-align: end;
}
</style>
