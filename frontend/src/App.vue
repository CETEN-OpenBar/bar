<script setup>
import Button from "./components/Button.vue";
import { config } from "./lib/config/config.js";

const callback = async () => {
  console.log(await config());
};
</script>

<template>
  <h1 class="text-3xl font-bold underline mb-5">
    Hello world!
  </h1>
  <Button class="text-white bg-blue-500 hover:bg-blue-700" :click="callback" text="Coucou" />
</template>
