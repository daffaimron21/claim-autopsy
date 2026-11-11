<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Locale } from '../i18n'

const { locale, t } = useI18n()

function set(l: Locale) {
  locale.value = l
  localStorage.setItem('locale', l)
}
</script>

<template>
  <div class="locale" role="group" :aria-label="t('locale.label')">
    <button v-for="l in (['en', 'id'] as Locale[])" :key="l" :class="{ on: locale === l }"
      :aria-pressed="locale === l" @click="set(l)">{{ t(`locale.${l}`) }}</button>
  </div>
</template>
