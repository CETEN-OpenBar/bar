<script setup>
import Button from "./components/Button.vue";
import { api } from "./lib/config/config";
import { AccountsApiFactory } from "./lib/api";
import { ref } from "vue";

let account = ref(null);

const callback = async () => {
  AccountsApiFactory({
    basePath: api(),
  })
    .getAccount({
      withCredentials: true,
    })
    .then((res) => (account.value = res.data.account))
    .catch((err) => console.log("oof"));
};
</script>

<template>
  <h1 class="text-3xl font-bold underline mb-5">Hello world!</h1>
  <Button
    class="text-white bg-blue-500 hover:bg-blue-700"
    :click="callback"
    text="Coucou"
  />
  <h1 class="text-xl" >
    <p>Carte : {{ account.card_id }}</p>
    <p>Balance : {{ account.balance }}</p>
    <p>State : {{ account.state }}</p>
  </h1>
</template>
