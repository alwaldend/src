<template>
  <div class="grid flex flex-column align-items-center">
    <div
      class="flex justify-content-center align-items-center bg-pink-500 border-circle"
      style="height: 3.2rem; width: 3.2rem"
    >
      <i class="pi pi-fw pi-exclamation-circle text-2xl text-50"></i>
    </div>
    <h1 class="font-bold text-5xl text-900 mb-2">Error Occured</h1>
    <span class="text-600">Requested resource is not available.</span>
    <!--
        https://stackoverflow.com/a/64108680
        vuejs 3 fucks up width for some reason
        (using data does not help)
        -->
    <img
      src="/layout/images/asset-error.svg"
      alt="Error"
      class="mt-5"
      :width="width"
      rel="preload"
    />
  </div>
</template>
<script>
export default {
  computed() {
    return { width: "80%" };
  },
};
</script>
