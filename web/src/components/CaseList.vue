<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { CaseItem } from '../api'

defineProps<{ items: CaseItem[]; selectedId: string | null }>()
const emit = defineEmits<{ select: [id: string] }>()
const { t } = useI18n()
</script>

<template>
  <nav class="case-list" :aria-label="t('cases.heading')">
    <h2>{{ t('cases.heading') }}</h2>
    <p v-if="!items.length" class="muted">{{ t('cases.loading') }}</p>
    <ul>
      <li v-for="c in items" :key="c.id">
        <button :class="{ on: c.id === selectedId }" :aria-current="c.id === selectedId" @click="emit('select', c.id)">
          <span class="mono id">{{ c.id }}</span>
          <span>{{ c.title }}</span>
        </button>
      </li>
    </ul>
  </nav>
</template>
