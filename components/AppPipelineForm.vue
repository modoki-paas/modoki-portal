<template>
  <v-dialog v-model="dialog" persistent max-width="600px">
    <v-form ref="form" v-model="valid">
    <v-card>
      <v-card-title>
        <span class="headline">New AppPipeline</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field v-model="name" label="Name" required :rules="[inputtedRule]
              "></v-text-field>
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
              <v-text-field v-model="subPath" label="Sub Path"></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field v-model="command" label="Commands" hint="separated by comma"></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="args" label="Args" hint="separated by comma"></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-row>
                <v-col>
                <span class="headline text-subtitle-1">Attributes</span>
                </v-col>
                <v-spacer></v-spacer>
                <v-btn icon @click="minus">
                  <v-icon>mdi-minus</v-icon>
                </v-btn>
                <v-btn icon @click="plus">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </v-row>
            </v-col>
            <v-col cols="12" class="py-0">
              <v-row v-for="attr in attributes" :key="attr.index">
                <v-col class="py-0" cols="12" sm="12" md="6">
                  <v-text-field class="pt-1" v-model="attr.key" label="Key"></v-text-field>
                </v-col>
                <v-col class="py-0" cols="12" sm="12" md="6">
                  <v-text-field class="pt-1" v-model="attr.value" label="Value"></v-text-field>
                </v-col>
              </v-row>
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
  DevTsuzuModokiV1alpha1AppPipeline
} from "@modoki-paas/kubernetes-fetch-client"
import {Installation, Repository, listInstallations, listRepositories} from "~/util/installations";

const namespace = "default"

export default Vue.extend({
  props: {
    dialog: Boolean,
  },
  data() {
    return {
      valid: false,
      name: "",
      installations: [] as Installation[],
      repositories: [] as Repository[],
      owner: "",
      repository: "",
      subPath: "",
      command: "",
      args: "",
      attributes: [] as {index: number; key: string; value: string}[],
      inputtedRule: (value: any) => !!value || "Cannot be empty",
      selectedRule: (value: any) => !!value || "選択してください",
      integer: (value: string) => /^[1-9][0-9]*$/.test(value) && value.length < 9 || "Invalid Number"
    }
  },
  async mounted() {
    const conf = new Configuration({
      fetchApi: fetch,
    })

    await this.getInstallation();
  },
  computed: {
    owners() {
      return (this as any).installations.map((ins: Installation) => ins.account.login)
    },
    repositoriesNames() {
      return ((this as any).repositories as Repository[]).map(r => r.name);
    },
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
    parseApplication(): DevTsuzuModokiV1alpha1AppPipeline {
      const attributes = Object.fromEntries(this.attributes.map(m => [m.key, m.value]));

      return {
        apiVersion: "modoki.tsuzu.dev/v1alpha1",
        kind: "AppPipeline",
        metadata: {
          name: this.name,
          namespace: namespace,
        },
        spec: {
          domainBase: "*.modoki.misw.jp",
          applicationTemplate: {
            spec: {
              args: this.args?.length ? this.args.split(",") : undefined,
              command: this.command?.length ? this.command.split(",") : undefined,
              attributes: this.attributes?.length ? attributes : undefined,
            }
          },
          base: {
            github: {
              owner: this.owner,
              repo: this.repository,
              secretName: "tsuzu",
            },
            subPath: this.subPath?.length ? this.subPath : undefined,
          },
          image: {
            name: `registry.misw.jp/${namespace}/${this.name}`,
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
    plus() {
      this.attributes.push({
        index: this.attributes.length,
        key: "",
        value: "",
      })
    },
    minus() {
      this.attributes.pop();
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
