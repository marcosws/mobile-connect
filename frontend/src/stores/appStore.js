import { defineStore } from "pinia"

export const useAppStore = defineStore("apps", {

  state: () => ({
    apps: [],
    loading: false
  }),

  actions: {

    async loadApps() {

      this.loading = true

      try {

        const response = await fetch(
          "http://localhost:8080/apps"
        )

        this.apps = await response.json()

      } catch (err) {

        console.error(err)

      } finally {

        this.loading = false

      }

    },
    async deleteApp(id) {

  try {

    await fetch(
      `http://localhost:8080/apps/${id}`,
      {
        method: "DELETE"
      }
    )

    await this.loadApps()

  } catch (err) {

    console.error(err)

  }

}

  }

})


