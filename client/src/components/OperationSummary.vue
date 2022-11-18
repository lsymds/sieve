<script lang="ts">
import { computed, defineComponent, PropType } from "vue";
import { format } from "date-fns";
import { Operation } from "../types/Operation";
import Badge from "./Badge.vue";

export default defineComponent({
  components: { Badge },
  props: {
    operation: {
      type: Object as PropType<Operation>,
      required: true,
    },
  },
  setup(props) {
    const createdAtReadable = computed(() =>
      format(new Date(props.operation.createdAt), "HH:mm:ss")
    );

    return {
      operation: props.operation,
      createdAtReadable,
    };
  },
});
</script>

<template>
  <div class="w-full bg-white border p-4 rounded">
    <div class="flex flex-row justify-between">
      <div class="flex flex-row gap-2">
        <strong>{{ operation.request.fullUrl }}</strong>
        <Badge
          :type="operation.response ? 'SUCCESS' : 'PLAIN'"
          :text="operation.response ? 'COMPLETED' : 'IN PROGRESS'"
        />
      </div>
      <span>{{ createdAtReadable }}</span>
    </div>
  </div>
</template>
