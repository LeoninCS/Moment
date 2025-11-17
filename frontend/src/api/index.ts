import axios from "axios";

const instance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api",
    timeout: 5000,
});

export function uploadCode(data: {
    author: string;
    language: string;
    content: string;
    destroy_time: number;
}) {
    return instance.post("/codeshare/upload", data);
}

export function getCodeByHash(hash: string) {
  return instance.get(`/codeshare/code/${hash}`)
}

export default instance;