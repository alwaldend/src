<template>
  <div class="grid">
    <div class="col-12 md:col-6">
      <div class="card">
        <h4>Component Scale</h4>
        <div class="config-scale">
          <Button
            icon="pi pi-minus"
            @click="decrementScale()"
            class="p-button-text"
            :disabled="scale === scales[0]"
          />
          <i
            class="pi pi-circle-fill"
            v-for="s of scales"
            :class="{ 'scale-active': s === scale }"
            :key="s"
          />
          <Button
            icon="pi pi-plus"
            @click="incrementScale()"
            class="p-button-text"
            :disabled="scale === scales[scales.length - 1]"
          />
        </div>
      </div>
      <div class="card">
        <AppInputStyleSwitch />
      </div>
    </div>

    <div class="col-12 md:col-6">
      <div class="card">
        <h4>Ripple Effect</h4>
        <InputSwitch
          :modelValue="rippleActive"
          @update:modelValue="onRippleChange"
        />
      </div>
      <div class="card">
        <h4>Menu Type</h4>
        <div class="p-formgroup-inline">
          <div class="field-radiobutton">
            <RadioButton
              id="static"
              name="layoutMode"
              value="static"
              v-model="layoutMode"
            />
            <label for="static">Static</label>
          </div>
          <div class="field-radiobutton">
            <RadioButton
              id="overlay"
              name="layoutMode"
              value="overlay"
              v-model="layoutMode"
            />
            <label for="overlay">Overlay</label>
          </div>
        </div>
      </div>
    </div>

    <div class="col-12">
      <div class="card">
        <h4>Themes</h4>
        <div class="grid free-themes">
          <div v-for="item in computedThemes" :key="item.title" class="col-6">
            <Fieldset :legend="item.title" :toggleable="true">
              <div class="flex flex-wrap justify-content-center">
                <div
                  class="flex flex-column align-items-center m-3"
                  v-for="item_inner in item.themes"
                  :key="item_inner.name"
                >
                  <Button
                    type="button"
                    class="p-button-text p-button-plain"
                    @click="changeTheme(item_inner.name, item_inner.dark)"
                  >
                    <img :src="item_inner.source" :alt="item_inner.name" />
                  </Button>
                  <span>{{ item_inner.title }}</span>
                </div>
              </div>
            </Fieldset>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    changeTheme(theme, dark) {
      this.$store.commit("theme/theme", { name: theme, dark: dark });
    },
    decrementScale() {
      this.scale--;
      this.applyScale();
    },
    incrementScale() {
      this.scale++;
      this.applyScale();
    },
    applyScale() {
      document.documentElement.style.fontSize = this.scale + "px";
    },
    onRippleChange(value) {
      this.$primevue.config.ripple = value;
    },
  },
  computed: {
    layoutMode: {
      get() {
        return this.$store.getters["theme/layoutMode"];
      },
      set(value) {
        this.$store.commit("theme/layoutMode", value);
      },
    },
    rippleActive() {
      return this.$primevue.config.ripple;
    },
    // have to normalize the list, because v-for does not allow to call methods
    // inside a template
    // https://forum.vuejs.org/t/call-method-inside-a-v-for-loop/37790
    computedThemes() {
      return this.$data.data_themes.map((theme_cluster) => {
        theme_cluster.themes = theme_cluster.themes.map((theme) => {
          theme.title = theme.name
            .split("-")
            .splice(1)
            .map((word) => word[0].toUpperCase() + word.substr(1))
            .join(" ");
          var source_name = theme.source_name ? theme.source_name : theme.name;
          theme.source =
            "/images/themes/" + source_name + "." + theme.extension;
          return theme;
        });
        return theme_cluster;
      });
    },
  },
  data() {
    return {
      scale: 14,
      scales: [12, 13, 14, 15, 16],
      data_themes: [
        {
          title: "Bootstrap",
          themes: [
            { name: "bootstrap4-light-blue", extension: "svg" },
            { name: "bootstrap4-light-purple", extension: "svg" },

            { name: "bootstrap4-dark-blue", extension: "svg", dark: true },
            { name: "bootstrap4-dark-purple", extension: "svg", dark: true },
          ],
        },
        {
          title: "PrimeOne Design - 2022",
          themes: [
            { name: "lara-light-indigo", extension: "png" },
            { name: "lara-light-blue", extension: "png" },
            { name: "lara-light-purple", extension: "png" },
            { name: "lara-light-teal", extension: "png" },

            { name: "lara-dark-indigo", extension: "png", dark: true },
            { name: "lara-dark-blue", extension: "png", dark: true },
            { name: "lara-dark-purple", extension: "png", dark: true },
            { name: "lara-dark-teal", extension: "png", dark: true },
          ],
        },
        {
          title: "Material Design",
          themes: [
            { name: "md-light-indigo", extension: "svg" },
            { name: "md-light-deeppurple", extension: "svg" },

            { name: "md-dark-indigo", extension: "svg", dark: true },
            { name: "md-dark-deeppurple", extension: "svg", dark: true },
          ],
        },
        {
          title: "Material Design Compact",
          themes: [
            {
              name: "mdc-light-indigo",
              extension: "svg",
              source_name: "md-light-indigo",
            },
            {
              name: "mdc-light-deeppurple",
              extension: "svg",
              source_name: "md-light-deeppurple",
            },

            {
              name: "mdc-dark-indigo",
              extension: "svg",
              source_name: "md-dark-indigo",
              dark: true,
            },
            {
              name: "mdc-dark-deeppurple",
              extension: "svg",
              source_name: "md-dark-deeppurple",
              dark: true,
            },
          ],
        },
        {
          title: "PrimeOne Design - 2021",
          themes: [
            { name: "saga-blue", extension: "png" },
            { name: "saga-green", extension: "png" },
            { name: "saga-orange", extension: "png" },
            { name: "saga-purple", extension: "png" },

            { name: "vela-blue", extension: "png", dark: true },
            { name: "vela-green", extension: "png", dark: true },
            { name: "vela-orange", extension: "png", dark: true },
            { name: "vela-purple", extension: "png", dark: true },

            { name: "arya-blue", extension: "png", dark: true },
            { name: "arya-green", extension: "png", dark: true },
            { name: "arya-orange", extension: "png", dark: true },
            { name: "arya-purple", extension: "png", dark: true },
          ],
        },
        {
          title: "Tailwind",
          themes: [{ name: "tailwind-light", extension: "png" }],
        },
        {
          title: "Fluent UI",
          themes: [{ name: "fluent-light", extension: "png" }],
        },
      ],
    };
  },
};
</script>

<style lang="scss" scoped>
.free-themes {
  img {
    width: 50px;
    border-radius: 4px;
  }

  span {
    font-size: 0.875rem;
    margin-top: 0.25rem;
  }
}

.config-scale {
  display: flex;
  align-items: center;

  .p-button {
    margin-right: 0.5rem;
  }

  i {
    margin-right: 0.5rem;
    font-size: 0.75rem;
    color: var(--text-color-secondary);

    &.scale-active {
      font-size: 1.25rem;
      color: var(--primary-color);
    }
  }
}
</style>
