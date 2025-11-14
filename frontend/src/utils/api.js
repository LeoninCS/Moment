// src/utils/api.js
import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

export const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 15000
})

/**
 * CodeShare 相关
 * Go 后端：
 *  cg := r.Group("/codeshare")
 *    cg.POST("/upload", ...)
 *    cg.GET("/code/:hash", ...)
 */
export async function uploadCode({ author, language, content, ttl }) {
  // ⚠️ 一定要是 /codeshare（全小写），不要 /CodeShare
  const { data } = await api.post('/codeshare/upload', {
    author,
    language,
    content,
    ttl
  })
  // 预期：{ message, hash, url } 或 { hash }
  return data
}

export async function getCode(hash) {
  const { data } = await api.get(
    `/codeshare/code/${encodeURIComponent(hash)}`
  )
  // 预期：{ author, language, content, hash, destroy_time }
  return data
}

/**
 * TextCrypto 相关
 * Go 后端：
 *  tg := r.Group("/textcrypto")
 *    tg.GET("/random-screat-key", ...)
 *    tg.POST("/encrypt-text", ...)
 *    tg.POST("/decrypt-text", ...)
 */

// 获取随机密钥
export async function getRandomSecretKey() {
  const { data } = await api.get('/textcrypto/random-screat-key')
  if (typeof data === 'string') return data.trim()
  return (data.screat_key || data.key || '').trim()
}

// 加密文本
export async function encryptText({ text, screat_key }) {
  const { data } = await api.post('/textcrypto/encrypt-text', {
    text,
    screat_key
  })
  return data
}

// 解密文本
export async function decryptText({ encrypted_text, screat_key }) {
  const { data } = await api.post('/textcrypto/decrypt-text', {
    encrypted_text,
    screat_key
  })
  return data
}
