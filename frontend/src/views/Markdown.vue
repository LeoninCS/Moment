<!-- src/views/Markdown.vue -->
<template>
  <div class="md-root">
    <div class="md-hero">
      <div class="md-text">
        <h1 class="md-title">在线 Markdown 协同编辑</h1>
        <p class="md-subtitle">
          创建一个带哈希链接的文档，支持多端实时同步编辑，随时导出 HTML / PDF / MD。
        </p>

        <ul class="md-bullets">
          <li>🧑‍💻 左侧编辑 Markdown，右侧实时预览效果</li>
          <li>🤝 把链接发给同事，即可一起在线协作编辑</li>
          <li>📄 一键导出为 HTML / PDF / Markdown 文件</li>
        </ul>

        <div class="md-actions">
          <button
            class="md-start-btn"
            :disabled="loading"
            @click="handleStart"
          >
            <span v-if="!loading">✏️ 创建新文档</span>
            <span v-else>生成中...</span>
          </button>

          <p v-if="error" class="md-error">
            {{ error }}
          </p>

          <p v-if="hash" class="md-hint">
            已生成文档：<code>{{ hash }}</code>，正在跳转到编辑页面…
          </p>
        </div>
      </div>

      <div class="md-hero-preview">
        <div class="md-preview-window">
          <div class="md-preview-header">
            <span></span><span></span><span></span>
          </div>
          <div class="md-preview-body">
            <div class="md-preview-left">
              <div class="md-line" v-for="i in 8" :key="i"></div>
            </div>
            <div class="md-preview-right">
              <div class="md-block-title"></div>
              <div class="md-block-line" v-for="i in 5" :key="i"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { createMarkdownDoc } from "../api/markdown.ts";

const router = useRouter();

const loading = ref(false);
const error = ref("");
const hash = ref("");

const handleStart = async () => {
  if (loading.value) return;
  loading.value = true;
  error.value = "";

  try {
    const res = await createMarkdownDoc();
    hash.value = res.data.hash;

    router.push({
      name: "MarkdownEditorView",
      params: { hash: hash.value },
    });
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || "生成失败，请稍后重试";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* 外层只是做一个宽度和间距控制，不再铺满全屏背景 */
.md-root {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 12px 8px 24px;
  box-sizing: border-box;
}

/* 真正的小卡片：白色背景 + 圆角 + 阴影 */
.md-hero {
  display: flex;
  gap: 40px;
  align-items: center;
  background: var(--moment-card-bg, #ffffff);
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
  padding: 24px 24px 26px;
  box-sizing: border-box;
}

.md-text {
  flex: 1.2;
}

.md-title {
  margin: 0 0 8px;
  font-size: 28px;
  font-weight: 800;
  color: #0f172a;
}

.md-subtitle {
  margin: 0 0 18px;
  font-size: 15px;
  color: #4b5563;
}

.md-bullets {
  margin: 0 0 24px;
  padding-left: 18px;
  font-size: 14px;
  color: #374151;
}

.md-bullets li + li {
  margin-top: 4px;
}

.md-actions {
  margin-top: 8px;
}

.md-start-btn {
  min-width: 220px;
  padding: 11px 22px;
  border-radius: 999px;
  border: none;
  outline: none;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  background: linear-gradient(120deg, #2563eb, #4f46e5);
  color: #ffffff;
  box-shadow: 0 14px 40px rgba(37, 99, 235, 0.5);
  transition: transform 0.08s ease, box-shadow 0.08s ease, opacity 0.2s;
}

.md-start-btn:disabled {
  opacity: 0.7;
  cursor: default;
  box-shadow: none;
}

.md-start-btn:not(:disabled):hover {
  transform: translateY(-1px);
}

.md-start-btn:not(:disabled):active {
  transform: translateY(1px);
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.45);
}

.md-error {
  margin-top: 10px;
  color: #dc2626;
  font-size: 13px;
}

.md-hint {
  margin-top: 10px;
  font-size: 13px;
  color: #4b5563;
}

.md-hint code {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

/* 右侧预览插画 */
.md-hero-preview {
  flex: 1;
  display: flex;
  justify-content: center;
}

.md-preview-window {
  width: 100%;
  max-width: 420px;
  border-radius: 18px;
  background: #0f172a;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.6);
  overflow: hidden;
}

.md-preview-header {
  height: 28px;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.md-preview-header span {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #475569;
}

.md-preview-body {
  display: flex;
  padding: 12px;
  gap: 10px;
  background: radial-gradient(circle at top left, #1f2937, #020617);
}

.md-preview-left {
  flex: 1;
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.9);
  padding: 10px 8px;
}

.md-line {
  height: 6px;
  margin-bottom: 6px;
  border-radius: 999px;
  background: linear-gradient(90deg, #4b5563, #9ca3af);
  opacity: 0.85;
}

.md-preview-right {
  flex: 1;
  border-radius: 12px;
  background: #f9fafb;
  padding: 10px 10px 16px;
}

.md-block-title {
  width: 70%;
  height: 14px;
  border-radius: 999px;
  background: #e5e7eb;
  margin-bottom: 10px;
}

.md-block-line {
  width: 100%;
  height: 8px;
  border-radius: 999px;
  background: #e5e7eb;
  margin-bottom: 6px;
}

@media (max-width: 900px) {
  .md-root {
    padding: 12px 4px 20px;
  }
  .md-hero {
    flex-direction: column;
    align-items: stretch;
  }
  .md-hero-preview {
    order: -1;
    margin-bottom: 12px;
  }
}
</style>

