<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ApiError, getCase, listCases, runAutopsy } from './api'
import type { CaseDetail, CaseItem, Report } from './api'
import CaseFile from './components/CaseFile.vue'
import CaseList from './components/CaseList.vue'
import LocaleToggle from './components/LocaleToggle.vue'
import AutopsyReport from './views/AutopsyReport.vue'

const { t, locale } = useI18n()
const items = ref<CaseItem[]>([])
const selectedId = ref<string | null>(null)
const detail = ref<CaseDetail | null>(null)
const report = ref<Report | null>(null)
const running = ref(false)
const errorKey = ref<string | null>(null)

watch(locale, (l) => { document.documentElement.lang = l }, { immediate: true })

async function guarded(fn: () => Promise<void>) {
  errorKey.value = null
  try {
    await fn()
  } catch (e) {
    errorKey.value = e instanceof ApiError ? e.message : 'error.network'
  }
}

async function select(id: string) {
  selectedId.value = id
  report.value = null
  detail.value = null
  await guarded(async () => {
    const d = await getCase(id)
    if (selectedId.value === id) detail.value = d
  })
}

async function run() {
  if (!selectedId.value || running.value) return
  running.value = true
  const id = selectedId.value
  await guarded(async () => {
    const r = await runAutopsy(id)
    if (selectedId.value === id) report.value = r // drop a result for a case the user already left
  })
  running.value = false
}

onMounted(() => guarded(async () => {
  items.value = await listCases()
  if (items.value.length) await select(items.value[0].id)
}))
</script>

<template>
  <div class="app">
    <header class="top no-print">
      <div>
        <h1>{{ t('app.title') }}</h1>
        <p class="tagline">{{ t('app.tagline') }}</p>
      </div>
      <LocaleToggle />
    </header>
    <p class="synthetic no-print">{{ t('app.synthetic') }}</p>
    <p v-if="errorKey" class="error" role="alert">{{ t(errorKey) }}</p>

    <main class="layout">
      <aside class="sidebar no-print">
        <CaseList :items="items" :selected-id="selectedId" @select="select" />
      </aside>
      <div class="work">
        <CaseFile v-if="detail" :detail="detail" />
        <button v-if="detail" class="primary no-print" :disabled="running" @click="run">
          {{ running ? t('autopsy.running') : t('autopsy.run') }}
        </button>
        <AutopsyReport v-if="report" :report="report" />
        <p v-else-if="detail" class="muted no-print">{{ t('report.empty') }}</p>
      </div>
    </main>
  </div>
</template>
