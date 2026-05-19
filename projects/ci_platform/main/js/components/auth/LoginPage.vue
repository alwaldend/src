<template>
  <div>
    <div class="text-center mb-5">
      <div class="text-900 text-3xl font-medium mb-3">Welcome Back</div>
      <span class="text-600 font-medium line-height-3">
        Don't have an account?
      </span>
      <router-link
        class="font-medium no-underline ml-2 text-blue-500 cursor-pointer"
        :to="{ name: 'auth_registration' }"
      >
        Create today!
      </router-link>
    </div>
    <form
      class="p-fluid w-full md:w-10 mx-auto"
      @submit.prevent="handleSubmit(validator.$invalid)"
    >
      <div class="field">
        <div class="p-float-label">
          <InputText
            id="username1"
            autofocus
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
        <div class="p-float-label">
          <Password
            id="password1"
            v-model="validator.password.$model"
            :toggleMask="true"
            :feedback="false"
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

      <div class="flex align-items-center justify-content-between mb-5">
        <div class="flex align-items-center">
          <Checkbox
            id="rememberme1"
            v-model="remember_me"
            :binary="true"
            class="mr-2"
          ></Checkbox>
          <label for="rememberme1">Remember me</label>
        </div>
        <router-link
          class="font-medium no-underline ml-2 text-right cursor-pointer"
          style="color: var(--primary-color)"
          :to="{ name: 'auth_forgot_password' }"
        >
          Forgot password?
        </router-link>
      </div>
      <Button
        label="Sign In"
        type="submit"
        class="w-full p-3 text-xl"
        :loading="submitted"
      ></Button>
    </form>
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
  data() {
    return {
      username: "",
      password: "",
      remember_me: false,
      submitted: false,
    };
  },
  validations() {
    return {
      username: {
        required,
      },
      password: {
        required,
      },
    };
  },
  computed: {
    isSubmitted() {
      return this.$data.submitted;
    },
  },
  methods: {
    handleSubmit(isFormInvalid) {
      if (this.submitted || isFormInvalid) {
        return;
      }
      this.submitted = true;
      this.$store
        .dispatch("user/login", {
          username: this.username,
          password: this.password,
          remember_me: this.remember_me,
        })
        .then(() => this.$router.push({ name: "root" }))
        .catch((error) => Error(error, this.$toast.add))
        .finally(() => (this.submitted = false));
    },
  },
};
</script>
