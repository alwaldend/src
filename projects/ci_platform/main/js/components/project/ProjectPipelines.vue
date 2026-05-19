<template>
  <div class="card">
    <Dropdown
      v-model="currentBranchData"
      :options="branches"
      optionLabel="name"
      :filter="true"
      :loading="loadingBranches"
      @show.once="loadBranches"
      placeholder="Select a branch"
    >
      <template #value="slotProps">
        <div v-if="slotProps.value">
          {{ slotProps.value.name }}
        </div>
        <span v-else>
          {{ slotProps.placeholder }}
        </span>
      </template>
      <template #option="slotProps">
        <div>{{ slotProps.option.name }}</div>
      </template>
    </Dropdown>
    <Button
      label="Update"
      class="p-button-outlined mx-3"
      :loading="updatingRemoteBranch"
      @click="updateRemoteBranch"
    ></Button>
  </div>
  <div v-if="loadingBranch" class="card">
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
    <Skeleton class="mb-2" borderRadius="16px"></Skeleton>
  </div>
  <DataTable
    v-if="!loadingBranch"
    :value="pipelines"
    responsiveLayout="scroll"
    class="card"
  >
    <Column field="name" header="Name">
      <template #body="slotProps">
        <router-link
          :to="{
            name: 'project_pipeline',
            params: {
              pipeline_branch: currentBranch.name,
              pipeline_name: slotProps.data.name,
            },
          }"
        >
          <span>{{ slotProps.data.name }}</span>
        </router-link>
      </template>
    </Column>
  </DataTable>
</template>

<script>
import { mapGetters } from "vuex";

import Error from "../../lib/main/Error";

export default {
  data() {
    return {
      currentBranchData: {},
      branchesData: [],
      loadingBranch: false,
      loadingBranches: false,
      updatingRemoteBranch: false,
    };
  },
  computed: {
    ...mapGetters("project", {
      defaultBranch: "currentDefaultBranchName",
      projectID: "currentProjectID",
      branchesMap: "currentBranches",
    }),
    currentBranch() {
      return this.currentBranchData;
    },
    pipelines() {
      return this.currentBranchData.config
        ? this.currentBranchData.config.pipelines
        : [];
    },
    branches() {
      return this.branchesData;
    },
  },
  watch: {
    dropdownChosenBranch(newBranch) {
      this.loadBranch(newBranch.name);
    },
  },
  methods: {
    updateRemoteBranch() {
      if (this.updatingRemoteBranch) {
        return;
      }
      this.updatingRemoteBranch = true;
      this.$store
        .dispatch("project/updateRemoteBranch", this.currentBranchData.name)
        .then(() =>
          this.$toast.add({
            severity: "success",
            summary: "branch updated",
            life: 3000,
          })
        )
        .catch(Error)
        .finally(() => (this.updatingRemoteBranch = false));
    },
    loadBranch(branchName, setCurrent) {
      if (this.loadingBranch) {
        return;
      }
      this.loadingBranch = true;
      this.$store
        .dispatch("project/loadBranch", branchName)
        .then((response) => {
          if (setCurrent) {
            this.currentBranchData = response.data.data;
          }
          if (!this.branchesData.length) {
            this.branchesData = [this.currentBranchData];
          }
          this.loadingBranch = false;
        })
        .catch(Error);
    },
    loadBranches() {
      if (this.loadingBranches) {
        return;
      }
      this.loadingBranches = true;
      this.$store
        .dispatch("project/loadBranches")
        .then((response) => {
          this.loadingBranches = false;
          this.branchesData = response.data.data;
        })
        .catch(Error);
    },
  },
  beforeMount() {
    this.loadBranch(this.defaultBranch, true);
  },

  beforeRouteUpdate() {
    this.loadBranch(this.defaultBranch, true);
  },
};
</script>
