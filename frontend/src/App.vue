<template>
  <div class="layout">

    <Sidebar />

    <div class="main">

      <TopBar />

      <div class="page">
        <router-view />
      </div>

    </div>

  </div>
</template>

<script setup>
import Sidebar from "./components/Sidebar.vue"
import TopBar from "./components/TopBar.vue"
</script>

<style>
body {
  margin: 0;
  font-family: Arial, sans-serif;
}

.layout {
  display: flex;
  height: 100vh;
}

.main {
  flex: 1;

  display: flex;
  flex-direction: column;
  color:rgb(80, 80, 255); /** teste */
}

.page {
  flex: 1;
  padding: 20px;
  background: #f5f7fa;
}

html, body, #app {
  height: 100%;
  margin: 0;
}
</style>