// src/api/http.ts
import axios from "axios";

const baseURL =
  import.meta.env.VITE_API_BASE_URL
  || (import.meta.env.DEV ? "http://localhost:8080/api" : "/api");

const http = axios.create({
  baseURL,
  timeout: 5000,
});

export default http;
