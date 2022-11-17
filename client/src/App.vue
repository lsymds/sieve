<script lang="ts">
import { defineComponent } from "@vue/runtime-core";
import NoOperations from "./components/NoOperations.vue";
import OperationSummary from "./components/OperationSummary.vue";
import { useOperationsStore } from "./stores/operations_store";

export default defineComponent({
  components: { NoOperations, OperationSummary },
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
  <NoOperations v-if="operations.length === 0" />
  <div :key="operation" v-else v-for="operation of operations">
    <OperationSummary :operation="operation" />
  </div>
</template>
