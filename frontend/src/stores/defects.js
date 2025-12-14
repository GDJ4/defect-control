import { defineStore } from 'pinia'
import api from '../services/api'

export const useDefectsStore = defineStore('defects', {
  state: () => ({
    items: [],
    current: null,
    comments: [],
    loading: false,
    error: null,
  }),
  getters: {
    total: (state) => state.items.length,
    critical: (state) =>
      state.items.filter((item) => item.priority === 'CRITICAL').length,
    inProgress: (state) =>
      state.items.filter((item) => item.status === 'IN_PROGRESS').length,
  },
  actions: {
    async fetch(filter = {}) {
      this.loading = true
      this.error = null
      try {
        const { data } = await api.getDefects(filter)
        this.items = data?.items ?? []
      } catch (error) {
        this.error =
          error?.response?.data?.message ||
          'Не удалось загрузить список дефектов'
        this.items = []
      } finally {
        this.loading = false
      }
    },
    async create(payload) {
      const { data } = await api.createDefect(payload)
      this.items.unshift(data)
      return data
    },
    async fetchOne(id) {
      const { data } = await api.getDefect(id)
      this.current = data
      this.comments = data.comments ?? []
      return data
    },
    async fetchComments(id) {
      const { data } = await api.getDefectComments(id)
      this.comments = data.items ?? []
      return this.comments
    },
    async addComment(id, payload) {
      const { data } = await api.addDefectComment(id, payload)
      this.comments.push(data)
      return data
    },
    async addAttachment(id, file) {
      const { data } = await api.addDefectAttachment(id, file)
      if (this.current) {
        this.current.attachments = [data, ...(this.current.attachments ?? [])]
      }
      return data
    },
  },
})
