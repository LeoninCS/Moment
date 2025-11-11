import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

export const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 15000
})

export async function uploadCode({ author, language, content, ttl }) {
  const { data } = await api.post('/CodeShare/upload', { author, language, content, ttl })
  return data // { message, hash, url }
}

export async function getCode(hash) {
  const { data } = await api.get(`/CodeShare/code/${encodeURIComponent(hash)}`)
  return data // { author, language, content, hash, destroy_time }
}
