<template>
  <div class="layout-topbar">
    <router-link :to="{ name: 'root' }" class="layout-topbar-logo mx-1 p-2">
      <img alt="Logo" :src="$store.getters['theme/logo']" />
    </router-link>
    <Menubar :model="topbarLeftMenuItems" class="flex-grow-1" />
    <Menubar :model="topbarCenterMenuItems" class="layout-topbar-menu-right" />
    <span class="p-input-icon-left">
      <i class="pi pi-search" />
      <InputText type="text" v-model="topbarSearch" placeholder="Search" />
    </span>
    <Menubar :model="topbarRightMenuItems" class="layout-topbar-menu-right" />

    <Dialog
      header="Confirmation"
      v-model:visible="displayLogout"
      :style="{ width: '350px' }"
      :modal="true"
    >
      <div class="flex align-items-center justify-content-center">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span>Are you sure you want to log out?</span>
      </div>
      <template #footer>
        <Button
          label="No"
          icon="pi pi-times"
          @click="logoutAction"
          class="p-button-text"
        />
        <Button
          label="Yes"
          icon="pi pi-check"
          @click="logoutAction(true)"
          class="p-button-text"
          autofocus
        />
      </template>
    </Dialog>
  </div>
</template>

<script>
export default {
  methods: {
    logoutAction(logout) {
      this.$data.displayLogout = false;
      if (!logout) {
        return;
      }
      this.$store.dispatch("user/clean");
      this.$router.push({ name: "auth_login" });
    },
  },
  data() {
    return {
      topbarSearch: "",
      displayLogout: false,
      topbarLeftMenuItems: [
        {
          label: "Browse",
          icon: "pi pi-bars",
          items: [
            {
              label: "Projects",
              icon: "pi pi-folder",
              to: { name: "browse_projects" },
            },
            {
              label: "Groups",
              icon: "pi pi-share-alt",
              to: { name: "browse_groups" },
            },
          ],
        },
        {
          label: "Admin",
          icon: "pi pi-shield",
          items: [],
          visible: () => {
            return this.$store.getters["user/admin"];
          },
        },
      ],
      topbarCenterMenuItems: [
        {
          icon: "pi pi-plus-circle",
          items: [
            {
              label: "New project",
              to: { name: "project_create" },
            },
          ],
        },
      ],
      topbarRightMenuItems: [
        {
          icon: "pi pi-question-circle",
          items: [
            {
              label: "API",
              icon: "pi pi-key",
              to: { name: "help_swagger" },
            },
            {
              label: "Source code",
              icon: "pi pi-external-link",
              url: "https://gitlab.com/kongrentian-group/ci_platform",
            },
          ],
        },
        {
          icon: "pi pi-user",
          items: [
            {
              label: "Profile",
              to: { name: "user_profile" },
            },
            {
              label: "Preferences",
              to: { name: "user_appearance" },
            },
            {
              separator: true,
            },
            {
              label: "Logout",
              command: () => (this.$data.displayLogout = true),
            },
          ],
        },
      ],
    };
  },
};
</script>
