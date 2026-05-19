<template>
  <router-view></router-view>
  <Toast position="top-center" />
</template>

<script>
import { mapState } from "vuex";

import EventBus from "../../lib/main/EventBus";

export default {
  computed: mapState("theme", {
    theme: "theme",
  }),
  mounted() {
    EventBus.on("app-toast-add", (toast_message) => {
      this.$toast.add(toast_message);
    });
  },
  watch: {
    theme(newValue, oldValue) {
      const elementId = "theme-link";
      const linkElement = document.getElementById(elementId);
      const cloneLinkElement = linkElement.cloneNode(true);
      const newThemeUrl = linkElement
        .getAttribute("href")
        .replace(oldValue.name, newValue.name);

      cloneLinkElement.setAttribute("id", elementId + "-clone");
      cloneLinkElement.setAttribute("href", newThemeUrl);
      cloneLinkElement.addEventListener("load", () => {
        linkElement.remove();
        cloneLinkElement.setAttribute("id", elementId);
      });
      linkElement.parentNode.insertBefore(
        cloneLinkElement,
        linkElement.nextSibling
      );
    },
  },
};
</script>
