<script>
// Params for button : size, color, text & callback
export default {
  props: {
    text: {
      type: String,
      default: "Button",
    },
    click: {
      type: Function,
      default: () => {},
    },
  },
};
</script>

<template>
  <a
    class="inline-block px-4 py-2 rounded hover:cursor-pointer"
    type="button"
    @click="click()"
  >
    {{ text }}
  </a>
</template>
