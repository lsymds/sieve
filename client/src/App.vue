<script lang="ts">
import { computed, defineComponent } from "@vue/runtime-core";
import NoOperations from "./components/NoOperations.vue";
import OperationSummary from "./components/OperationSummary.vue";
import { useOperationsStore } from "./stores/operations_store";

export default defineComponent({
  components: { NoOperations, OperationSummary },
  name: "App",
  setup() {
    const operationsStore = useOperationsStore();

    const operations = computed(() =>
      Object.keys(operationsStore.operations)
        .map((k) => operationsStore.operations[k])
        .reverse()
    );

    const hasOperations = computed(
      () => Object.keys(operationsStore.operations).length > 0
    );

    return {
      hasOperations,
      operations: operations,
    };
  },
});
</script>

<template>
  <NoOperations v-if="!hasOperations" />
  <div v-else class="flex flex-col gap-4">
    <h1 class="text-xl">Operations</h1>
    <OperationSummary
      :key="operation"
      v-for="operation of operations"
      :operation="operation"
    />
  </div>
</template>
