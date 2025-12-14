import { defineStore } from 'pinia'
import api from '../services/api'

const ACCESS_TOKEN_KEY = 'defect_access_token'
const REFRESH_TOKEN_KEY = 'defect_refresh_token'
const USER_KEY = 'defect_user'

const loadInitialState = () => {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY) || ''
  const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY) || ''
  const userRaw = localStorage.getItem(USER_KEY)
  const user = userRaw ? JSON.parse(userRaw) : null
  return { token, refreshToken, user }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    ...loadInitialState(),
    loading: false,
    error: null,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
  },
  actions: {
    async bootstrap() {
      if (this.refreshToken && !this.token) {
        try {
          await this.refresh()
        } catch (error) {
          this.logout()
        }
      }
    },
    async login(email, password) {
      this.loading = true
      this.error = null
      try {
        const { data } = await api.login({ email, password })
        this.persistSession(data)
        return data
      } catch (error) {
        this.error =
          error?.response?.data?.message || 'Не удалось войти. Проверьте данные.'
        throw error
      } finally {
        this.loading = false
      }
    },
    async register(payload) {
      this.loading = true
      this.error = null
      try {
        const { data } = await api.register(payload)
        this.persistSession(data)
        return data
      } catch (error) {
        this.error =
          error?.response?.data?.message ||
          'Не удалось выполнить регистрацию.'
        throw error
      } finally {
        this.loading = false
      }
    },
    async refresh() {
      if (!this.refreshToken) throw new Error('no refresh token')
      const { data } = await api.refresh({ refreshToken: this.refreshToken })
      this.persistSession(data)
      return data
    },
    async logout() {
      try {
        if (this.token || this.refreshToken) {
          await api.logout({ refreshToken: this.refreshToken })
        }
      } catch (error) {
        console.warn('logout failed', error)
      }
      this.token = ''
      this.refreshToken = ''
      this.user = null
      this.error = null
      localStorage.removeItem(ACCESS_TOKEN_KEY)
      localStorage.removeItem(REFRESH_TOKEN_KEY)
      localStorage.removeItem(USER_KEY)
    },
    async changePassword(currentPassword, newPassword) {
      return api.changePassword({ currentPassword, newPassword })
    },
    persistSession(data) {
      this.token = data.accessToken
      this.refreshToken = data.refreshToken
      this.user = data.user
      localStorage.setItem(ACCESS_TOKEN_KEY, this.token)
      localStorage.setItem(REFRESH_TOKEN_KEY, this.refreshToken)
      localStorage.setItem(USER_KEY, JSON.stringify(this.user))
    },
  },
})
