<template>
  <v-container>
    <v-row dense>
      <v-col cols="3">
        <span class="text-h4">Pods</span>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-data-table
          item-key="metadata.name"
          :headers="headers"
          :items="calcedPods"
          mobile-breakpoint="0"
        >
          <!-- <template v-slot:[`item.actions`]="{ item }">
            <v-icon
              class="mr-2"
              @click.stop="openApp(item.spec.domains[0])"
            >
              mdi-open-in-new
            </v-icon>
          </template> -->
        </v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'
import { fetch } from "~/util/proxy";
import { AppsV1Api, Configuration, ConfigurationParameters, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1AppPipeline, DevTsuzuModokiV1alpha1AppPipelineSpecBase, V1Pod } from "@modoki-paas/kubernetes-fetch-client";

export default Vue.extend({
  components: {
    Logo,
    VuetifyLogo,
  },
  props: {
    pods: Array,
  },
  data() {
    return {
      headers: [
        {
          text: "Name",
          value: "metadata.name",
        },
        {
          text: "Phase",
          value: "status.phase",
        },
        {
          text: "Message",
          value: "status.message",
        },
        {
          text: "Image",
          value: "spec.containers[0].image",
          align: 'left'
        },
      ]
    }
  },
  computed: {
    calcedPods() {
      return ((this as any).pods as V1Pod[]).map(pod => ({
          ...pod,
          ready: `${pod.status?.containerStatuses?.filter(c => c.ready).length ?? 0}/${pod.status?.containerStatuses?.length ?? 0}`,
        }))
    }
  },
})
</script>

<style>
.text-end {
  text-align: end;
}
.v-data-table {
  white-space : nowrap;
}
</style>
