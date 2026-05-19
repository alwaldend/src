<template>
  <div
    class="flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden"
  >
    <div class="grid justify-content-center p-2 lg:p-0" style="min-width: 80%">
      <div
        class="surface-card p-4 border-round w-full lg:w-6"
        :class="shadowClass"
      >
        <div
          class="col-12"
          style="
            border-radius: 56px;
            padding: 0.3rem;
            background: linear-gradient(
              180deg,
              var(--primary-color),
              rgba(33, 150, 243, 0) 30%
            );
          "
        >
          <div
            class="h-full w-full m-0 py-7 px-4"
            style="
              border-radius: 53px;
              background: linear-gradient(
                180deg,
                var(--surface-50) 38.9%,
                var(--surface-0)
              );
            "
          >
            <div class="text-center mb-5">
              <img
                :src="logoSrc"
                alt="Tianyi logo"
                class="mb-5"
                style="width: 81px; height: 60px"
              />
            </div>
            <div @form-page-logo-color="setLogoColor">
              <router-view></router-view>
              <div
                v-if="$route.meta.showGoBackLink"
                class="col-12 mt-5 text-center"
              >
                <i
                  class="pi pi-fw pi-arrow-left text-blue-500 mr-2"
                  style="vertical-align: center"
                ></i>
                <a @click="$router.go(-1)" href="#" class="text-blue-500"
                  >Go Back</a
                >
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    setLogoColor(color) {
      console.log(color);
      this.logo_color = color;
    },
  },
  data() {
    return { logo_color: "" };
  },
  computed: {
    shadowClass() {
      return (
        "shadow-" +
        (this.$route.meta.shadow ? this.$route.meta.shadow : this.shadow)
      );
    },
    logoSrc() {
      var color;
      if (this.$route.meta.logoColor) {
        color = this.$route.meta.logoColor;
      } else if (this.$store.getters["theme/isDark"]) {
        color = "white";
      } else {
        color = "dark";
      }
      return "/layout/images/logo-" + color + ".svg";
    },
  },
  props: {
    shadow: {
      type: Number,
      default: 2,
    },
  },
};
</script>

<style lang="scss">
.pi-eye {
  transform: scale(1.6);
  margin-right: 1rem;
}

.pi-eye-slash {
  transform: scale(1.6);
  margin-right: 1rem;
}
</style>
