// src/utils/api.js
import axios from 'axios'

// 统一变量名：VITE_API_BASE_URL
const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 15000
})

// ========== CodeShare API ==========
export async function uploadCode({ author, language, content, ttl }) {
  const { data } = await api.post('/codeshare/upload', {
    author,
    language,
    content,
    ttl
  })
  return data
}

export async function getCode(hash) {
  const { data } = await api.get(`/codeshare/code/${encodeURIComponent(hash)}`)
  return data
}

// ========== TextCrypto API ==========
export async function getRandomSecretKey() {
  const { data } = await api.get('/textcrypto/random-screat-key')
  if (typeof data === 'string') return data.trim()
  return (data.screat_key || data.key || '').trim()
}

export async function encryptText({ text, screat_key }) {
  const { data } = await api.post('/textcrypto/encrypt-text', {
    text,
    screat_key
  })
  return data
}

export async function decryptText({ encrypted_text, screat_key }) {
  const { data } = await api.post('/textcrypto/decrypt-text', {
    encrypted_text,
    screat_key
  })
  return data
}
