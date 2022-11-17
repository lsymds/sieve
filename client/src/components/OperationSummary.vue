<script lang="ts">
import { defineComponent, PropType } from "vue";
import { Operation } from "../types/Operation";

export default defineComponent({
  props: {
    operation: Object as PropType<Operation>,
  },
  setup(props) {
    return {
      operation: props.operation,
    };
  },
});
</script>

<template>
  {{ operation }}
</template>
