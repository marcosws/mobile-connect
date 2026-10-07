<script setup>
import { ref } from "vue"
import { onMounted } from "vue"
import { storeToRefs } from "pinia"
import { useAppStore } from "../stores/appStore"

import AppCard from "../components/AppCard.vue"
import AppUpload from "../components/AppUpload.vue"
import InstallAppDialog from "../components/InstallAppDialog.vue"

const selectedApp = ref(null)

const store = useAppStore()

const { apps, loading } = storeToRefs(store)

onMounted(() => {
  store.loadApps()
})

async function deleteApp(id) {

  if (!confirm("Delete this application?")) {
    return
  }

  await store.deleteApp(id)

}
function installApp(app) {

    selectedApp.value = app

}

function closeDialog() {

    selectedApp.value = null

}

</script>

<template>

  <div>

    <h1>Applications</h1>

    <button @click="store.loadApps()">
      Refresh
    </button>

    <AppUpload
      @uploaded="store.loadApps()"
    />

    <div class="app-grid">

    <AppCard
        v-for="app in apps"
        :key="app.id"
        :app="app"
        @delete="deleteApp"
        @install="installApp"
    />

    </div>



  </div>

      <InstallAppDialog
        v-if="selectedApp"
        :app="selectedApp"
        @close="closeDialog"
    />

</template>

<style scoped>
.app-grid {

  display: grid;

  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));

  gap: 20px;

  margin-top: 20px;
}
</style>