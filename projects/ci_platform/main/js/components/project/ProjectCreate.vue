<template>
  <div
    class="flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden"
  >
    <div class="grid justify-content-center p-2 lg:p-0" style="min-width: 80%">
      <div
        class="surface-card p-4 border-round w-full lg:w-6 shadow-2 flex justify-content-center"
      >
        <form @submit.prevent="handleSubmit(validator.$invalid)" class="my-2">
          <div class="field">
            <div class="p-float-label">
              <InputText
                id="project_name"
                autofocus
                v-model="validator.project_name.$model"
                :class="{ 'p-invalid': validator.project_name.$invalid }"
              />
              <label
                for="project_name"
                :class="{ 'p-error': validator.project_name.$invalid }"
              >
                Project name
              </label>
            </div>
            <small
              v-if="
                validator.project_name.$invalid ||
                validator.project_name.$pending.$response
              "
              class="p-error"
              >{{
                validator.project_name.required.$message.replace("Value", "It")
              }}</small
            >
          </div>

          <div class="field">
            <div class="p-float-label">
              <InputText
                id="project_source"
                autofocus
                v-model="validator.project_source.$model"
                :class="{ 'p-invalid': validator.project_source.$invalid }"
              />
              <label
                for="project_source"
                :class="{ 'p-error': validator.project_source.$invalid }"
              >
                Project source
              </label>
            </div>
            <small
              v-if="
                validator.project_source.$invalid ||
                validator.project_source.$pending.$response
              "
              class="p-error"
              >{{
                validator.project_source.required.$message.replace(
                  "Value",
                  "It"
                )
              }}</small
            >
          </div>

          <div class="field">
            <div class="p-float-label">
              <InputText
                id="project_default_branch"
                v-model="project_default_branch"
              />
              <label for="project_default_branch"> Default branch </label>
            </div>
          </div>

          <Button
            label="Create project"
            type="submit"
            class="p-3 text-xl"
            :loading="submitted"
          ></Button>
        </form>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.field {
  margin-bottom: 2rem;
}
</style>

<script>
import { required } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import Error from "../../lib/main/Error";

export default {
  setup: () => ({ validator: useVuelidate() }),
  validations() {
    return {
      project_name: {
        required,
      },
      project_source: {
        required,
      },
    };
  },
  data() {
    return {
      submitted: false,
      project_name: "",
      project_source: "",
      project_default_branch: "master",
    };
  },
  methods: {
    handleSubmit(isFormInvalid) {
      if (this.submitted || isFormInvalid) {
        return;
      }
      this.submitted = true;
      this.$store
        .dispatch("project/createProject", {
          name: this.project_name,
          source: this.project_source,
          defaultBranch: this.project_default_branch,
        })
        .then((response) =>
          this.$router.push({
            name: "project",
            params: {
              project_path: response.data.data.path,
            },
          })
        )
        .catch(Error)
        .finally(() => (this.submitted = false));
    },
  },
};
</script>
