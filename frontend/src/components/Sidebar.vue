<<template>
  <aside class="sidebar" :class="{ collapsed: isCollapsed }">
    <div class="sidebar-header">
      <h2 v-if="!isCollapsed">Mobile Connect</h2>
      <button class="toggle-btn" @click="isCollapsed = !isCollapsed" :aria-label="isCollapsed ? 'Expandir menu' : 'Encolher menu'">
        <span class="arrow">&lsaquo;&lsaquo;</span>
      </button>
    </div>

    <nav>
      <router-link to="/" title="Home">
        <span class="icon">          
          <img
            :src="homeLogo"
            class="menu-logo"
          /></span>
        <span class="label" v-if="!isCollapsed">Home</span>
      </router-link>
      <router-link to="/devices" title="Devices">
        <span class="icon">          <img
            :src="devicesLogo"
            class="menu-logo"
          /></span>
        <span class="label" v-if="!isCollapsed">Devices</span>
      </router-link>
      <router-link to="/apps" title="Apps">
        <span class="icon">          <img
            :src="appsLogo"
            class="menu-logo"
          /></span>
        <span class="label" v-if="!isCollapsed">Apps</span>
      </router-link>
      <router-link to="/settings" title="Settings">
        <span class="icon">          <img
            :src="settingsLogo"
            class="menu-logo"
          /></span>
        <span class="label" v-if="!isCollapsed">Settings</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup>
import { ref } from 'vue'

import homeLogo from "../assets/home32x32.png"
import devicesLogo from "../assets/device32x32.png"
import appsLogo from "../assets/apps32x32.png"
import settingsLogo from "../assets/settings32x32.png"

const isCollapsed = ref(false)
</script>

<style scoped>
.sidebar {
  width: 220px;
  height: 100vh;
  background: rgb(80, 80, 255);
  color: white;
  padding: 20px;
  box-sizing: border-box;
  transition: width 0.25s ease;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 70px;
  padding: 20px 12px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  min-height: 32px;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 1.1rem;
  white-space: nowrap;
}

.toggle-btn {
  background: rgba(255, 255, 255, 0.15);
  border: none;
  color: white;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s ease, transform 0.25s ease;
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.sidebar.collapsed .toggle-btn {
  transform: rotate(180deg);
}

.arrow {
  font-size: 1.6rem;
  line-height: 1;
}

nav {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

a {
  color: white;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 4px;
  border-radius: 6px;
  white-space: nowrap;
  overflow: hidden;
}

.icon {
  font-size: 1.2rem;
  flex-shrink: 0;
  width: 20px;
  text-align: center;
}

.label {
  overflow: hidden;
}

.router-link-active {
  font-weight: bold;
  background: rgba(255, 255, 255, 0.12);
}

.menu-logo{
  width: 80%;
  height: auto; 
}

</style>