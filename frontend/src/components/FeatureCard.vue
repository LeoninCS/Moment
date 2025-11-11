<template>
  <div class="card">
    <div class="badge" v-if="badge">{{ badge }}</div>
    <h3 style="margin-top:8px">{{ title }}</h3>
    <p>{{ description }}</p>
    <div style="height:12px"></div>
    <slot name="action">
      <router-link :to="to" class="btn primary" v-if="to">进入</router-link>
    </slot>
  </div>
</template>

<script setup>
defineProps({
  title: String,
  description: String,
  to: String,
  badge: String
})
</script>
