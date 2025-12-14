import axios from 'axios'

const baseURL =
  import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

const client = axios.create({
  baseURL,
  timeout: 5000,
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('defect_access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      console.error('API error', error.response.status, error.response.data)
    } else {
      console.error('API error', error.message)
    }
    return Promise.reject(error)
  }
)

const api = {
  register(payload) {
    return client.post('/auth/register', payload)
  },
  login(payload) {
    return client.post('/auth/login', payload)
  },
  refresh(payload) {
    return client.post('/auth/refresh', payload)
  },
  logout(payload) {
    return client.post('/auth/logout', payload)
  },
  changePassword(payload) {
    return client.post('/auth/password', payload)
  },
  getProjects(params = {}) {
    return client.get('/projects', { params })
  },
  createProject(payload) {
    return client.post('/projects', payload)
  },
  getDefects(params = {}) {
    return client.get('/defects', { params })
  },
  createDefect(payload) {
    return client.post('/defects', payload)
  },
  getDefect(id) {
    return client.get(`/defects/${id}`)
  },
  getDefectComments(id) {
    return client.get(`/defects/${id}/comments`)
  },
  addDefectComment(id, payload) {
    return client.post(`/defects/${id}/comments`, payload)
  },
  addDefectAttachment(id, file) {
    const formData = new FormData()
    formData.append('file', file)
    return client.post(`/defects/${id}/attachments`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export default api
