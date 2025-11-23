// src/api/markdown.ts
import http from "./http";

export interface MarkdownNewResponse {
  hash: string;
}

export interface MarkdownDocResponse {
  hash: string;
  content: string;
}

// GET /markdown/new  创建新文档
export function createMarkdownDoc() {
  return http.get<MarkdownNewResponse>("/markdown/new");
}

// GET /markdown/:hash  获取文档
export function fetchMarkdownDoc(hash: string) {
  return http.get<MarkdownDocResponse>(`/markdown/${hash}`);
}

// POST /markdown/update  更新文档
// body: { hash, content }
export function updateMarkdownDoc(hash: string, content: string) {
  return http.post("/markdown/update", { hash, content });
}
