<template>
  <section>
    <h1 style="margin:0 0 12px 0;">文本加解密</h1>
    <p class="helper">
      输入密钥和文本，选择加密或解密，点击提交后，下方将展示只读结果。
    </p>
    <div class="hr"></div>

    <!-- 表单区域 -->
    <form class="form" @submit.prevent="onSubmit">
      <!-- 密钥输入 + 随机密钥按钮 -->
      <div>
        <label>密钥</label>
        <div class="key-row">
          <input
            class="input"
            v-model.trim="secretKey"
            placeholder="请输入密钥，或点击右侧钥匙生成"
            required
          />
          <button
            type="button"
            class="btn icon-btn"
            @click="onGenerateKey"
            :disabled="loading"
            title="生成随机密钥"
          >
            🔑
          </button>
        </div>
        <p class="helper">后端使用同一个密钥进行加密和解密，请妥善保存。</p>
      </div>

      <!-- 文本输入 -->
      <div>
        <label>文本内容</label>
        <textarea
          class="textarea"
          v-model="inputText"
          placeholder="在此输入要加密的原文或要解密的密文..."
          :maxlength="maxLen"
          required
        />
        <p class="helper">{{ inputText.length }} / {{ maxLen }} 字符</p>
      </div>

      <!-- 模式选择：加密 / 解密 -->
      <div>
        <label>操作类型</label>
        <div class="mode-toggle">
          <label :class="{ active: mode === 'encrypt' }">
            <input
              type="radio"
              name="crypto-mode"
              value="encrypt"
              v-model="mode"
            />
            加密
          </label>
          <label :class="{ active: mode === 'decrypt' }">
            <input
              type="radio"
              name="crypto-mode"
              value="decrypt"
              v-model="mode"
            />
            解密
          </label>
        </div>
        <p class="helper">
          · 选择 <strong>加密</strong>：文本视为明文，调用后端 <code>Encrypt</code>（<code>/textcrypto/encrypt-text</code>）。<br />
          · 选择 <strong>解密</strong>：文本视为密文，调用后端 <code>Decrypt</code>（<code>/textcrypto/decrypt-text</code>）。
        </p>
      </div>

      <!-- 操作按钮 -->
      <div class="actions">
        <button class="btn primary" :disabled="loading">
          {{ loading ? '处理中...' : (mode === 'encrypt' ? '加密并显示结果' : '解密并显示结果') }}
        </button>
        <button type="button" class="linklike" @click="reset" :disabled="loading">
          重置
        </button>
        <span class="helper" v-if="error" style="color:#d33; margin-left:8px;">
          {{ error }}
        </span>
      </div>
    </form>

    <div style="height:20px"></div>
    <div class="hr"></div>

    <!-- 结果展示区域（原 TextView 内容整合进来） -->
    <section style="margin-top:16px;">
      <h2 style="margin:0 0 12px 0;">文本查看</h2>
      <p class="helper">
        此区域为只读视图，展示{{ modeLabel }}后的文本内容。
      </p>

      <div v-if="!resultText">
        暂未产生结果，请先在上方完成加密或解密操作。
      </div>
      <div v-else>
        <div class="meta">
          <div class="item">操作类型：{{ modeLabel }}</div>
        </div>

        <div style="height:14px"></div>

        <div class="textwrap" :data-theme="theme">
          <div class="toolbar">
            <button class="theme-btn" @click="toggleTheme">
              {{ theme === 'dark' ? '☀️ 浅色' : '🌙 深色' }}
            </button>
            <button class="copy-btn" @click="copyText">
              {{ copied ? '已复制!' : '复制' }}
            </button>
          </div>

          <pre class="content">{{ resultText }}</pre>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup>
import { ref, computed } from 'vue'
import { getRandomSecretKey, encryptText, decryptText } from '../utils/api'

const maxLen = 100000

const secretKey = ref('')
const inputText = ref('')
const mode = ref('encrypt') // 'encrypt' | 'decrypt'
const error = ref('')
const loading = ref(false)

// 展示用的结果文本（原 TextView.rawText）
const resultText = ref('')
const theme = ref('dark')
const copied = ref(false)

const modeLabel = computed(() =>
  mode.value === 'encrypt' ? '加密' : '解密'
)

function reset () {
  secretKey.value = ''
  inputText.value = ''
  error.value = ''
  mode.value = 'encrypt'
  resultText.value = ''
}

// 生成随机密钥
async function onGenerateKey () {
  error.value = ''
  try {
    loading.value = true
    const key = await getRandomSecretKey()
    if (!key) {
      throw new Error('未从服务端获取到密钥')
    }
    secretKey.value = key
  } catch (e) {
    error.value = e.message || '生成密钥失败'
  } finally {
    loading.value = false
  }
}

function toggleTheme () {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
}

async function copyText () {
  if (!resultText.value) return
  try {
    await navigator.clipboard.writeText(resultText.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 1500)
  } catch {
    alert('复制失败，请手动复制')
  }
}

// 提交加密/解密，然后在本页下方展示结果
async function onSubmit () {
  error.value = ''
  resultText.value = ''

  if (!secretKey.value || !inputText.value) {
    error.value = '请填写密钥和文本内容。'
    return
  }

  try {
    loading.value = true

    let text = ''

    if (mode.value === 'encrypt') {
      const res = await encryptText({
        text: inputText.value,
        screat_key: secretKey.value
      })
      console.log('encrypt res:', res)
      text =
        res.encrypted_text ||
        res.encryptedText ||
        res.cipher_text ||
        res.cipherText ||
        res.text ||
        ''
    } else {
      const res = await decryptText({
        encrypted_text: inputText.value,
        screat_key: secretKey.value
      })
      console.log('decrypt res:', res)
      text =
        res.decrypted_text ||
        res.decryptedText ||
        res.plain_text ||
        res.plainText ||
        res.text ||
        ''
    }

    if (!text) {
      error.value = '请求成功，但未获得结果字段。请看控制台返回结构。'
      return
    }

    resultText.value = text
  } catch (e) {
    error.value = e.message || '请求失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.key-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.icon-btn {
  min-width: 42px;
  padding: 0 10px;
}

/* 模式切换样式 */
.mode-toggle {
  display: inline-flex;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #ddd;
  background: #f7f7f7;
}

.mode-toggle label {
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
  user-select: none;
  display: flex;
  align-items: center;
  gap: 4px;
}

.mode-toggle input {
  margin: 0;
}

.mode-toggle label.active {
  background: #1d7dfa;
  color: #fff;
}

/* ====== 以下是原 TextView 的样式，直接搬过来 ====== */
.meta {
  display: flex;
  gap: 16px;
  font-size: 14px;
}
.meta .item {
  color: #666;
}

/* 外层容器，带主题变量 */
.textwrap {
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

/* 明暗模式变量 */
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

/* 文本区域 */
.content {
  margin: 0;
  padding: 12px 16px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 14px;
  line-height: 1.5;
  border-radius: 0 0 8px 8px;
}
</style>
