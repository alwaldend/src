<template>
  <div :class="containerClass" @click="onWrapperClick">
    <AppTopBar />

    <div
      class="layout-sidebar flex flex-column"
      :class="sidebarCustomZIndex"
      @click="onSidebarClick"
    >
      <div class="layout-menu-container flex-grow-1">
        <AppSubmenu
          :items="sidebarItems"
          class="layout-menu"
          :root="$route.meta.sidebarAsRoot"
          @menuitem-click="onMenuItemClick"
          @keydown="onMenuKeyDown"
        />
      </div>
      <Button
        class="p-button-rounded p-button-primary p-button-outlined align-self-center"
        :icon="sidebarToggleIconClass"
        @click="onMenuToggle($event)"
      />
    </div>

    <div
      v-if="(!isSidebarVisible || mobileMenuActive) && sidebarItems"
      class="layout-sidebar-button-show-container"
    >
      <Button
        class="p-button-rounded p-button-primary"
        :icon="sidebarToggleIconClass"
        @click="onMenuToggle($event)"
      />
    </div>

    <div class="layout-main-container" :class="contentClass">
      <div class="layout-main">
        <router-view />
      </div>
      <AppFooter />
    </div>

    <transition name="layout-mask">
      <div
        class="layout-mask p-component-overlay"
        v-if="mobileMenuActive"
      ></div>
    </transition>
  </div>
</template>

<script>
import AppTopBar from "./AppTopbar.vue";
import AppFooter from "./AppFooter.vue";
import AppSubmenu from "./AppSubmenu.vue";
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      staticMenuInactive: false,
      overlayMenuActive: false,
      mobileMenuActive: false,
    };
  },
  watch: {
    $route() {
      this.menuActive = false;
      this.$toast.removeAllGroups();
    },
  },
  methods: {
    onMenuItemClick(event) {
      this.$emit("menuitem-click", event);
    },
    onMenuKeyDown(event) {
      const nodeElement = event.target;
      if (event.code === "Enter" || event.code === "Space") {
        nodeElement.click();
        event.preventDefault();
      }
    },
    onWrapperClick() {
      if (!this.menuClick) {
        this.overlayMenuActive = false;
        this.mobileMenuActive = false;
      }

      this.menuClick = false;
    },
    onMenuToggle(event) {
      this.menuClick = true;

      if (this.isDesktop()) {
        if (this.layoutMode === "overlay") {
          if (this.mobileMenuActive === true) {
            this.overlayMenuActive = true;
          }

          this.overlayMenuActive = !this.overlayMenuActive;
          this.mobileMenuActive = false;
        } else if (this.layoutMode === "static") {
          this.staticMenuInactive = !this.staticMenuInactive;
        }
      } else {
        this.mobileMenuActive = !this.mobileMenuActive;
      }

      event.preventDefault();
    },
    onSidebarClick() {
      this.menuClick = true;
    },
    addClass(element, className) {
      if (element.classList) element.classList.add(className);
      else element.className += " " + className;
    },
    removeClass(element, className) {
      if (element.classList) element.classList.remove(className);
      else
        element.className = element.className.replace(
          new RegExp(
            "(^|\\b)" + className.split(" ").join("|") + "(\\b|$)",
            "gi"
          ),
          " "
        );
    },
    isDesktop() {
      return window.innerWidth >= 992;
    },
  },
  computed: {
    ...mapGetters("theme", {
      layoutMode: "layoutMode",
    }),

    sidebarItems() {
      return this.$route.meta.sidebarItems;
    },
    sidebarToggleIconClass() {
      return (
        "pi " +
        (this.isSidebarVisible
          ? "pi-angle-double-left"
          : "pi-angle-double-right")
      );
    },
    // we need to dynamically add z-index, because if it's too small, then
    // it will no be able to cover menubar in the overlay mode,
    // and if it's too big, then in normal mode menubar popups will be covered
    sidebarCustomZIndex() {
      return this.overlayMenuActive ? "layout-sidebar-z-index-overlay" : "";
    },
    isSidebarVisible() {
      if (this.isDesktop()) {
        if (this.layoutMode === "static") return !this.staticMenuInactive;
        else if (this.layoutMode === "overlay") return this.overlayMenuActive;
      }

      return true;
    },
    containerClass() {
      return [
        "layout-wrapper",
        {
          "layout-overlay": this.layoutMode === "overlay",
          "layout-static": this.layoutMode === "static",
          "layout-static-sidebar-inactive":
            this.layoutMode === "static" &&
            (this.staticMenuInactive || !this.sidebarItems),
          "layout-overlay-sidebar-active":
            this.layoutMode === "overlay" &&
            this.overlayMenuActive &&
            this.sidebarItems,
          "layout-mobile-sidebar-active":
            this.mobileMenuActive && this.sidebarItems,
          "p-input-filled": this.$primevue.config.inputStyle === "filled",
          "p-ripple-disabled": this.$primevue.config.ripple === false,
        },
      ];
    },
  },
  beforeUpdate() {
    if (this.mobileMenuActive)
      this.addClass(document.body, "body-overflow-hidden");
    else this.removeClass(document.body, "body-overflow-hidden");
  },
  components: {
    AppTopBar,
    AppFooter,
    AppSubmenu,
  },
};
</script>

<style lang="scss">
.p-toast.p-toast-top-right {
  z-index: 1000;
  top: 7rem;
}
</style>
