<template>
  <v-container>
    <v-row dense>
      <v-col cols="3">
        <span class="text-h4">AppPipeline</span>
      </v-col>
      <v-spacer></v-spacer>
      <v-col style="text-align: right" cols="3">
        <v-btn large @click="dialog=true">New AppPipeline</v-btn>
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
          :items="calcedAppPipelines"
          mobile-breakpoint="0"
          @click:row="click"
        >
          <template v-slot:[`item.actions`]="{ item }">
            <v-icon
              class="mr-2"
              @click.stop="openApp(item.spec.domains[0])"
            >
              mdi-open-in-new
            </v-icon>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <app-pipeline-form @close="close" :dialog="dialog"></app-pipeline-form>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'
import { fetch } from "~/util/proxy";
import AppPipelineForm from "~/components/AppPipelineForm.vue";
import { AppsV1Api, Configuration, ConfigurationParameters, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1AppPipeline, DevTsuzuModokiV1alpha1AppPipelineSpecBase } from "@modoki-paas/kubernetes-fetch-client";

export default Vue.extend({
  components: {
    Logo,
    VuetifyLogo,
    AppPipelineForm,
  },
  props: {},
  data() {
    return {
      modokiApi: undefined as (ModokiTsuzuDevV1alpha1Api | undefined),
      dialog: false,
      appPipelines: [] as DevTsuzuModokiV1alpha1AppPipeline[],
      headers: [
        {
          text: "Name",
          value: "metadata.name",
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
      ]
    }
  },
  async created(): Promise<void> {
    const conf = new Configuration({
      fetchApi: fetch,
    })

    const modokiApi = new ModokiTsuzuDevV1alpha1Api(conf);

    this.modokiApi = modokiApi;

    await this.reload();
  },
  mounted(): void {
    this.$nuxt.$emit(
      "headerInfo", {
        title: "modoki portal",
        prev: undefined,
        tabs: [],
      }
    )
  },
  computed: {
    calcedAppPipelines(): (DevTsuzuModokiV1alpha1AppPipeline & {baseString: string})[] {
      return (this.appPipelines as DevTsuzuModokiV1alpha1AppPipeline[]).map(ap => ({
          ...ap,
          baseString: this.calcBase(ap.spec?.base),
        }))
    }
  },
  methods: {
    calcBase(b?: DevTsuzuModokiV1alpha1AppPipelineSpecBase): string {
      return `${b?.github.owner}/${b?.github.repo}`
    },
    click(item: any): void {
      console.log(item);
      this.$router.push(`/apppipeline/${item.metadata.name}`)
    },
    openApp(domain: string): void {
      window.open("http://" + domain, '_blank');
    },
    async reload(): Promise<void> {
      if(this.modokiApi)
        this.appPipelines = (await this.modokiApi.listAppPipelineForAllNamespaces({})).items;
    },
    async close(ap : DevTsuzuModokiV1alpha1AppPipeline | undefined): Promise<void> {
      this.dialog = false;
      if(ap && this.modokiApi) {
        console.log(ap)

        await this.modokiApi.createNamespacedAppPipeline({
          body: ap,
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
