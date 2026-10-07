import { defineStore } from "pinia";
import { getDevices } from "../services/deviceService";

export const useDeviceStore = defineStore("devices", {
  state: () => ({
    devices: [],
    loading: false
  }),

  actions: {
    async loadDevices() {
      this.loading = true;

      try {
        this.devices = await getDevices();
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    }
  }
});