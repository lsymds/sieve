<script lang="ts">
import { defineComponent } from "@vue/runtime-core";
import { useOperationsStore } from "./stores/operations_store";

export default defineComponent({
  name: "App",
  setup() {
    const operationsStore = useOperationsStore();

    return {
      operations: operationsStore.operations,
    };
  },
});
</script>

<template>
  <div :key="operation" v-for="operation of operations">
    {{ operation }}
  </div>
</template>
