import { defineStore } from "pinia";
import { Operation } from "../types/Operation";

export const useOperationsStore = defineStore("operations", {
  state: () => ({
    operations: [] as Operation[],
  }),
  actions: {
    async getOperation(id: string): Promise<void> {
      const response = await fetch(`http://localhost:8080/_/operations/${id}`);
      this.operations.push(await response.json());
    },
  },
});
