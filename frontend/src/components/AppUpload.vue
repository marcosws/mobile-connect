<script setup>

import { ref } from "vue"

const emit = defineEmits([
  "uploaded"
])

const selectedFile = ref(null)
const uploading = ref(false)
const message = ref("")

function onFileSelected(event) {

  selectedFile.value = event.target.files[0]

}

async function uploadApp() {

  if (!selectedFile.value) {
    return
  }

  uploading.value = true
  message.value = ""

  const formData = new FormData()

  formData.append(
    "apk",
    selectedFile.value
  )

  try {

    const response = await fetch(
      "http://localhost:8080/apps",
      {
        method: "POST",
        body: formData
      }
    )

    if (!response.ok) {
      throw new Error("Upload failed")
    }

    message.value =
      "Application uploaded successfully."

    selectedFile.value = null

    emit("uploaded")

  } catch (err) {

    console.error(err)

    message.value =
      "Upload failed."

  } finally {

    uploading.value = false

  }

}

</script>

<template>

<div class="upload-card">

    <h3>Upload APK</h3>

    <input
        type="file"
        accept=".apk"
        @change="onFileSelected"
    />

    <p v-if="selectedFile">

        {{ selectedFile.name }}

    </p>

    <button
        @click="uploadApp"
        :disabled="!selectedFile || uploading"
    >

        {{ uploading ? "Uploading..." : "Upload" }}

    </button>

    <p
        v-if="message"
        class="message"
    >

        {{ message }}

    </p>

</div>

</template>

<style scoped>

.upload-card{

    background:white;

    border:1px solid #ddd;

    border-radius:12px;

    padding:20px;

    margin-bottom:25px;

    box-shadow:0 2px 8px rgba(0,0,0,.08);
}

.upload-card h3{

    margin-bottom:15px;

    color:#1976d2;
}

.upload-card input{

    margin-bottom:15px;
}

.upload-card button{

    margin-top:10px;

    padding:10px 18px;

    border:none;

    border-radius:6px;

    cursor:pointer;

    background:#1976d2;

    color:white;

    font-weight:600;
}

.upload-card button:hover{

    background:#1565c0;
}

.upload-card button:disabled{

    background:#999;

    cursor:not-allowed;
}

.message{

    margin-top:15px;

    color:green;
}

</style>