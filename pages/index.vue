<template>
  <v-container>
    <v-row dense>
      <v-col cols="12">
        <p class="text-h4">Apps</p>
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
          :items="apps"
          mobile-breakpoint="0"
        >
          <template v-slot:[`item.actions`]="{ item }">
            <v-icon
              class="mr-2"
            >
              mdi-open-in-new
            </v-icon>
            <v-icon
              class="mr-2"
            >
              mdi-dots-horizontal
            </v-icon>
          </template>
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
import { AppsV1Api, Configuration, ConfigurationParameters, ModokiTsuzuDevV1alpha1Api, DevTsuzuModokiV1alpha1Application } from "@modoki-paas/kubernetes-fetch-client";

export default Vue.extend({
  components: {
    Logo,
    VuetifyLogo,
  },
  data() {
    return {
      apps: [],
      headers: [
        {
          text: "Name",
          value: "metadata.name",
        },
        {
          text: "Domain",
          value: "status.domains",
          align: 'right'
        },
        {
          text: "Status",
          value: "status.status",
          align: 'right'
        },
        {
          text: "Actions",
          value: "actions",
          align: 'right'
        },
      ]
    } as {
      apps: DevTsuzuModokiV1alpha1Application[],
    }
  },

  async created() {
    const conf = new Configuration({
      fetchApi: fetch,
    })

    const modokiApi = new ModokiTsuzuDevV1alpha1Api(conf);

    this.apps = (await modokiApi.listApplicationForAllNamespaces({})).items;
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
  }
})
</script>

<style>
.text-end {
  text-align: end;
}
</style>
