<script setup>

const props = defineProps({
  app: {
    type: Object,
    required: true
  }
})

const emit = defineEmits([
  "delete",
  "install"
])

import androidLogo from "../assets/package32x32.png"

function formatSize(bytes) {

  return (bytes / 1024 / 1024).toFixed(2) + " MB"

}

</script>

<template>

  <div class="app-card">

    <div class="phone-frame">

      <img
        :src="androidLogo"
        class="android-logo"
        alt="Android"
      />

    </div>

    <h3>
      {{ app.originalName }}
    </h3>

    <p class="package-name">
      {{ app.packageName }}
    </p>

    <p>
      <strong>Version:</strong>
      {{ app.versionName }}
      ({{ app.versionCode }})
    </p>

    <p>
      <strong>SDK:</strong>
      {{ app.minSdk }} → {{ app.targetSdk }}
    </p>

    <p>
      <strong>Size:</strong>
      {{ formatSize(app.size) }}
    </p>

    <div class="actions">

    <button
    @click="emit('install', app)"
    >
    Install
    </button>

    <button
    class="danger"
    @click="emit('delete', app.id)"
    >
    Delete
    </button>

    </div>

  </div>

</template>

<style scoped>

.app-card {

  padding: 20px;

  background: white;

  border: 1px solid #ddd;
  border-radius: 12px;

  text-align: center;

  box-shadow: 0 2px 8px rgba(0,0,0,.08);

  transition: .2s;
}

.app-card:hover {

  transform: translateY(-3px);

  box-shadow: 0 6px 16px rgba(0,0,0,.12);

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

.android-logo {

  width: 32px;
  height: 32px;

  object-fit: contain;
}

h3 {

  margin-bottom: 8px;

  font-size: 16px;

  color: #1976d2;
}

.package-name {

  font-size: 12px;

  color: #666;

  word-break: break-word;

  margin-bottom: 12px;
}

p {

  margin: 6px 0;

  font-size: 13px;

  color: #444;
}

.actions {

  display: flex;

  justify-content: center;

  gap: 10px;

  margin-top: 18px;
}

.actions button {

  padding: 8px 14px;

  border: none;

  border-radius: 6px;

  cursor: pointer;

  font-weight: 600;

  transition: .2s;
}

.actions button:first-child {

  background: #1976d2;

  color: white;
}

.actions button:first-child:hover {

  background: #125ca1;
}

.danger {

  background: #d32f2f;

  color: white;
}

.danger:hover {

  background: #b71c1c;
}

</style>