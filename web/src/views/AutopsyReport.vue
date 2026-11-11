<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Finding, Report } from '../api'
import FindingCard from '../components/FindingCard.vue'
import RuleDrawer from '../components/RuleDrawer.vue'

defineProps<{ report: Report }>()
const { t } = useI18n()
const open = ref<Finding | null>(null)
const print = () => window.print()
</script>

<template>
  <section class="report">
    <header class="report-head">
      <h2>{{ t('report.heading') }}</h2>
      <button class="no-print" @click="print">{{ t('report.print') }}</button>
    </header>

    <p v-if="report.degraded" class="notice" data-test="degraded" role="status">{{ t('report.degraded') }}</p>

    <div class="verdict-row">
      <span class="label">{{ t('report.verdict') }}</span>
      <span class="stamp" :class="report.verdict" data-test="verdict">{{ t(`verdict.${report.verdict}`) }}</span>
    </div>

    <div v-if="report.primary" data-test="primary">
      <h3>{{ t('report.primary') }}</h3>
      <FindingCard :finding="report.primary" @why="open = report.primary" />
    </div>
    <p v-else class="no-cause" data-test="no-cause">{{ t('report.noCause') }}</p>

    <div v-if="report.secondary.length" data-test="secondary">
      <h3>{{ t('report.secondary') }}</h3>
      <FindingCard v-for="f in report.secondary" :key="f.rule_id + f.message_key" :finding="f" @why="open = f" />
    </div>

    <RuleDrawer v-if="open" :rule-id="open.rule_id" :source="open.source" @close="open = null" />
  </section>
</template>
