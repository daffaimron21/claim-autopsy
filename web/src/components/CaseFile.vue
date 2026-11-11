<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { CaseDetail } from '../api'

defineProps<{ detail: CaseDetail }>()
const { t } = useI18n()
</script>

<template>
  <section class="case-file">
    <h2>{{ t('file.heading') }}: <span class="mono">{{ detail.id }}</span></h2>
    <p class="title">{{ detail.title }}</p>
    <dl>
      <div><dt>{{ t('file.cbg') }}</dt><dd class="mono">{{ detail.cbg_code }}</dd></div>
      <div>
        <dt>{{ t('file.diagnoses') }}</dt>
        <dd class="mono">
          <span v-for="d in detail.diagnoses" :key="d.code" class="code">{{ d.code }}<em v-if="d.primary"> ({{ t('file.primary') }})</em></span>
        </dd>
      </div>
      <div>
        <dt>{{ t('file.procedures') }}</dt>
        <dd class="mono">
          <span v-for="p in detail.procedures" :key="p.code" class="code">{{ p.code }}</span>
          <span v-if="!detail.procedures.length" class="muted">{{ t('file.none') }}</span>
        </dd>
      </div>
      <div><dt>{{ t('file.admit') }}</dt><dd class="mono">{{ detail.admit_date }}</dd></div>
      <div><dt>{{ t('file.discharge') }}</dt><dd class="mono">{{ detail.discharge_date }}</dd></div>
      <div><dt>{{ t('file.careClass') }}</dt><dd class="mono">{{ detail.care_class }}</dd></div>
      <div class="wide"><dt>{{ t('file.returnReason') }}</dt><dd>{{ detail.return_reason }}</dd></div>
      <div class="wide">
        <dt>{{ t('file.summary') }} <span class="muted">({{ t('file.summaryNote') }})</span></dt>
        <dd class="summary">{{ detail.resume_text }}</dd>
      </div>
    </dl>
  </section>
</template>
