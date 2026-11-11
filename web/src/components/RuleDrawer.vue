<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Source } from '../api'

const props = defineProps<{ ruleId: string; source: Source }>()
const emit = defineEmits<{ close: [] }>()
const { t } = useI18n()
const key = computed(() => props.ruleId.replace(/-/g, '_'))

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}
onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <aside class="drawer" role="dialog" aria-modal="true" data-test="drawer">
    <button class="close no-print" @click="emit('close')">{{ t('drawer.close') }}</button>
    <p class="mono rule-id">{{ ruleId }}</p>
    <h3>{{ t(`rule.${key}.title`) }}</h3>
    <p>{{ t(`rule.${key}.desc`) }}</p>
    <template v-if="source === 'rule'">
      <p class="note">{{ t('drawer.deterministic') }}</p>
      <p class="note simplified">{{ t('drawer.simplified') }}</p>
    </template>
    <p v-else class="note">{{ t('drawer.modelNote') }}</p>
  </aside>
</template>
