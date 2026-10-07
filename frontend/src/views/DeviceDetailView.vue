<script setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import DeviceScreen from "../components/DeviceScreen.vue";
import DeviceControls from "../components/DeviceControls.vue"

const route = useRoute();
const deviceId = route.params.id;
const device = ref(null);
const selectedFile = ref(null)
const installMessage = ref("")
const packageName = ref("")
const packages = ref([])
const selectedPackage = ref("")

const props = defineProps({
  deviceId: {
    type: String,
    required: true
  }
})

onMounted(async () => {

  const response = await fetch(
    `http://localhost:8080/devices/${deviceId}`
  );

  device.value = await response.json();

  await loadPackages()

});

function onFileSelected(event) {
  selectedFile.value =
    event.target.files[0]
}

async function installApk() {

  if (!selectedFile.value) {
    return
  }

  const formData = new FormData()

  formData.append(
    "apk",
    selectedFile.value
  )

  try {

    await fetch(
      `http://localhost:8080/devices/${deviceId}/install`,
      {
        method: "POST",
        body: formData
      }
    )

    installMessage.value =
      "APK installed successfully"

  } catch (err) {

    installMessage.value =
      "Installation failed"

    console.error(err)
  }
}

async function uninstallApp() {

  if (!selectedPackage.value) {
    return
  }

  try {

    await fetch(
      `http://localhost:8080/devices/${deviceId}/uninstall`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          package: selectedPackage.value
        })
      }
    )

    alert("Application removed")

    await loadPackages()

  } catch (err) {

    console.error(err)

    alert("Uninstall failed")
  }
}

async function launchApp() {

  if (!selectedPackage.value) {
    return
  }

  try {

    await fetch(
      `http://localhost:8080/devices/${deviceId}/launch`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          package: selectedPackage.value
        })
      }
    )

  } catch (err) {

    console.error(err)

  }
}

async function stopApp() {

  if (!selectedPackage.value) {
    return
  }

  try {

    await fetch(
      `http://localhost:8080/devices/${deviceId}/stop`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          package: selectedPackage.value
        })
      }
    )

  } catch (err) {

    console.error(err)

  }
}

async function loadPackages() {

  try {

    const response = await fetch(
      `http://localhost:8080/devices/${deviceId}/packages`
    )

    packages.value =
      await response.json()

  } catch (err) {

    console.error(err)

  }
}

</script>

<template>
  <div class="page">

    <!-- ESQUERDA -->
    <div class="left-panel">

      <h2>Remote Session</h2>

     <!-- <p><strong>Device ID: {{ deviceId }} </strong></p> -->

      <DeviceScreen :device-id="deviceId" />

      <DeviceControls :device-id="deviceId" />

    </div>

    <!-- DIREITA -->
    <div class="right-panel">

      <h2>Device Details</h2>

      <div v-if="device">

      <h3>Information</h3>

        <table id="device-details">
          <tbody>
            <tr><th>ID</th><td>{{ device.id }}</td></tr>
            <tr><th>Manufacturer</th><td>{{ device.manufacturer }}</td></tr>
            <tr><th>Model</th><td>{{ device.model }}</td></tr>
            <tr><th>Android</th><td>{{ device.androidVersion }}</td></tr>
            <tr><th>SDK</th><td>{{ device.sdk }}</td></tr>
          </tbody>
        </table>

      </div>

      <h3>APK Management</h3>

      <input
        type="file"
        accept=".apk"
        @change="onFileSelected"
      />

      <button
        @click="installApk"
        :disabled="!selectedFile"
      >
        Install APK
      </button>

      <p v-if="selectedFile">
        {{ selectedFile.name }}
      </p>

      <p v-if="installMessage">
        {{ installMessage }}
      </p>

      <hr>

      <!-- Uninstall -->

      <input
        v-model="packageName"
        placeholder="com.android.chrome"
      />

      <h3>Installed Applications</h3>

      <select v-model="selectedPackage">

        <option value="">
          Select application
        </option>

        <option
          v-for="pkg in packages"
          :key="pkg"
          :value="pkg"
        >
          {{ pkg }}
        </option>

      </select>

      <div class="app-actions">

        <button @click="launchApp">
          Launch App
        </button>

        <button @click="stopApp">
          Stop App
        </button>

        <button @click="uninstallApp">
          Uninstall APK
        </button>

      </div>

    </div>

  </div>
</template>
<style scoped>

button {
  background-color: #7c7c7c;
  color: #ffffff;
  border: none;
  height: 30px;
  width: 120px;
  border-radius: 8px;
}


#device-details {
  border-collapse: collapse;
  margin: 0 auto; /* centraliza horizontalmente */
  color: #807e7e;
  border: 1px solid #ccc;
  width: 60%; /* necessário para o margin: auto funcionar */
}

  th, td {
    border: 1px solid #ccc;
    padding: 8px 12px;
    text-align: left;
  }

th {
    background-color: #f0f0f0;
  }

h2 {
    text-align: center;
    margin-bottom: 20px;
  }

  p {
    text-align: center;
    font-size: 10px;
    margin-bottom: 20px;
    color: #666;
  }

  .page {
  display: flex;
  gap: 20px;
  align-items: flex-start;
  padding: 20px;
}

.left-panel {
  flex: 2;
}

.right-panel {
  flex: 1;
  position: sticky;
  top: 20px;
}

.apk-install {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ddd;
}

.apk-install button {
  margin-top: 10px;
}

.app-actions {
  margin-top: 10px;
  display: flex;
  gap: 10px;
}

.app-actions button {
  padding: 8px 12px;
}

</style>