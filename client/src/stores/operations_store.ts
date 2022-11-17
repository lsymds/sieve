import { defineStore } from "pinia";
import { Operation } from "../types/Operation";

export const useOperationsStore = defineStore("operations", {
  state: () => ({
    operations: [] as Operation[],
  }),
});
