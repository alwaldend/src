<template>
  <div>
    <div class="card" v-show="!projects">
      <Skeleton height="1rem" width="100%" />
      <Skeleton height="1rem" width="100%" />
      <Skeleton height="1rem" width="100%" />
    </div>
    <div class="card" v-for="project in projects" :key="project.id">
      <div class="flex flex-column md:flex-row align-items-center p-3 w-full">
        <Avatar
          :label="project.name[0]"
          :image="project.image"
          size="xlarge"
          class="my-4 md:my-0 shadow-2 mr-5"
        />
        <div class="flex-1 text-center md:text-left">
          <router-link
            class="font-bold text-2xl"
            :to="{ name: 'project', params: { project_path: project.path } }"
            >{{ project.path.replace("/", " / ") }}
          </router-link>
          <div class="mb-3">{{ project.description }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../../lib/main/Axios";

import Error from "../../lib/main/Error";

export default {
  data() {
    return {
      projects: [],
    };
  },
  mounted() {
    axios
      .get("/api/v1/projects")
      .then((response) => (this.$data.projects = response.data.data))
      .catch(Error);
  },
};
</script>
