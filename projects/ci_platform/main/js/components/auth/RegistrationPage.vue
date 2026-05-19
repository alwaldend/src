<template>
  <div>
    <div class="text-center mb-5">
      <div class="text-900 text-3xl font-medium mb-3">Register</div>
    </div>
    <form
      class="p-fluid w-full md:w-10 mx-auto"
      @submit.prevent="handleSubmit(validator.$invalid)"
    >
      <div class="field">
        <div class="p-float-label">
          <InputText
            id="username1"
            v-model="validator.username.$model"
            :class="{ 'p-invalid': validator.username.$invalid }"
          />
          <label
            for="username1"
            :class="{ 'p-error': validator.username.$invalid }"
          >
            Username
          </label>
        </div>
        <small
          v-if="
            validator.username.$invalid ||
            validator.username.$pending.$response
          "
          class="p-error"
          >{{
            validator.username.required.$message.replace("Value", "It")
          }}</small
        >
      </div>

      <div class="field">
        <div class="p-float-label p-input-icon-right">
          <i class="pi pi-envelope" />
          <InputText
            id="email1"
            v-model="validator.email.$model"
            :class="{ 'p-invalid': validator.email.$invalid }"
          />
          <label for="email1" :class="{ 'p-error': validator.email.$invalid }">
            Email
          </label>
        </div>
        <span v-if="validator.email.$error">
          <span v-for="(error, index) of validator.email.$errors" :key="index">
            <small class="p-error">{{
              error.$message.replace("Value", "It")
            }}</small>
          </span>
        </span>
        <small
          v-else-if="
            validator.email.$invalid || validator.email.$pending.$response
          "
          class="p-error"
          >{{
            validator.email.required.$message.replace("Value", "It")
          }}</small
        >
      </div>

      <div class="field">
        <div class="p-float-label">
          <Password
            id="password1"
            v-model="validator.password.$model"
            :toggleMask="true"
            :feedback="true"
            :class="{ 'p-invalid': validator.password.$invalid }"
          ></Password>
          <label
            for="password1"
            :class="{ 'p-error': validator.password.$invalid }"
            >Password</label
          >
        </div>
        <small
          v-if="
            validator.password.$invalid ||
            validator.password.$pending.$response
          "
          class="p-error"
          >{{
            validator.password.required.$message.replace("Value", "Password")
          }}</small
        >
      </div>

      <Button label="Submit" type="submit" class="w-full p-3 text-xl"></Button>
    </form>
  </div>
</template>

<style lang="scss" scoped>
.field {
  margin-bottom: 2rem;
}
</style>

<script>
import { email, required } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import axios from "axios";
import Error from "../../lib/main/Error";

export default {
  setup: () => ({ validator: useVuelidate() }),
  data() {
    return {
      username: "",
      password: "",
      email: "",
      submitted: false,
    };
  },
  validations() {
    return {
      username: {
        required,
      },
      email: {
        required,
        email,
      },
      password: {
        required,
      },
    };
  },
  methods: {
    handleSubmit(isFormInvalid) {
      if (this.submitted || isFormInvalid) {
        return;
      }
      this.submitted = true;

      axios
        .post("/api/v1/users", {
          username: this.username,
          password: this.password,
          email: this.email,
        })
        .then(() => this.$router.push({ name: "auth_login" }))
        .catch((error) => Error(error, this.$toast.add))
        .finally(() => (this.submitted = false));
    },
  },
};
</script>
