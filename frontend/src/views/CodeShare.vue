<template>
  <section>
    <h1 style="margin:0 0 12px 0;">代码分享</h1>
    <p class="helper">填写信息并提交，系统将返回一个哈希链接。任何人访问该链接即可<strong>只读</strong>查看代码。</p>
    <div class="hr"></div>

    <form class="form" @submit.prevent="onSubmit">
      <div>
        <label>作者</label>
        <input class="input" v-model.trim="author" placeholder="请输入作者名称（必填）" maxlength="64" required />
      </div>

      <div>
        <label>语言</label>
        <select class="select" v-model="language" required>
          <option disabled value="">选择语言</option>
          <option v-for="opt in languages" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
      </div>

      <div>
        <label>保存时间</label>
        <select class="select" v-model.number="ttl" required>
          <option v-for="opt in ttls" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
        <p class="helper">后端默认 1 小时（3600 秒）。</p>
      </div>

      <div>
        <label>代码内容</label>
        <textarea class="textarea" v-model="content" placeholder="在此粘贴或输入代码..." :maxlength="maxLen" />
        <p class="helper">{{ content.length }} / {{ maxLen }} 字符</p>
      </div>

      <div class="actions">
        <button class="btn primary" :disabled="submitting">提交并生成链接</button>
        <button type="button" class="linklike" @click="reset" :disabled="submitting">重置</button>
        <span class="helper" v-if="error">{{ error }}</span>
      </div>
    </form>
  </section>
</template>

<script setup>
import { ref } from 'vue'
import { uploadCode } from '../utils/api'
import { useRouter } from 'vue-router'

const router = useRouter()
const maxLen = 100000

const author = ref('')
const language = ref('')
const content = ref('')
const ttl = ref(3600)
const error = ref('')
const submitting = ref(false)

const languages = [
  { value: 'plaintext', label: 'Plain Text' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'go', label: 'Go' },
  { value: 'c', label: 'C' },
  { value: 'cpp', label: 'C++' },
  { value: 'rust', label: 'Rust' },
  { value: 'php', label: 'PHP' },
  { value: 'ruby', label: 'Ruby' },
  { value: 'swift', label: 'Swift' },
  { value: 'kotlin', label: 'Kotlin' },
  { value: 'shell', label: 'Shell' }
]

const ttls = [
  { value: 600, label: '10 分钟' },
  { value: 3600, label: '1 小时（默认）' },
  { value: 86400, label: '24 小时' },
  { value: 604800, label: '7 天' }
]

function reset(){
  author.value=''; language.value=''; content.value=''; ttl.value=3600; error.value=''
}

async function onSubmit(){
  error.value=''
  if(!author.value || !language.value || !content.value){
    error.value = '请完整填写作者、语言和代码内容。'
    return
  }
  if(content.value.length > maxLen){
    error.value = '代码内容超过最大长度限制。'
    return
  }
  try{
    submitting.value = true
    const res = await uploadCode({
      author: author.value,
      language: language.value,
      content: content.value,
      ttl: ttl.value
    })
    if(res?.hash){
      router.push(`/code/${encodeURIComponent(res.hash)}`)
    }else if(res?.url){
      // 兼容返回 url 的情况：/code/<hash>
      const parts = String(res.url).split('/')
      const h = parts[parts.length-1]
      router.push(`/code/${encodeURIComponent(h)}`)
    }else{
      error.value = '上传成功，但未获得哈希。'
    }
  }catch(e){
    error.value = e?.response?.data?.message || e.message || '上传失败'
  }finally{
    submitting.value = false
  }
}
</script>
