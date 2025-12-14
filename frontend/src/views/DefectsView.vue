<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useDefectsStore } from '../stores/defects'
import { useProjectsStore } from '../stores/projects'
import { useAuthStore } from '../stores/auth'

const DEFAULT_ENGINEER = '22222222-2222-2222-2222-222222222222'

const defectsStore = useDefectsStore()
const projectsStore = useProjectsStore()
const authStore = useAuthStore()

const filter = reactive({
  status: '',
  priority: '',
  projectId: '',
})

const newDefect = reactive({
  projectId: '',
  title: '',
  description: '',
  priority: 'MEDIUM',
  severity: 'MAJOR',
  assigneeId: DEFAULT_ENGINEER,
  dueDate: '',
})

const statuses = ['NEW', 'IN_PROGRESS', 'IN_REVIEW', 'CLOSED', 'CANCELED']
const priorities = ['LOW', 'MEDIUM', 'HIGH', 'CRITICAL']

const selectedId = ref(null)
const detailLoading = ref(false)
const commentBody = ref('')
const attachmentInput = ref(null)

const isAuthed = computed(() => authStore.isAuthenticated)

const selectedDefect = computed(() => {
  if (!selectedId.value) return null
  if (defectsStore.current && defectsStore.current.id === selectedId.value) {
    return defectsStore.current
  }
  return null
})

const applyFilters = () => {
  defectsStore.fetch({
    status: filter.status || undefined,
    priority: filter.priority || undefined,
    projectId: filter.projectId || undefined,
  })
}

const submitDefect = async () => {
  if (!isAuthed.value || !newDefect.projectId || !newDefect.title) return
  await defectsStore.create({ ...newDefect })
  Object.assign(newDefect, {
    projectId: '',
    title: '',
    description: '',
    priority: 'MEDIUM',
    severity: 'MAJOR',
    assigneeId: DEFAULT_ENGINEER,
    dueDate: '',
  })
}

const openDetail = async (item) => {
  if (!isAuthed.value) return
  selectedId.value = item.id
  detailLoading.value = true
  await defectsStore.fetchOne(item.id)
  detailLoading.value = false
}

const submitComment = async () => {
  if (!isAuthed.value || !selectedId.value || !commentBody.value.trim()) return
  await defectsStore.addComment(selectedId.value, {
    body: commentBody.value,
  })
  commentBody.value = ''
}

const uploadAttachment = async () => {
  if (
    !isAuthed.value ||
    !attachmentInput.value?.files?.length ||
    !selectedId.value
  )
    return
  const file = attachmentInput.value.files[0]
  await defectsStore.addAttachment(selectedId.value, file)
  attachmentInput.value.value = ''
}

watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      projectsStore.fetch()
      defectsStore.fetch()
    } else {
      defectsStore.$reset()
      projectsStore.$reset()
      selectedId.value = null
    }
  },
  { immediate: true }
)
</script>

<template>
  <section v-if="isAuthed" class="defects-grid">
    <div class="defects-main">
      <header class="defects__header">
        <div>
          <p class="section-label">Список дефектов</p>
          <h2>Работа с заявками</h2>
        </div>
      </header>

      <form class="create-card" @submit.prevent="submitDefect">
        <div class="form-row">
          <label>
            Проект
            <select v-model="newDefect.projectId" required>
              <option value="">Выберите проект</option>
              <option v-for="project in projectsStore.items" :key="project.id" :value="project.id">
                {{ project.name }}
              </option>
            </select>
          </label>
          <label>
            Заголовок
            <input v-model="newDefect.title" placeholder="Опишите проблему" required />
          </label>
        </div>
        <div class="form-row">
          <label>
            Приоритет
            <select v-model="newDefect.priority">
              <option v-for="priority in priorities" :key="priority" :value="priority">
                {{ priority }}
              </option>
            </select>
          </label>
          <label>
            Серьёзность
            <select v-model="newDefect.severity">
              <option v-for="priority in priorities" :key="priority" :value="priority">
                {{ priority }}
              </option>
            </select>
          </label>
          <label>
            Срок
            <input v-model="newDefect.dueDate" type="date" />
          </label>
        </div>
        <label>
          Описание
          <textarea v-model="newDefect.description" rows="2" placeholder="Детали дефекта"></textarea>
        </label>
        <div class="form-actions">
          <button class="primary-btn" type="submit">Создать</button>
        </div>
      </form>

      <form class="filter-bar" @submit.prevent="applyFilters">
        <label>
          Проект
          <select v-model="filter.projectId">
            <option value="">Все</option>
            <option v-for="project in projectsStore.items" :key="project.id" :value="project.id">
              {{ project.name }}
            </option>
          </select>
        </label>
        <label>
          Статус
          <select v-model="filter.status">
            <option value="">Все</option>
            <option v-for="status in statuses" :key="status" :value="status">
              {{ status }}
            </option>
          </select>
        </label>

        <label>
          Приоритет
          <select v-model="filter.priority">
            <option value="">Все</option>
            <option v-for="priority in priorities" :key="priority" :value="priority">
              {{ priority }}
            </option>
          </select>
        </label>

        <button class="secondary-btn" type="submit">Применить</button>
      </form>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Проект</th>
              <th>Заголовок</th>
              <th>Приоритет</th>
              <th>Статус</th>
              <th>Срок</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="defectsStore.loading">
              <td colspan="6" class="table-state">Загрузка...</td>
            </tr>
            <tr v-else-if="defectsStore.error">
              <td colspan="6" class="table-state table-state--error">
                {{ defectsStore.error }}
              </td>
            </tr>
            <tr v-else v-for="item in defectsStore.items" :key="item.id">
              <td>{{ item.project || '—' }}</td>
              <td>{{ item.title }}</td>
              <td>
                <span class="pill" :data-priority="item.priority">
                  {{ item.priority }}
                </span>
              </td>
              <td>{{ item.status }}</td>
              <td>{{ item.dueDate ?? '—' }}</td>
              <td>
                <button class="link-btn" type="button" @click="openDetail(item)">
                  Открыть
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <aside class="defect-detail" v-if="selectedId">
      <div v-if="detailLoading" class="panel-state">Загрузка карточки...</div>
      <div v-else-if="selectedDefect">
        <header>
          <p class="section-label">Карточка дефекта</p>
          <h3>{{ selectedDefect.title }}</h3>
          <p class="detail-meta">
            {{ selectedDefect.project }} · {{ selectedDefect.status }}
          </p>
        </header>

        <p class="detail-description">{{ selectedDefect.description || 'Нет описания' }}</p>

        <div class="detail-block">
          <h4>Комментарии</h4>
          <div class="comments">
            <p v-if="!defectsStore.comments.length" class="muted">Комментариев пока нет</p>
            <article v-for="comment in defectsStore.comments" :key="comment.id" class="comment">
              <p class="comment__author">{{ comment.author || 'Без автора' }}</p>
              <p class="comment__body">{{ comment.body }}</p>
            </article>
          </div>
          <form class="comment-form" @submit.prevent="submitComment">
            <textarea v-model="commentBody" rows="2" placeholder="Добавить комментарий"></textarea>
            <button class="secondary-btn" type="submit">Отправить</button>
          </form>
        </div>

        <div class="detail-block">
          <h4>Вложения</h4>
          <p v-if="!selectedDefect.attachments?.length" class="muted">Файлы не прикреплены</p>
          <ul class="attachments">
            <li v-for="att in selectedDefect.attachments" :key="att.id">
              <a :href="att.downloadUrl" target="_blank">
                {{ att.filename }}
              </a>
              ({{ (att.sizeBytes / 1024).toFixed(1) }} КБ)
            </li>
          </ul>
          <form class="attachment-form" @submit.prevent="uploadAttachment">
            <input ref="attachmentInput" type="file" />
            <button class="secondary-btn" type="submit">Загрузить</button>
          </form>
        </div>
      </div>
      <div v-else class="panel-state">Выберите дефект в таблице</div>
    </aside>
  <p v-else class="muted">Авторизуйтесь, чтобы работать с дефектами.</p>
</section>
</template>

<style scoped>
.defects-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

@media (max-width: 1200px) {
  .defects-grid {
    grid-template-columns: 1fr;
  }
}

.defects__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-label {
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.85rem;
  color: #94a3b8;
}

.create-card,
.filter-bar,
.defect-detail,
.detail-block {
  background: rgba(15, 23, 42, 0.7);
  border-radius: 16px;
  border: 1px solid rgba(148, 163, 184, 0.12);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
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
select,
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

.secondary-btn {
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 0.6rem 1.2rem;
  background: transparent;
  color: #e2e8f0;
  cursor: pointer;
}

.filter-bar {
  flex-direction: row;
  flex-wrap: wrap;
}

.filter-bar label {
  flex: 1;
  min-width: 180px;
}

.table-wrapper {
  overflow-x: auto;
  background: rgba(15, 23, 42, 0.9);
  border-radius: 16px;
  border: 1px solid rgba(148, 163, 184, 0.1);
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 0.85rem 1rem;
  text-align: left;
}

thead {
  background: rgba(148, 163, 184, 0.1);
}

.table-state {
  text-align: center;
  padding: 2rem 0;
}

.table-state--error {
  color: #f87171;
}

.pill {
  padding: 0.3rem 0.7rem;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
}

.pill[data-priority='HIGH'],
.pill[data-priority='CRITICAL'] {
  background: rgba(248, 113, 113, 0.15);
  color: #fecaca;
}

.pill[data-priority='MEDIUM'] {
  background: rgba(251, 191, 36, 0.2);
  color: #fcd34d;
}

.pill[data-priority='LOW'] {
  background: rgba(34, 197, 94, 0.2);
  color: #86efac;
}

.link-btn {
  background: transparent;
  border: none;
  color: #60a5fa;
  cursor: pointer;
}

.defect-detail .detail-meta {
  color: #94a3b8;
  margin-top: 0.25rem;
}

.detail-description {
  margin: 0;
  color: #e2e8f0;
}

.comments {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.comment {
  background: rgba(15, 23, 42, 0.6);
  border-radius: 10px;
  padding: 0.75rem;
  border: 1px solid rgba(148, 163, 184, 0.1);
}

.comment__author {
  font-size: 0.85rem;
  color: #94a3b8;
  margin: 0 0 0.3rem;
}

.comment__body {
  margin: 0;
}

.comment-form textarea {
  width: 100%;
}

.attachments {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.attachment-form {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.muted {
  color: #94a3b8;
}

.panel-state {
  text-align: center;
  color: #94a3b8;
}
</style>
