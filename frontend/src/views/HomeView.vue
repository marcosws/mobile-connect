<template>
  <div>

    <h1>Dashboard</h1>

    <div class="cards">

      <div class="card">
        <h3>Devices Online</h3>
        <span>0</span>
      </div>

      <div class="card">
        <h3>Apps</h3>
        <span>0</span>
      </div>

      <div class="card">
        <h3>Sessions</h3>
        <span>0</span>
      </div>

    </div>

  </div>
</template>

<style scoped>
.cards {
  display: flex;
  gap: 20px;
}

.card {
  background: white;

  padding: 20px;

  border-radius: 10px;

  min-width: 180px;

  box-shadow: 0 2px 8px rgba(0,0,0,.08);
}

.card span {
  font-size: 32px;
  font-weight: bold;
}
</style>