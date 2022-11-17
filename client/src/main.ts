import { createApp } from "vue";
import { createPinia } from "pinia";
import "./style.css";
import App from "./App.vue";
import { useOperationsStore } from "./stores/operations_store";

const pinia = createPinia();
const app = createApp(App);

let websocket = new WebSocket("ws://localhost:8080/_/ws");

websocket.onopen = function () {
  console.log("Websocket connection established");
};

websocket.onmessage = function (evt: MessageEvent<any>) {
  console.log(evt);
  const { type, data } = JSON.parse(evt.data) as {
    type: "operation" | "ping";
    data: any;
  };

  if (type === "operation") {
    const operationsStore = useOperationsStore();
    operationsStore.operations.push(data);
  }
};

app.use(pinia);
app.mount("#app");
