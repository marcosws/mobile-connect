<script setup>

import { ref, onMounted } from "vue"

const props = defineProps({
  app: {
    type: Object,
    required: true
  }
})

const emit = defineEmits([
  "close"
])

const devices = ref([])
const selectedDevice = ref("")
const installing = ref(false)

onMounted(loadDevices)

async function loadDevices() {

  try {

    const response = await fetch(
      "http://localhost:8080/devices"
    )

    devices.value = await response.json()

  } catch (err) {

    console.error(err)

  }

}

async function installApp() {

  if (!selectedDevice.value) {
    return
  }

  installing.value = true

  try {

    const response = await fetch(
      `http://localhost:8080/apps/${props.app.id}/devices/${selectedDevice.value}/install`,
      {
        method: "POST"
      }
    )

    if (!response.ok) {
      throw new Error("Installation failed")
    }

    alert("Application installed successfully.")

    emit("close")

  } catch (err) {

    console.error(err)

    alert("Installation failed.")

  } finally {

    installing.value = false

  }

}

</script>

<template>

<div class="overlay">

    <div class="dialog">

        <h2>Install Application</h2>

        <hr>

        <p>

            <strong>Application</strong>

        </p>

        <p>

            {{ app.originalName }}

        </p>

        <p class="package">

            {{ app.packageName }}

        </p>

        <label>

            Select Device

        </label>

        <select v-model="selectedDevice">

            <option value="">
                Select...
            </option>

            <option
                v-for="device in devices"
                :key="device.id"
                :value="device.id"
            >

                {{ device.manufacturer }}
                {{ device.model }}

                (Android {{ device.androidVersion }})

            </option>

        </select>

        <div class="buttons">

            <button
                class="cancel"
                @click="emit('close')"
            >

                Cancel

            </button>

            <button
                class="install"
                @click="installApp"
                :disabled="!selectedDevice || installing"
            >

                {{ installing ? "Installing..." : "Install" }}

            </button>

        </div>

    </div>

</div>

</template>

<style scoped>

.overlay{

    position:fixed;

    inset:0;

    background:rgba(0,0,0,.45);

    display:flex;

    justify-content:center;

    align-items:center;

    z-index:999;
}

.dialog{

    width:420px;

    background:white;

    border-radius:12px;

    padding:24px;

    box-shadow:0 10px 30px rgba(0,0,0,.25);
}

.dialog h2{

    margin-top:0;

    color:#1976d2;
}

.package{

    color:#666;

    font-size:13px;

    word-break:break-all;

    margin-bottom:20px;
}

select{

    width:100%;

    padding:10px;

    margin-top:10px;

    margin-bottom:20px;
}

.buttons{

    display:flex;

    justify-content:flex-end;

    gap:10px;
}

.cancel{

    background:#888;

    color:white;
}

.install{

    background:#1976d2;

    color:white;
}

button{

    border:none;

    border-radius:6px;

    padding:10px 18px;

    cursor:pointer;
}

button:disabled{

    opacity:.5;

    cursor:not-allowed;
}

</style>