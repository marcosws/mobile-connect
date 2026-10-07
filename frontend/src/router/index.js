import { createRouter, createWebHistory } from "vue-router";

import HomeView from "../views/HomeView.vue";
import DevicesView from "../views/DevicesView.vue";
import AppsView from "../views/AppsView.vue";
import SettingsView from "../views/SettingsView.vue";
import DeviceDetailView from "../views/DeviceDetailView.vue";

const routes = [
  {
    path: "/",
    name: "home",
    component: HomeView,
  },
  {
    path: "/devices",
    name: "devices",
    component: DevicesView,
  },
  {
    path: "/apps",
    name: "apps",
    component: AppsView,
  },
  {
    path: "/settings",
    name: "settings",
    component: SettingsView,
  },
  {
    path: "/devices/:id",
    name: "device-details",
    component: DeviceDetailView
  },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});