<script setup>
import { RouterLink, RouterView } from 'vue-router'
import AuthPanel from './components/AuthPanel.vue'
</script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div>
        <p class="app-header__label">Система управления дефектами</p>
        <h1 class="app-header__title">Контроль строительных объектов</h1>
      </div>
      <div class="app-header__meta">
        <span class="badge">v0.2.0</span>
        <AuthPanel />
      </div>
    </header>

    <nav class="app-nav">
      <RouterLink to="/" exact-active-class="app-nav__link--active" class="app-nav__link">
        Дашборд
      </RouterLink>
      <RouterLink to="/defects" exact-active-class="app-nav__link--active" class="app-nav__link">
        Дефекты
      </RouterLink>
      <RouterLink to="/projects" exact-active-class="app-nav__link--active" class="app-nav__link">
        Объекты
      </RouterLink>
    </nav>

    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
  padding: 2rem;
  background: linear-gradient(180deg, #0f172a 0%, #111827 100%);
  color: #e2e8f0;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.app-header__label {
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #94a3b8;
}

.app-header__title {
  font-size: 2rem;
  margin: 0.2rem 0 0;
}

.app-header__meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  align-items: flex-end;
}

.badge {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border-radius: 999px;
  padding: 0.35rem 0.9rem;
  font-size: 0.85rem;
  font-weight: 600;
}

.app-nav {
  margin-top: 1.5rem;
  display: flex;
  gap: 1rem;
}

.app-nav__link {
  padding: 0.5rem 1.25rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.05);
  text-decoration: none;
  color: inherit;
  transition: background 0.3s ease;
}

.app-nav__link--active,
.app-nav__link:hover {
  background: rgba(59, 130, 246, 0.25);
  color: #fff;
}

.app-main {
  margin-top: 2rem;
}
</style>
