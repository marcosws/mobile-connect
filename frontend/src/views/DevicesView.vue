<script setup>
import { onMounted } from "vue";
import { storeToRefs } from "pinia";
import androidLogo from "../assets/android32x32.png"

import { useDeviceStore } from "../stores/deviceStore";

const store = useDeviceStore();

const { devices, loading } = storeToRefs(store);

onMounted(() => {
  store.loadDevices();
});
</script>

<template>
  <div>
    <h1>Devices</h1>
    <button @click="store.loadDevices()">
      Refresh
    </button>

    <p v-if="loading">
      Loading...
    </p>

    <div
      v-else
      class="device-grid"
    >

      <div
        v-for="device in devices"
        :key="device.id"
        class="device-card"
      >
      <div class="phone-frame">
        <div class="android-logo">
          <img
            :src="androidLogo"
            class="android-logo"
          />
        </div>
      </div>

        <h3>
          {{ device.manufacturer }}
        </h3>

        <p class="model">
          {{ device.model }}
        </p>

        <p class="device-id">
          {{ device.id }}
        </p>

        <p>
          Android {{ device.androidVersion }}
        </p>

        <p>
          SDK {{ device.sdk }}
        </p>

        <p
          :class="
            device.status === 'device'
            ? 'status-online'
            : 'status-offline'
          "
        >
          {{
            device.status === 'device'
            ? '🟢 Online'
            : '🔴 Offline'
          }}
        </p>
        <router-link
          :to="`/devices/${device.id}`"
          target="_blank"
          class="open-button"
        >
          Open
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>

button{
  background-color: #7c7c7c;
  color: #ffffff;
  border: none;
  height: 30px;
  width: 120px;
  border-radius: 8px;
}

.device-grid {
  margin-top: 20px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 20px;
}

.device-card {
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 12px;
  background: white;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: 0.2s;
}

.device-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.android-logo {
  font-size: 72px;
  margin-bottom: 10px;
  object-fit: contain;
}

.device-card h3 {
  font-size: 15px;
  margin-bottom: 4px;
}

.model {
  font-size: 14px;
  font-weight: 600;
  color: #1976d2;
}

.device-card p {
  font-size: 13px;
  margin: 4px 0;
}

.device-id {
  font-size: 11px;
  color: #666;
  word-break: break-all;
}

.status-online {
  color: green;
  font-weight: bold;
}

.status-offline {
  color: red;
  font-weight: bold;
}

.open-button {
  display: inline-block;
  margin-top: 10px;
  padding: 8px 16px;
  background: #1976d2;
  color: white;
  text-decoration: none;
  border-radius: 6px;
}

.open-button:hover {
  background: #1565c0;
}

.phone-frame {

    width: 110px;
    height: 180px;

    margin: 0 auto 15px;

    border: 3px solid #555;
    border-radius: 18px;

    background: #f8f8f8;

    display: flex;
    align-items: center;
    justify-content: center;

    position: relative;
}

.phone-frame::before {

    content: "";

    position: absolute;

    top: 8px;

    width: 40px;
    height: 4px;

    background: #555;

    border-radius: 4px;
}

.phone-frame::after {

    content: "";

    position: absolute;

    bottom: 8px;

    width: 12px;
    height: 12px;

    border: 2px solid #555;

    border-radius: 50%;
}

</style>