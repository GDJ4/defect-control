import { defineStore } from 'pinia'
import api from '../services/api'

export const useProjectsStore = defineStore('projects', {
  state: () => ({
    items: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetch() {
      this.loading = true
      this.error = null
      try {
        const { data } = await api.getProjects()
        this.items = data.items ?? []
      } catch (error) {
        this.error =
          error?.response?.data?.message ||
          'Не удалось загрузить список проектов'
        this.items = []
      } finally {
        this.loading = false
      }
    },
    async create(payload) {
      const { data } = await api.createProject(payload)
      this.items.unshift(data)
      return data
    },
  },
})
