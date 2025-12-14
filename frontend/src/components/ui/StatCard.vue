<template>
  <div class="stat-card">
    <p class="stat-card__label">{{ label }}</p>
    <p class="stat-card__value">{{ value }}</p>
    <p v-if="hint" class="stat-card__hint">{{ hint }}</p>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  value: {
    type: [String, Number],
    required: true,
  },
  hint: {
    type: String,
    default: '',
  },
})
</script>

<style scoped>
.stat-card {
  border-radius: 12px;
  padding: 1.25rem;
  background: linear-gradient(135deg, #1f2937, #111827);
  color: #fff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.35);
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.stat-card__label {
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.7);
}

.stat-card__value {
  font-size: 2rem;
  font-weight: 600;
}

.stat-card__hint {
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.6);
}
</style>
