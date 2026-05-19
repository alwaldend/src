<template>
  <div v-if="loadingPipeline" class="card">
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
  </div>
  <div v-if="!loadingPipeline" class="card">
    <div class="flex justify-content-between">
      <div class="font-medium text-3xl text-900 mb-3">{{ pipeline.name }}</div>
      <Button
        label="Run"
        class="p-button-outlined"
        :loading="pipelineLaunched"
        @click="createPipeline"
      ></Button>
    </div>
    <div class="text-500 mb-5">
      {{ pipeline.description ? pipeline.description : "no description" }}
    </div>
    <DataTable :value="jobs" responsiveLayout="scroll">
      <Column field="name" header="Name">
        <template #body="slotProps">
          <span>{{ slotProps.data.name }}</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import Error from "../../lib/main/Error";

export default {
  data() {
    return {
      pipelineData: {},
      pipelineLaunched: false,
      loadingPipeline: false,
      branchData: {},
    };
  },
  computed: {
    pipeline() {
      return this.pipelineData;
    },
    ...mapGetters("project", {
      branch: "branch",
      project_id: "id",
    }),
    jobs() {
      return this.branchData.config.jobs;
    },
  },
  methods: {
    createPipeline() {
      if (this.pipelineLaunched) {
        return;
      }
      this.pipelineLaunched = true;
      this.$store
        .dispatch("project/createPipeline", {
          branchName: this.$route.params.pipeline_branch,
          pipelineName: this.$route.params.pipeline_name,
        })
        .then((response) =>
          this.$toast.add({
            severity: "success",
            summary: "pipeline launched",
            detail: response.data.web_url,
            life: 10000,
          })
        )
        .catch(Error)
        .finally(() => (this.pipelineLaunched = false));
    },
    loadPipeline() {
      this.loadingPipeline = true;
      this.$store
        .dispatch("project/loadBranch", this.$route.params.pipeline_branch)
        .then((response) => {
          this.branchData = response.data.data;
          response.data.data.config.pipelines.some((element) => {
            if (element.name !== this.$route.params.pipeline_name) {
              return false;
            }
            this.pipelineData = element;
            return true;
          });
          this.loadingPipeline = false;
        })
        .catch(Error);
    },
  },
  beforeMount() {
    this.loadPipeline();
  },
  beforeRouteUpdate() {
    this.loadPipeline();
  },
};
</script>
