<template>
  <div class="grid flex-column align-items-center">
    <span class="text-blue-500 font-bold text-3xl">404</span>
    <h1 class="text-900 font-bold text-3xl lg:text-5xl mb-2">
      Looks like you are lost
    </h1>
    <span class="text-gray-600">Requested resource is not available.</span>
    <div
      class="col-12 flex align-items-center py-5 mt-6 border-300 border-bottom-1"
    >
      <div
        class="flex justify-content-center align-items-center bg-cyan-400 border-round"
        style="height: 3.5rem; width: 3.5rem"
      >
        <i class="pi pi-fw pi-table text-50 text-2xl"></i>
      </div>
      <div class="ml-4">
        <router-link to="/">
          <p class="text-900 lg:text-xl font-medium mb-0">
            Frequently Asked Questions
          </p>
        </router-link>
        <span class="text-gray-600 lg:text-xl"
          >Ultricies mi quis hendrerit dolor.</span
        >
      </div>
    </div>
    <div
      class="col-12 flex align-items-center py-5 border-300 border-bottom-1"
    >
      <div
        class="flex justify-content-center align-items-center bg-orange-400 border-round"
        style="height: 3.5rem; width: 3.5rem"
      >
        <i class="pi pi-fw pi-question-circle text-50 text-2xl"></i>
      </div>
      <div class="ml-4">
        <router-link to="/">
          <p class="text-900 lg:text-xl font-medium mb-0">Solution Center</p>
        </router-link>
        <span class="text-gray-600 lg:text-xl"
          >Phasellus faucibus scelerisque eleifend.</span
        >
      </div>
    </div>
    <div
      class="col-12 flex align-items-center py-5 border-300 border-bottom-1"
    >
      <div
        class="flex justify-content-center align-items-center bg-indigo-400 border-round"
        style="height: 3.5rem; width: 3.5rem"
      >
        <i class="pi pi-fw pi-unlock text-50 text-2xl"></i>
      </div>
      <div class="ml-4">
        <router-link to="/">
          <p class="text-900 lg:text-xl font-medium mb-0">
            Permission Manager
          </p>
        </router-link>
        <span class="text-gray-600 lg:text-xl"
          >Accumsan in nisl nisi scelerisque</span
        >
      </div>
    </div>
  </div>
</template>
