<template>
  <Breadcrumb
    :home="{ icon: 'pi pi-home', to: { name: 'root' } }"
    :model="path"
    class="card"
  />
  <div v-if="loadingProject" class="card">
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
  </div>
  <div v-if="!loadingProject">
    <div
      class="card flex align-items-start flex-column lg:justify-content-between lg:flex-row"
    >
      <div>
        <div class="font-medium text-3xl text-900">
          {{ currentProjectName }}
        </div>
        <div class="flex align-items-center text-700 flex-wrap">
          <span>
            Project ID: {{ currentProjectID }}
            <Button
              icon="pi pi-copy"
              class="p-button-text p-button-plain h-1rem w-2rem"
              v-tooltip.bottom="'Copy project id'"
              @click="writeToClipboard(currentProjectID)"
            />
          </span>
        </div>
        <div class="flex align-items-center text-700 flex-wrap">
          <div class="mr-5 flex align-items-center mt-3">
            <i class="pi pi-globe mr-2"></i>
            <span>placeholder</span>
          </div>
          <div class="flex align-items-center mt-3">
            <i class="pi pi-clock mr-2"></i>
            <span>placeholder</span>
          </div>
        </div>
        <div
          v-if="currentProjectDescription"
          class="flex align-items-center text-700 flex-wrap mt-3"
        >
          <span>{{ currentProjectDescription }}</span>
        </div>
      </div>
    </div>
    <router-view />
  </div>
</template>

<script>
import { mapGetters } from "vuex";

import Error from "../../lib/main/Error";

export default {
  data() {
    return {
      loadingProject: false,
    };
  },
  beforeRouteUpdate() {
    this.loadProject();
  },
  beforeMount() {
    this.loadProject();
  },
  methods: {
    writeToClipboard: (text) => navigator.clipboard.writeText(text),
    loadProject() {
      this.loadingProject = true;
      // load the current project
      this.$store
        .dispatch("project/loadProject", this.$route.params.project_path)
        .then(() => (this.loadingProject = false))
        .catch(Error);
    },
  },
  computed: {
    ...mapGetters("project", {
      currentProjectID: "currentID",
      currentProjectName: "currentName",
      currentProjectDescription: "currentDescription",
    }),
    path() {
      var oldPath = "";
      return this.$route.params.project_path.split("/").map((element) => {
        return {
          label: element,
          to: {
            name: "project",
            params: {
              project_path: oldPath ? [oldPath, element].join("/") : element,
            },
          },
        };
      });
    },
  },
};
</script>
