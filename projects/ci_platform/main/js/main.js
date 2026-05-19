import "primevue/resources/primevue.min.css";
import "primeflex/primeflex.css";
import "primeicons/primeicons.css";
import "prismjs/themes/prism-coy.css";

import "./assets/styles/layout.scss";
import "./assets/demo/flags/flags.css";

import { createApp } from "vue";
import { VueCookies } from "./lib/main/Cookies";

import PrimeVue from "primevue/config";
import BadgeDirective from "primevue/badgedirective";
import ConfirmationService from "primevue/confirmationservice";
import Ripple from "primevue/ripple";
import StyleClass from "primevue/styleclass";
import Tooltip from "primevue/tooltip";
import ToastService from "primevue/toastservice";

import Router from "./lib/main/Router";
import AppWrapper from "./components/app/AppWrapper.vue";
import CodeHighlight from "./lib/app/AppCodeHighlight";
import UniversalComponents from "./lib/main/UniversalComponents";
import VuexStore from "./lib/store";

Router.beforeEach(function (to, from, next) {
  window.scrollTo(0, 0);
  next();
});

const app = createApp(AppWrapper);
// initialise user store from local storage
VuexStore.dispatch("user/initializeStore");

app.use(PrimeVue, { ripple: true, inputStyle: "outlined" });
app.use(ConfirmationService);
app.use(ToastService);
app.use(Router);
app.use(VuexStore);
app.use(VueCookies);

const directives = {
  tooltip: Tooltip,
  ripple: Ripple,
  code: CodeHighlight,
  badge: BadgeDirective,
  styleclass: StyleClass,
};

for (const [key, value] of Object.entries(UniversalComponents)) {
  app.component(key, value);
}

for (const [key, value] of Object.entries(directives)) {
  app.directive(key, value);
}

app.mount("#app");
