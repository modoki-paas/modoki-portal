<template>
  <v-dialog v-model="dialog" persistent max-width="600px">
    <v-form ref="form" v-model="valid">
    <v-card>
      <v-card-title>
        <span class="headline">New RemoteSync</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field v-model="name" label="Name" required :rules="[inputtedRule]
              "></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-select
                :items="applicationNames"
                :disabled="applicationNames.length === 0"
                label="Application"
                required
                v-model="appName"
                :rules="[selectedRule]"
              ></v-select>
            </v-col>



            <v-col cols="12" sm="6">
              <v-select
                :items="owners"
                :disabled="owners.length === 0"
                label="Owner"
                required
                v-model="owner"
                :rules="[selectedRule]"
              ></v-select>
            </v-col>
            <v-col cols="12" sm="6">
              <v-select
                :items="repositoriesNames"
                :disabled="repositoriesNames.length === 0"
                label="Repository"
                required
                v-model="repository"
                :rules="[selectedRule]"
              ></v-select>
            </v-col>

            <v-col cols="12">
              <v-radio-group v-model="mode" row :rules="[selectedRule]">
                <v-radio label="Branch" value="branch"></v-radio>
                <v-radio label="Pull Request" value="pr"></v-radio>
                <v-radio label="Revision" value="sha"></v-radio>
              </v-radio-group>
            </v-col>

            <v-col cols="12">
              <v-text-field
                :label="githubBaseModeLabel"
                :type="mode==='pr' ? 'number' : undefined"
                :rules="mode==='pr' ? [integer] : undefined"
                :prepend-icon="mode==='pr' ? 'mdi-pound' : undefined"
              >
              </v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field v-model="subPath" label="Sub Path"></v-text-field>
            </v-col>

          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="cancel">Close</v-btn>
        <v-btn color="blue darken-1" :disabled="!valid" text @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
    </v-form>
  </v-dialog>
</template>

<script lang="ts">
import Vue from 'vue'
import { fetch } from "~/util/proxy";
import {
  DevTsuzuModokiV1alpha1Application,
  ModokiTsuzuDevV1alpha1Api,
  Configuration,
  ConfigurationParameters,
  DevTsuzuModokiV1alpha1RemoteSync
} from "@modoki-paas/kubernetes-fetch-client"
import {Installation, Repository, listInstallations, listRepositories} from "~/util/installations"

export default Vue.extend({
  props: {
    dialog: Boolean,
  },
  data() {
    return {
      mode: "branch",
      modokiApi: undefined as (ModokiTsuzuDevV1alpha1Api | undefined),
      applications: [] as DevTsuzuModokiV1alpha1Application[],
      valid: false,
      name: "",
      installations: [] as Installation[],
      repositories: [] as Repository[],
      owner: "",
      repository: "",
      appName: "",
      subPath: "",
      inputtedRule: (value: any) => !!value || "Cannot be empty",
      selectedRule: (value: any) => !!value || "選択してください",
      integer: (value: string) => /^[1-9][0-9]*$/.test(value) && value.length < 9 || "Invalid Number"
    }
  },
  async mounted() {
    const conf = new Configuration({
      fetchApi: fetch,
    })

    const modokiApi = new ModokiTsuzuDevV1alpha1Api(conf);
    this.modokiApi = modokiApi;

    await this.getInstallation();
    await this.getApplications();
  },
  computed: {
    owners(): string[] {
      return this.installations.map((ins: Installation) => ins.account.login)
    },
    repositoriesNames(): string[] {
      return (this.repositories as Repository[]).map(r => r.name);
    },
    applicationNames(): (string | undefined)[] {
      return (this.applications as DevTsuzuModokiV1alpha1Application[]).map(app => app.metadata?.name).filter(x => x);
    },
    githubBaseModeLabel() {
      if(this.mode === "branch") {
        return "Branch Name"
      }else if(this.mode === "pr") {
        return "Pull Request Number"
      }else {
        return "Commit Hash"
      }
    }
  },
  methods: {
    async getInstallation() {
      this.installations = await listInstallations();
    },
    async getRepositories() {
      if(this.owner.length === 0 ) {
        this.repositories = [];
      }else {
        const installationID = this.installations.filter(ins => ins.account.login === this.owner)[0].id;
        this.repositories = await listRepositories(installationID);
      }
    },
    async getApplications() {
      this.applications = (await this.modokiApi!.listNamespacedApplication({
        namespace: this.$route.params.namespace,
      })).items;
    },
    parseApplication(): DevTsuzuModokiV1alpha1RemoteSync {
      return {
        apiVersion: "modoki.tsuzu.dev/v1alpha1",
        kind: "RemoteSync",
        metadata: {
          name: this.name,
          namespace: this.$route.params.namespace,
        },
        spec: {
          applicationRef: {
            name: this.appName,
          },
          base: {
            github: {
              owner: this.owner,
              repo: this.repository,
              secretName: "tsuzu",
            },
            subPath: this.subPath.length ? this.subPath : undefined,
          },
          image: {
            name: `registry.misw.jp/${this.$route.params.namespace}/${this.name}`,
            secretName: "misw-registry",
          }
        },
      }
    },
    cancel() {
      this.$emit("close", undefined)
    },
    save() {
      this.$emit("close", this.parseApplication())
    },
    reset() {
      if(this.$refs.form)
        (this.$refs.form as any).reset();
    }
  },
  watch: {
    dialog: function(after: boolean, before: boolean) {
      if(!before && after) {
        this.reset();
      }
    },
    owner: async function(after: string, before: string) {
      if(before !== after) {
        await this.getRepositories();
      }
    },
  }
})
</script>
