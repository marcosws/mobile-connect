<script setup>

import { ref } from "vue"

const textInput = ref("")

const props = defineProps({
  deviceId: {
    type: String,
    required: true
  }
})

function sendKey(command) {
  fetch(`http://localhost:8080/devices/${props.deviceId}/shell`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      command
    })
  }).catch(err => {
    console.error("Shell command failed:", err)
  })
}

function sendText() {

  if (!textInput.value.trim()) {
    return
  }

  fetch(
    `http://localhost:8080/devices/${props.deviceId}/shell`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        command: `input text "${textInput.value.replaceAll(" ", "%s")}"`
      })
    }
  )

  textInput.value = ""
}

</script>

<template>
  <div class="remote-session">

    <DeviceScreen :device-id="deviceId" />

    <div class="controls-panel">

      <div class="navigation-buttons">
        <button @click="sendKey('input keyevent 4')">
          ⟵ Back
        </button>

        <button @click="sendKey('input keyevent 3')">
          ⌂ Home
        </button>

        <button @click="sendKey('input keyevent 187')">
          ☰ Recent
        </button>
      </div>

      <div class="text-input-section">
        <input
          v-model="textInput"
          placeholder="Type text..."
          @keyup.enter="sendText"
        />

        <button @click="sendText">
          Send
        </button>
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

.remote-session {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.controls-panel {

  width: 100%;
  max-width: 400px;

  margin-top: 10px;

  border: 1px solid #ccc;
  border-radius: 8px;

  overflow: hidden;

  background: #f8f8f8;
}


.navigation-buttons {

  display: flex;
  justify-content: center;

  gap: 10px;

  padding: 12px;

  border-bottom: 1px solid #ddd;
}

.text-input-section {

  display: flex;

  gap: 10px;

  padding: 12px;
}

.text-input-section input {

  flex: 1;

  padding: 8px;

  border: 1px solid #ccc;
  border-radius: 4px;
}
</style>