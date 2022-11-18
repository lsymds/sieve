import { defineStore } from "pinia";
import { Operation } from "../types/Operation";

export const useOperationsStore = defineStore("operations", {
  state: () => ({
    operations: {} as Record<string, Operation>,
  }),
  actions: {
    async getOperation(id: string): Promise<void> {
      const response = await fetch(`http://localhost:8080/_/operations/${id}`);
      const operation = (await response.json()) as Operation;

      this.operations[operation.id] = operation;
    },
  },
});
