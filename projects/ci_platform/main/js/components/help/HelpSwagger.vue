<template>
  <div class="grid">
    <div class="col-12">
      <div class="card">
        <div class="card" v-show="showSkeleton">
          <div class="border-round border-1 surface-border p-4">
            <div class="flex mb-3">
              <Skeleton shape="circle" size="4rem" class="mr-2"></Skeleton>
              <div>
                <Skeleton width="10rem" class="mb-2"></Skeleton>
                <Skeleton width="5rem" class="mb-2"></Skeleton>
                <Skeleton height=".5rem"></Skeleton>
              </div>
            </div>
            <Skeleton width="100%" height="150px"></Skeleton>
            <div class="flex justify-content-between mt-3">
              <Skeleton width="4rem" height="2rem"></Skeleton>
              <Skeleton width="4rem" height="2rem"></Skeleton>
            </div>
          </div>
        </div>
        <div @load="showSkeleton = false" ref="swagger"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@import "swagger-ui/dist/swagger-ui.css";
</style>

<script>
import SwaggerUI from "swagger-ui";
import axios from "axios";
import Error from "../../lib/main/Error";
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      swagger_src: "/api/v1/swagger/swagger.json",
      showSkeleton: true,
      swagger: null,
    };
  },
  mounted() {
    // https://swagger.io/docs/open-source-tools/swagger-ui/usage/installation/
    axios
      .get(this.$data.swagger_src)
      .then((response) => {
        this.swagger = SwaggerUI({
          domNode: this.$refs.swagger,
          spec: response.data,
        });
        this.preauthorize();
        this.$data.showSkeleton = false;
      })
      .catch((error) => {
        Error(this.$toast, error);
      });
  },
  watch: {
    token() {
      this.preauthorize();
    },
  },
  computed: {
    ...mapGetters("user", {
      token: "token",
    }),
    outerHeight() {
      return window.outerHeight;
    },
  },
  methods: {
    preauthorize() {
      if (!this.$data.swagger) {
        return;
      }
      this.$data.swagger.preauthorizeApiKey(
        "ApiKeyAuth",
        `Bearer ${this.token}`
      );
    },
  },
};
</script>
