<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Finding } from '../api'

defineProps<{ finding: Finding }>()
const emit = defineEmits<{ why: [] }>()
const { t } = useI18n()
const SEVERITY = ['minor', 'major', 'blocking'] as const
</script>

<template>
  <article class="finding" :class="`sev-${SEVERITY[finding.severity]}`" data-test="finding">
    <header>
      <span class="chip severity">{{ t(`severity.${SEVERITY[finding.severity]}`) }}</span>
      <span class="chip" :class="`src-${finding.source}`" :data-test="`source-${finding.source}`">{{ t(`source.${finding.source}`) }}</span>
      <span class="mono rule-id">{{ finding.rule_id }}</span>
    </header>
    <p class="message">{{ t(finding.message_key, finding.params ?? {}) }}</p>
    <ul v-if="finding.evidence?.length" class="evidence" :aria-label="t('report.evidence')">
      <li v-for="(e, i) in finding.evidence" :key="i">
        <span class="label">{{ t(`field.${e.field}`) }}</span>
        <q class="mono">{{ e.quote }}</q>
      </li>
    </ul>
    <p class="fix"><strong>{{ t('report.fix') }}.</strong> {{ t(finding.fix_key) }}</p>
    <button class="link no-print" data-test="why" @click="emit('why')">{{ t('drawer.show') }}</button>
  </article>
</template>
