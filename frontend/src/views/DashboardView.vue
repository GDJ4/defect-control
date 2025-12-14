<script setup>
import { computed, watch } from 'vue'
import StatCard from '../components/ui/StatCard.vue'
import { useDefectsStore } from '../stores/defects'
import { useProjectsStore } from '../stores/projects'
import { useAuthStore } from '../stores/auth'

const defectsStore = useDefectsStore()
const projectsStore = useProjectsStore()
const authStore = useAuthStore()

const isAuthed = computed(() => authStore.isAuthenticated)

watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      defectsStore.fetch({ limit: 10 })
      projectsStore.fetch()
    } else {
      defectsStore.$reset()
      projectsStore.$reset()
    }
  },
  { immediate: true }
)

const stats = computed(() => [
  {
    label: 'Всего дефектов',
    value: defectsStore.total,
    hint: 'Всего записей',
  },
  {
    label: 'Проекты',
    value: projectsStore.items.length,
    hint: 'Активных объектов',
  },
  {
    label: 'Критичные',
    value: defectsStore.critical,
    hint: 'Priority = CRITICAL',
  },
])
</script>

<template>
  <section v-if="isAuthed" class="dashboard">
    <div class="dashboard__stats">
      <StatCard
        v-for="stat in stats"
        :key="stat.label"
        :label="stat.label"
        :value="stat.value"
        :hint="stat.hint"
      />
    </div>

    <div class="dashboard__panel">
      <header class="panel-header">
        <div>
          <p class="panel-label">Последние дефекты</p>
          <h2 class="panel-title">Мониторинг SLA</h2>
        </div>
        <RouterLink class="panel-action" to="/defects">
          Открыть список
        </RouterLink>
      </header>

      <div v-if="defectsStore.loading" class="panel-state">Загрузка...</div>
      <div v-else-if="defectsStore.error" class="panel-state panel-state--error">
        {{ defectsStore.error }}
      </div>
      <ul v-else class="defect-list">
        <li v-for="defect in defectsStore.items" :key="defect.id" class="defect-list__item">
          <div>
            <p class="defect-list__title">{{ defect.title }}</p>
            <p class="defect-list__meta">
              {{ defect.project }} · до {{ defect.dueDate || 'не задано' }}
            </p>
          </div>
          <span class="status-pill" :data-status="defect.status">
            {{ defect.status }}
          </span>
        </li>
      </ul>
    </div>
  </section>
  <p v-else class="muted">Войдите, чтобы увидеть статистику.</p>
</template>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.dashboard__stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
}

.dashboard__panel {
  background: rgba(15, 23, 42, 0.9);
  border-radius: 16px;
  padding: 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.05);
  box-shadow: 0 20px 45px rgba(15, 23, 42, 0.45);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.panel-label {
  font-size: 0.85rem;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.panel-title {
  margin: 0.25rem 0 0;
}

.panel-action {
  background: #3b82f6;
  border: none;
  color: #fff;
  padding: 0.6rem 1.2rem;
  border-radius: 999px;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.panel-action:hover {
  opacity: 0.9;
}

.panel-state {
  padding: 1rem 0;
  text-align: center;
  color: #cbd5f5;
}

.panel-state--error {
  color: #f87171;
}

.defect-list {
  list-style: none;
  margin: 1rem 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.defect-list__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.7);
  border: 1px solid rgba(148, 163, 184, 0.15);
}

.defect-list__title {
  font-weight: 600;
  margin: 0 0 0.25rem;
}

.defect-list__meta {
  margin: 0;
  font-size: 0.85rem;
  color: #94a3b8;
}

.status-pill {
  padding: 0.35rem 0.75rem;
  border-radius: 999px;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.status-pill[data-status='IN_PROGRESS'] {
  background: rgba(59, 130, 246, 0.2);
  color: #93c5fd;
}

.status-pill[data-status='NEW'] {
  background: rgba(248, 113, 113, 0.15);
  color: #fecaca;
}
</style>
