<template>
  <div class="grid">
    <div class="col-3 sm:col-0"></div>
    <div class="col-6">
      <div class="card p-fluid">
        <div class="field grid">
          <label for="field_id">ID</label>
          <InputText id="field_id" v-model="id" type="text" />
        </div>
        <div class="field grid">
          <label for="field_admin">Admin</label>
          <InputText id="field_admin" v-model="admin" type="text" />
        </div>
        <div class="field grid">
          <label for="field_created_at">Created at</label>
          <InputText id="field_created_at" v-model="createdAt" type="text" />
        </div>
        <div class="field grid">
          <label for="field_updated_at">Updated at</label>
          <InputText id="field_updated_at" v-model="updatedAt" type="text" />
        </div>
        <div class="field grid">
          <label for="field_username">Username</label>
          <InputText id="field_username" v-model="username" type="text" />
        </div>
        <div class="field grid">
          <label for="field_email">Email</label>
          <InputText id="field_email" v-model="email" type="text" />
        </div>
        <div class="field flex flex-column">
          <label for="field_image">Image</label>
          <Avatar icon="pi pi-user" :image="image" size="xlarge" />
        </div>
        <div class="field grid">
          <label for="field_roles">Roles</label>
          <InputText id="field_roles" v-model="roles" type="text" />
        </div>
        <div class="field grid">
          <label for="field_token">Token</label>
          <InputText id="field_token" v-model="token" type="text" />
        </div>
      </div>
    </div>
    <div class="col-3 sm:col-0"></div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters("user", {
      id: "id",
      admin: "admin",
      createdAt: "created_at",
      updatedAt: "updated_at",
      username: "username",
      email: "email",
      bio: "bio",
      image: "image",
      roles: "roles",
      token: "token",
      user: "user",
    }),
  },
};
</script>
