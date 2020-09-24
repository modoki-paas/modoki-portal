<template>
  <v-app dark>
    <v-app-bar fixed app>
      <v-btn icon>
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-btn icon v-if="prev" @click="routeToPrev">
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>
      <v-toolbar-title v-text="title" />
      <v-spacer />

      <template v-slot:extension v-if="tabs.length !== 0">
        <v-tabs
          v-model="model"
          centered
        >
          <v-tab
            v-for="tab in tabs"
            :key="tab.key"
            :nuxt="true"
            :to="{path: tabsPath, hash: `#${tab.key}`}"
          >
          {{tab.name}}
          </v-tab>
        </v-tabs>
      </template>
    </v-app-bar>
    <v-main>
      <v-container>
        <nuxt />
      </v-container>
    </v-main>
    <v-footer absolute app>
      <span>&copy; 2020 modoki</span>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      title: 'modoki portal',
      prev: undefined as (()=>{}) | undefined,
      tabs: [] as {
        name: string;
        key: string;
      }[],
      tabsPath: "/",
      model: "",
    }
  },
  created() {
    this.$nuxt.$on('headerInfo', this.setHeaderInfo);
  },
  methods: {
    setHeaderInfo(prop: {title: string; prev?: ()=>{}; path: string; tabs: {name: string; key: string;}[]}) {
      this.title = prop.title;
      this.tabs = prop.tabs;
      this.prev = prop.prev;
      this.tabsPath = prop.path;
    },
    routeToPrev() {
      if(this.prev) {
        this.prev();
      }
    }
  }
})
</script>
