import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080"
});

export async function getDevices() {
  const response = await api.get("/devices");
  return response.data;
}