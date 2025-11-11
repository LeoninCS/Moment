<template>
  <section>
    <h1 style="margin:0 0 12px 0;">代码查看</h1>
    <p class="helper">此页面为只读视图，来自唯一哈希：<code>{{ hash }}</code></p>
    <div class="hr"></div>

    <div v-if="loading">正在加载...</div>
    <div v-else-if="error" class="helper">获取失败：{{ error }}</div>
    <div v-else-if="!data">未找到该代码片段，可能已过期或不存在。</div>
    <div v-else>
      <div class="meta">
        <div class="item">作者：{{ data.author }}</div>
        <div class="item">语言：{{ data.language }}</div>
        <div class="item" v-if="data.destroy_time">销毁时间：{{ formatTime(data.destroy_time) }}</div>
      </div>

      <div style="height:14px"></div>

      <div class="codewrap" :data-theme="theme">
        <div class="toolbar">
          <button class="theme-btn" @click="toggleTheme">
            {{ theme === 'dark' ? '☀️ 浅色' : '🌙 深色' }}
          </button>
          <button class="copy-btn" @click="copyCode">{{ copied ? '已复制!' : '复制' }}</button>
        </div>

        <pre class="line-numbers" :class="`language-${prismLang}`"><code :class="`language-${prismLang}`" v-html="highlightedHtml"></code></pre>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref, watch, nextTick, computed } from 'vue'
import { getCode } from '../utils/api'
import Prism from 'prismjs'

// 📏 行号插件
import 'prismjs/plugins/line-numbers/prism-line-numbers.css'
import 'prismjs/plugins/line-numbers/prism-line-numbers.js'

// ⚠️ 很关键：很多语言依赖 clike
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-java'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'
import 'prismjs/components/prism-rust'
import 'prismjs/components/prism-php'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-swift'
import 'prismjs/components/prism-kotlin'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-json'

// 🌈 使用更护眼的浅色主题 + 深色主题
import prismLightUrl from 'prismjs/themes/prism-solarizedlight.css?url'
import prismDarkUrl from 'prismjs/themes/prism-tomorrow.css?url'

const props = defineProps({ hash: String })
const hash = props.hash

const data = ref(null)
const loading = ref(true)
const error = ref('')
const copied = ref(false)
const prismLang = ref('none')
const theme = ref('dark')
let prismThemeLinkEl = null

// 语言映射表
const langMap = {
  plaintext: 'none',
  javascript: 'javascript',
  typescript: 'typescript',
  python: 'python',
  java: 'java',
  go: 'go',
  c: 'c',
  cpp: 'cpp',
  rust: 'rust',
  php: 'php',
  ruby: 'ruby',
  swift: 'swift',
  kotlin: 'kotlin',
  shell: 'bash',
  json: 'json',
  html: 'markup',
  xml: 'markup',
  markdown: 'markup'
}

// ======== 高亮 HTML（最稳定方案）========
const highlightedHtml = computed(() => {
  const src = data.value?.content ?? ''
  const lang = prismLang.value
  const grammar = Prism.languages[lang]
  try {
    if (!grammar) {
      return Prism.util.encode(src)
    }
    return Prism.highlight(src, grammar, lang)
  } catch (e) {
    return Prism.util.encode(src)
  }
})

// ======== 主题切换 ========
function ensurePrismThemeLink() {
  if (!prismThemeLinkEl) {
    prismThemeLinkEl = document.getElementById('prism-theme')
  }
  if (!prismThemeLinkEl) {
    prismThemeLinkEl = document.createElement('link')
    prismThemeLinkEl.rel = 'stylesheet'
    prismThemeLinkEl.id = 'prism-theme'
    document.head.appendChild(prismThemeLinkEl)
  }
}

function applyThemeCSS() {
  ensurePrismThemeLink()
  prismThemeLinkEl.href = theme.value === 'dark' ? prismDarkUrl : prismLightUrl
}

// ======== 通用方法 ========
function formatTime(ts) {
  try {
    return new Date(ts * 1000).toLocaleString()
  } catch {
    return ''
  }
}

async function load() {
  try {
    loading.value = true
    const res = await getCode(hash)
    data.value = res
    const langKey = (res?.language || '').toString().toLowerCase()
    prismLang.value = langMap[langKey] || 'none'
    await nextTick()
  } catch (e) {
    error.value = e?.response?.status === 404 ? '代码片段不存在或已过期' : (e?.message || '请求失败')
    data.value = null
  } finally {
    loading.value = false
  }
}

async function copyCode() {
  try {
    await navigator.clipboard.writeText(data.value?.content || '')
    copied.value = true
    setTimeout(() => (copied.value = false), 1500)
  } catch {
    alert('复制失败，请手动复制')
  }
}

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
}

// ======== 生命周期 ========
onMounted(async () => {
  applyThemeCSS()
  await load()
})

watch(() => props.hash, load)
watch(theme, async () => {
  applyThemeCSS()
  await nextTick()
})
watch(() => data.value?.content, async () => {
  await nextTick()
})
</script>

<style scoped>
.codewrap {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: flex-end;
  background: var(--toolbar-bg);
  padding: 6px 10px;
  gap: 6px;
}

.copy-btn,
.theme-btn {
  background: var(--btn-bg);
  color: var(--btn-color);
  border: none;
  border-radius: 4px;
  padding: 4px 10px;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.2s, transform 0.1s;
}
.copy-btn:hover,
.theme-btn:hover {
  background: var(--btn-bg-hover);
  transform: scale(1.05);
}
.copy-btn:active,
.theme-btn:active {
  transform: scale(0.95);
}

/* 明暗模式变量（容器 UI） */
[data-theme='dark'] {
  --toolbar-bg: #2d2d2d;
  --btn-bg: #3c3c3c;
  --btn-bg-hover: #555;
  --btn-color: #fff;
  background: #1e1e1e;
  color: #eee;
}

[data-theme='light'] {
  --toolbar-bg: #f5f5dc;
  --btn-bg: #ddd;
  --btn-bg-hover: #ccc;
  --btn-color: #000;
  background: #faf9f2;
  color: #000;
}

/* 修复行号左侧空白过多 */
pre.line-numbers {
  padding-left: 2.8em !important;
}

/* 行号与文本对齐 */
pre.line-numbers .line-numbers-rows {
  top: 0;
  left: 0;
  padding-top: 12px;
}

/* 兜底背景（主题 CSS 加载前） */
[data-theme='dark'] pre { background: #2d2d2d; color: #eee; }
[data-theme='light'] pre { background: #faf9f2; color: #222; }

pre {
  margin: 0;
  padding: 12px 16px;
  overflow-x: auto;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 14px;
  line-height: 1.5;
  border-radius: 0 0 8px 8px;
}
</style>
