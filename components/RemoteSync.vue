<template>
  <v-container v-if="rs !== null">
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
              <tr>
                <td class="text-center">Name</td>
                <td class="text-center"><div style="white-space: nowrap">{{rs.metadata.name}}</div></td>
              </tr>
              <tr>
                <td class="text-center">Image</td>
                <td class="text-center"><div style="white-space: nowrap">{{rs.spec.image.name}}</div></td>
              </tr>
              <tr>
                <td class="text-center">Base</td>
                <td class="text-center"><a @click="openApp(rs.spec)"> {{calcBase(rs.spec.base)}}</a></td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { fetch } from "~/util/proxy";
import { calcBase, openApp } from "~/util/remoteSyncUtil";
import { AppsV1Api, Configuration, ConfigurationParameters, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1AppPipeline, DevTsuzuModokiV1alpha1AppPipelineSpecBase, V1Pod, DevTsuzuModokiV1alpha1RemoteSync } from "@modoki-paas/kubernetes-fetch-client";

export default Vue.extend({
  props: {
    rs: Object,
  },
  methods: {
    calcBase(spec: any) {
      return calcBase(spec)
    },
    openApp(spec: any) {
      openApp(spec)
    }
  }
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
