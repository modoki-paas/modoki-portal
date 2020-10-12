<template>
  <v-container>
    <v-row dense>
      <v-col cols="3">
        <span class="text-h4">RemoteSync</span>
      </v-col>
      <v-spacer></v-spacer>
      <v-col style="text-align: right" cols="3">
        <v-btn large @click="dialog=true">New RemoteSync</v-btn>
      </v-col>
    </v-row>
    <!-- <v-row>
      <v-col cols="12">
        <v-simple-table>
          <template v-slot:default>
            <thead>
              <tr>
                <th class="text-left">Name</th>
                <th class="text-right">Domain</th>
                <th class="text-left">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="app in apps" :key="app.metadata.name">
                <td> {{app.metadata.name}}</td>
                <td class="text-right" v-if="app.status.domains.length === 0"> No domain</td>
                <td class="text-right" v-else-if="app.status.domains.length === 1"> {{app.status.domains[0]}}</td>
                <td class="text-right" v-else>{{app.status.domains.join(", ")}}</td>
                <td> {{app.status.status}}</td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-col>
    </v-row> -->
    <v-row>
      <v-col cols="12">
        <v-data-table
          item-key="spec.name"
          :headers="headers"
          :items="calcedRemoteSyncs"
          mobile-breakpoint="0"
          @click:row="click"
        >
          <template v-slot:[`item.actions`]="{ item }">
            <v-icon
              class="mr-2"
              @click.stop="openApp(item.spec)"
            >
              mdi-open-in-new
            </v-icon>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <remote-sync-form @close="close" :dialog="dialog"></remote-sync-form>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'
import { fetch } from "~/util/proxy";
import RemoteSyncForm from "~/components/RemoteSyncForm.vue";
import { calcBase, openApp } from "~/util/remoteSyncUtil";
import { AppsV1Api, Configuration, ConfigurationParameters, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1RemoteSync, DevTsuzuModokiV1alpha1RemoteSyncSpecBase, DevTsuzuModokiV1alpha1RemoteSyncSpec } from "@modoki-paas/kubernetes-fetch-client"

export default Vue.extend({
  components: {
    Logo,
    VuetifyLogo,
    RemoteSyncForm,
  },
  data() {
    return {
      modokiApi: undefined as (ModokiTsuzuDevV1alpha1Api | undefined),
      dialog: false as boolean,
      remoteSyncs: [] as DevTsuzuModokiV1alpha1RemoteSync[],
      headers: [
        {
          text: "Name",
          value: "metadata.name",
        },
        {
          text: "Application Name",
          value: "spec.applicationRef.name",
          align: 'right'
        },
        {
          text: "Base",
          value: "baseString",
          align: 'right'
        },
        {
          text: "Actions",
          value: "actions",
          align: 'right'
        },
      ] as {text: string; value: string, align?: string}[]
    }
  },
  async created() {
    const conf = new Configuration({
      fetchApi: fetch,
    })

    const modokiApi = new ModokiTsuzuDevV1alpha1Api(conf);

    this.modokiApi = modokiApi;

    await this.reload();
    // this.apps[0].metadata.name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    // this.apps[0].status?.domains.push("foo.modoki.misw.jp");
    // this.apps[0].status?.domains.push("bar.modoki.misw.jp");
    // this.apps.push(this.apps[0]);
    // this.apps.push(this.apps[0]);
    // this.apps.push(this.apps[0]);
    // this.apps.push(this.apps[0]);

    // const appsClient = new AppsV1Api(conf);
    // (await appsClient.listNamespacedDeployment({
    //   namespace: "modoki-operator-system",
    // })).items.forEach(dpl => {
    //   console.log(JSON.stringify(dpl))
    // })
  },
  mounted() {
    this.$nuxt.$emit(
      "headerInfo", {
        title: "modoki portal",
        prev: undefined,
        tabs: [],
      }
    )
  },
  computed: {
    calcedRemoteSyncs(): (DevTsuzuModokiV1alpha1RemoteSync & {baseString: string})[] {
      return this.remoteSyncs.map((rs: DevTsuzuModokiV1alpha1RemoteSync) => ({
          ...rs,
          baseString: calcBase(rs.spec?.base),
        }))
    }
  },
  methods: {
    click(item: any) {
      console.log(item);
      this.$router.push(`/app/${item.spec.applicationRef.name}`)
    },
    openApp(spec: DevTsuzuModokiV1alpha1RemoteSyncSpec) {
      openApp(spec)
    },
    async reload() {
      if(this.modokiApi)
        this.remoteSyncs = (await this.modokiApi.listRemoteSyncForAllNamespaces({})).items;
    },
    async close(rs : DevTsuzuModokiV1alpha1RemoteSync | undefined) {
      this.dialog = false;
      if(rs && this.modokiApi) {
        console.log(rs)

        await this.modokiApi.createNamespacedRemoteSync({
          body: rs,
          namespace: this.$route.params.namespace,
        })

        await this.reload();
      }
    },
  }
})
</script>

<style>
.text-end {
  text-align: end;
}
</style>
