<template>
  <v-dialog v-model="dialog" persistent max-width="600px">
    <v-form ref="form" v-model="valid">
    <v-card>
      <v-card-title>
        <span class="headline">New Application</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field v-model="name" label="Name" required></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="domains" label="Domains" required hint="separated by comma"></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="image" label="Image"></v-text-field>
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
import {DevTsuzuModokiV1alpha1Application} from "@modoki-paas/kubernetes-fetch-client"

export default Vue.extend({
  props: {
    dialog: Boolean,
  },
  data() {
    return {
      valid: false,
      name: "",
      domains: "",
      image: "",
      command: "",
      args: "",
      attributes: [] as {index: number; key: string; value: string}[],
    }
  },
  methods: {
    parseApplication(): DevTsuzuModokiV1alpha1Application {
      const attributes = Object.fromEntries(this.attributes.map(m => [m.key, m.value]));

      return {
        apiVersion: "modoki.tsuzu.dev/v1alpha1",
        kind: "Application",
        metadata: {
          name: this.name,
          namespace: this.$route.params.namespace,
        },
        spec: {
          args: this.args.length ? this.args.split(",") : undefined,
          command: this.command.length ? this.command.split(",") : undefined,
          domains: this.domains.split(","),
          image: this.image,
          attributes: this.attributes.length ? attributes : undefined,
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
        this.reset()
      }
    }
  }
})
</script>
