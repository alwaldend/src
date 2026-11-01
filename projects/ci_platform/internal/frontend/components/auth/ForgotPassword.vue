<template>
  <div class="w-full md:w-10 mx-auto">
    <label for="email1" class="block text-900 text-xl font-medium mb-2"
      >Email</label
    >
    <InputText
      id="email1"
      v-model="email"
      type="text"
      class="w-full mb-3"
      placeholder="Email"
      style="padding: 1rem"
    />

    <Button label="Send recover email" class="w-full p-3 text-xl"></Button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      password: "",
      remember_password: false,
    };
  },
};
</script>
