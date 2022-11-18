<script lang="ts">
import { defineComponent, PropType } from "vue";

type BadgeType = "SUCCESS" | "PLAIN";

export default defineComponent({
  props: {
    text: String,
    type: Object as PropType<BadgeType>,
  },
  setup(props) {},
});
</script>

<template>
  <div
    class="text-xs px-2 py-1 rounded-full"
    :class="{
      'bg-green-100': type === 'SUCCESS',
      'bg-gray-50': type === 'PLAIN',
    }"
  >
    {{ text }}
  </div>
</template>
