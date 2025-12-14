<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
authStore.bootstrap()

const isRegister = ref(false)
const email = ref('manager@systemacontrola.ru')
const password = ref('password')
const fullName = ref('')
const role = ref('engineer')
const roleOptions = [
  { value: 'engineer', label: 'Инженер' },
  { value: 'manager', label: 'Менеджер' },
  { value: 'observer', label: 'Наблюдатель' },
]
const currentPassword = ref('')
const newPassword = ref('')
const changeMessage = ref('')

const submit = async () => {
  if (!email.value || !password.value) return
  if (isRegister.value && !fullName.value) return
  try {
    if (isRegister.value) {
      await authStore.register({
        email: email.value,
        fullName: fullName.value,
        password: password.value,
        role: role.value,
      })
    } else {
      await authStore.login(email.value, password.value)
    }
  } catch (error) {
    console.warn('auth failed', error)
  }
}

const toggleMode = () => {
  isRegister.value = !isRegister.value
  changeMessage.value = ''
  authStore.error = null
}

const changePassword = async () => {
  if (!currentPassword.value || !newPassword.value) return
  try {
    await authStore.changePassword(currentPassword.value, newPassword.value)
    changeMessage.value = 'Пароль обновлён'
    currentPassword.value = ''
    newPassword.value = ''
  } catch (error) {
    changeMessage.value =
      error?.response?.data?.message || 'Не удалось сменить пароль'
  }
}
</script>

<template>
  <div class="auth-panel">
    <form v-if="!authStore.isAuthenticated" class="auth-form" @submit.prevent="submit">
      <input v-if="isRegister" v-model="fullName" type="text" placeholder="ФИО" required />
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Пароль" required />
      <select v-if="isRegister" v-model="role">
        <option v-for="option in roleOptions" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>
      <button class="primary-btn" type="submit" :disabled="authStore.loading">
        {{
          authStore.loading
            ? isRegister
              ? 'Регистрация...'
              : 'Вход...'
            : isRegister
              ? 'Зарегистрироваться'
              : 'Войти'
        }}
      </button>
      <button class="link-btn" type="button" @click="toggleMode">
        {{ isRegister ? 'У меня уже есть аккаунт' : 'Нет аккаунта? Регистрация' }}
      </button>
      <p v-if="authStore.error" class="auth-error">{{ authStore.error }}</p>
    </form>

    <div v-else class="auth-info">
      <div>
        <p>{{ authStore.user.fullName }} ({{ authStore.user.role }})</p>
        <form class="change-form" @submit.prevent="changePassword">
          <input
            v-model="currentPassword"
            type="password"
            placeholder="Текущий пароль"
            required
          />
          <input
            v-model="newPassword"
            type="password"
            placeholder="Новый пароль"
            required
          />
          <button class="secondary-btn" type="submit">Сменить пароль</button>
        </form>
        <p v-if="changeMessage" class="change-message">{{ changeMessage }}</p>
      </div>
      <button class="secondary-btn" type="button" @click="authStore.logout">
        Выйти
      </button>
    </div>
  </div>
</template>

<style scoped>
.auth-panel {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  align-items: flex-end;
}

.auth-form {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
}

.auth-form input,
.auth-form select,
.change-form input {
  border-radius: 8px;
  border: none;
  padding: 0.5rem 0.75rem;
  background: rgba(15, 23, 42, 0.8);
  color: #e2e8f0;
}

.auth-form select {
  min-width: 150px;
}

.change-form {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.auth-error,
.change-message {
  margin: 0.25rem 0 0;
  font-size: 0.8rem;
  color: #f87171;
}

.auth-info {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  color: #94a3b8;
}

.change-message {
  color: #fbbf24;
}

.primary-btn,
.secondary-btn {
  border: none;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  cursor: pointer;
}

.primary-btn {
  background: #3b82f6;
  color: #fff;
}

.secondary-btn {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #e2e8f0;
}

.link-btn {
  border: none;
  background: transparent;
  color: #93c5fd;
  cursor: pointer;
  text-decoration: underline;
}
</style>
