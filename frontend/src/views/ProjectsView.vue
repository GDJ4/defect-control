<script setup>
import { computed, reactive, watch } from 'vue'
import { useProjectsStore } from '../stores/projects'
import { useAuthStore } from '../stores/auth'

const store = useProjectsStore()
const authStore = useAuthStore()

const newProject = reactive({
  name: '',
  stage: '',
  description: '',
  startDate: '',
  endDate: '',
})

const isAuthed = computed(() => authStore.isAuthenticated)

const submit = async () => {
  if (!isAuthed.value || !newProject.name) return
  await store.create({ ...newProject })
  Object.assign(newProject, {
    name: '',
    stage: '',
    description: '',
    startDate: '',
    endDate: '',
  })
}

watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      store.fetch()
    } else {
      store.$reset()
    }
  },
  { immediate: true }
)
</script>

<template>
  <section v-if="isAuthed" class="projects">
    <header>
      <p class="section-label">Проекты</p>
      <h2>Управление объектами</h2>
    </header>

    <form class="create-card" @submit.prevent="submit">
      <div class="form-row">
        <label>
          Название
          <input v-model="newProject.name" placeholder="Например, ЖК Север" required />
        </label>
        <label>
          Этап
          <input v-model="newProject.stage" placeholder="Текущий этап" />
        </label>
      </div>
      <div class="form-row">
        <label>
          Начало
          <input v-model="newProject.startDate" type="date" />
        </label>
        <label>
          Завершение
          <input v-model="newProject.endDate" type="date" />
        </label>
      </div>
      <label>
        Описание
        <textarea v-model="newProject.description" rows="2" placeholder="Краткая справка"></textarea>
      </label>
      <div class="form-actions">
        <button class="primary-btn" type="submit">Сохранить</button>
      </div>
    </form>

    <div class="cards" v-if="store.items.length">
      <article v-for="project in store.items" :key="project.id" class="project-card">
        <div class="project-card__head">
          <h3>{{ project.name }}</h3>
          <span>{{ project.stage || 'Этап не указан' }}</span>
        </div>
        <p class="project-card__meta">{{ project.description || 'Без описания' }}</p>
        <p class="project-card__dates">
          {{ project.startDate ?? '—' }} → {{ project.endDate ?? '—' }}
        </p>
      </article>
    </div>
    <p v-else class="muted">Проекты ещё не добавлены</p>
  </section>
  <p v-else class="muted">Авторизуйтесь, чтобы просматривать проекты.</p>
</template>

<style scoped>
.projects {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.create-card {
  background: rgba(15, 23, 42, 0.8);
  border-radius: 16px;
  border: 1px solid rgba(148, 163, 184, 0.15);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

label {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  font-size: 0.85rem;
  color: #94a3b8;
}

input,
textarea {
  border-radius: 10px;
  border: none;
  padding: 0.6rem 0.8rem;
  background: rgba(15, 23, 42, 0.9);
  color: #e2e8f0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.primary-btn {
  background: #3b82f6;
  border: none;
  padding: 0.65rem 1.4rem;
  border-radius: 8px;
  color: #fff;
  cursor: pointer;
}

.cards {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.project-card {
  background: rgba(15, 23, 42, 0.8);
  border-radius: 16px;
  padding: 1.5rem;
  border: 1px solid rgba(148, 163, 184, 0.15);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.project-card__head span {
  color: #94a3b8;
}

.project-card__meta {
  margin: 0;
  color: #e2e8f0;
}

.project-card__dates {
  margin: 0;
  color: #94a3b8;
}

.muted {
  color: #94a3b8;
}
</style>
