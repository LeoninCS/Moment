<template>
  <div class="httptest-page">
    <header class="hero">
      <div>
        <p class="eyebrow">HTTP Test Lab</p>
        <h1>在线接口测试（类 Postman）</h1>
        <p class="sub">
          构造和保存多个请求，点击发送即可调用后端 / 前端接口。工作区支持 JSON 导入 / 导出，本地自动保存。
        </p>
        <div class="hero-actions">
          <button class="btn primary" @click="addRequest">+ 新请求</button>
          <button class="btn ghost" :disabled="!requests.length" @click="exportWorkspace">
            导出 JSON
          </button>
          <label class="btn ghost import-btn">
            导入 JSON
            <input ref="fileInput" type="file" accept="application/json" @change="handleImportFile" />
          </label>
          <span class="muted">数据会自动存到浏览器，不怕刷新</span>
        </div>
      </div>
    </header>

    <div class="workspace">
      <aside class="request-list">
        <div class="list-header">
          <div>
            <div class="pill">Workspace</div>
            <p class="list-sub">维护多个接口配置，点击切换</p>
          </div>
          <button class="icon-btn" title="新建请求" @click="addRequest">＋</button>
        </div>

        <div class="list-body">
          <div
            v-for="req in requests"
            :key="req.id"
            class="req-item"
            :class="{ active: req.id === activeId }"
          >
            <div class="req-main" @click="selectRequest(req.id)">
              <span class="method" :data-method="req.method">{{ req.method }}</span>
              <div class="req-text">
                <input
                  v-model="req.name"
                  class="req-name"
                  placeholder="请求名称"
                  @click.stop
                />
                <p class="req-url">{{ req.url || "未填写 URL" }}</p>
              </div>
            </div>
            <div class="req-actions">
              <button class="icon-btn" title="复制" @click.stop="duplicateRequest(req.id)">⧉</button>
              <button
                class="icon-btn danger"
                title="删除"
                @click.stop="deleteRequest(req.id)"
              >
                ✕
              </button>
            </div>
          </div>
        </div>
      </aside>

      <section class="request-area">
        <div v-if="activeRequest" class="request-card">
          <div class="request-row">
            <select v-model="activeRequest.method" class="method-select">
              <option v-for="m in methods" :key="m" :value="m">{{ m }}</option>
            </select>

            <input
              v-model="activeRequest.url"
              class="url-input"
              placeholder="https://api.example.com/v1/users"
            />

            <button class="btn primary" :disabled="sending" @click="sendActiveRequest">
              {{ sending ? "发送中..." : "发送" }}
            </button>
          </div>

          <div class="request-row secondary">
            <label class="stacked">
              <span>请求名称</span>
              <input v-model="activeRequest.name" placeholder="如：获取用户列表" />
            </label>
            <div class="row-actions">
              <button class="btn ghost small" @click="addHeaderRow">+ Header</button>
              <button class="btn ghost small" @click="prefillJsonTemplate">
                填充 JSON 模板
              </button>
            </div>
          </div>

          <div class="panel">
            <div class="panel-header">
              <div>
                <h3>Headers</h3>
                <p class="hint">可勾选启用，空行自动忽略</p>
              </div>
            </div>
            <div class="headers">
              <div v-for="header in activeRequest.headers" :key="header.id" class="header-row">
                <input v-model="header.enabled" type="checkbox" />
                <input v-model="header.key" placeholder="Header Key" />
                <input v-model="header.value" placeholder="Header Value" />
                <button class="icon-btn danger" title="移除" @click="removeHeaderRow(header.id)">
                  ✕
                </button>
              </div>
            </div>
          </div>

          <div class="panel">
            <div class="panel-header">
              <div>
                <h3>Body</h3>
                <p class="hint">POST / PUT / PATCH 建议使用 JSON；GET 默认忽略请求体</p>
              </div>
            </div>
            <textarea
              v-model="activeRequest.body"
              class="body-input"
              rows="10"
              placeholder='{"name":"DevDesk"}'
            ></textarea>
          </div>

          <div class="panel response-panel">
            <div class="panel-header">
              <div>
                <h3>响应</h3>
                <p class="hint">状态码、响应头与正文会显示在这里</p>
              </div>
              <div v-if="latestResponse" class="resp-meta">
                <span class="badge" :class="statusBadgeClass">
                  {{ latestResponse.status }}
                </span>
                <span class="muted">
                  耗时 {{ latestResponse.durationMs }} ms
                  <span v-if="responseSize"> · {{ responseSize }}</span>
                </span>
              </div>
            </div>

            <p v-if="sendError" class="error-text">{{ sendError }}</p>

            <div v-if="latestResponse" class="response-content">
              <div class="resp-block">
                <div class="resp-title">Body</div>
                <pre class="resp-body">{{ prettyBody }}</pre>
              </div>
              <div class="resp-block headers">
                <div class="resp-title">Headers</div>
                <div v-for="(vals, key) in latestResponse.header" :key="key" class="resp-header">
                  <strong>{{ key }}:</strong>
                  <span>{{ vals.join(", ") }}</span>
                </div>
              </div>
            </div>

            <div v-else class="placeholder">
              <p>还没有响应，填写 URL 后点击发送即可。</p>
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <p>当前工作区没有请求，点击下方按钮创建。</p>
          <button class="btn primary" @click="addRequest">新建请求</button>
        </div>
      </section>
    </div>

    <p v-if="importError" class="error-banner">
      {{ importError }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { sendHttpTest, type HttpTestResponse } from "../api/httptest";

type HeaderRow = {
  id: string;
  key: string;
  value: string;
  enabled: boolean;
};

type WorkspaceRequest = {
  id: string;
  name: string;
  method: string;
  url: string;
  body: string;
  headers: HeaderRow[];
  lastResponse?: AugmentedResponse;
};

type AugmentedResponse = HttpTestResponse & { durationMs: number };

const methods = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"];
const storageKey = "devdesk-httptest-workspace";

const requests = ref<WorkspaceRequest[]>([]);
const activeId = ref<string>("");
const latestResponse = ref<AugmentedResponse | null>(null);
const sending = ref(false);
const sendError = ref("");
const importError = ref("");
const fileInput = ref<HTMLInputElement | null>(null);

const activeRequest = computed(
  () => requests.value.find((r) => r.id === activeId.value) || null
);

const prettyBody = computed(() => {
  if (!latestResponse.value) return "";
  const raw = latestResponse.value.body || "";
  try {
    return JSON.stringify(JSON.parse(raw), null, 2);
  } catch (_) {
    return raw;
  }
});

const responseSize = computed(() => {
  if (!latestResponse.value) return "";
  const len = latestResponse.value.body?.length || 0;
  if (len > 1024) return `${(len / 1024).toFixed(1)} KB`;
  return `${len} B`;
});

const statusBadgeClass = computed(() => {
  if (!latestResponse.value) return "";
  const code = latestResponse.value.status_code;
  if (code >= 200 && code < 300) return "badge-green";
  if (code >= 300 && code < 500) return "badge-warn";
  return "badge-red";
});

function genId() {
  return crypto.randomUUID ? crypto.randomUUID() : `id-${Date.now()}-${Math.random()}`;
}

function createHeaderRow(seed?: Partial<HeaderRow>): HeaderRow {
  return {
    id: genId(),
    key: seed?.key || "",
    value: seed?.value || "",
    enabled: seed?.enabled ?? true,
  };
}

function createRequest(name?: string): WorkspaceRequest {
  return {
    id: genId(),
    name: name || `请求 ${requests.value.length + 1}`,
    method: "GET",
    url: "",
    body: "",
    headers: [createHeaderRow()],
  };
}

function selectRequest(id: string) {
  activeId.value = id;
  const req = requests.value.find((r) => r.id === id);
  latestResponse.value = req?.lastResponse || null;
  sendError.value = "";
}

function addRequest() {
  const req = createRequest();
  requests.value.push(req);
  selectRequest(req.id);
}

function duplicateRequest(id: string) {
  const source = requests.value.find((r) => r.id === id);
  if (!source) return;
  const copy: WorkspaceRequest = {
    ...source,
    id: genId(),
    name: `${source.name} 副本`,
    headers: source.headers.map((h) => ({ ...h, id: genId() })),
    lastResponse: undefined,
  };
  requests.value.push(copy);
  selectRequest(copy.id);
}

function deleteRequest(id: string) {
  const idx = requests.value.findIndex((r) => r.id === id);
  if (idx === -1) return;
  requests.value.splice(idx, 1);
  if (requests.value.length === 0) {
    activeId.value = "";
    latestResponse.value = null;
  } else if (activeId.value === id) {
    const next = requests.value[Math.max(0, idx - 1)];
    if (next) {
      selectRequest(next.id);
    }
  }
}

function addHeaderRow() {
  if (!activeRequest.value) return;
  activeRequest.value.headers.push(createHeaderRow());
}

function removeHeaderRow(headerId: string) {
  if (!activeRequest.value) return;
  activeRequest.value.headers = activeRequest.value.headers.filter(
    (h) => h.id !== headerId
  );
  if (activeRequest.value.headers.length === 0) {
    activeRequest.value.headers.push(createHeaderRow());
  }
}

function prefillJsonTemplate() {
  if (!activeRequest.value) return;
  if (activeRequest.value.body.trim()) return;
  activeRequest.value.body = JSON.stringify(
    { message: "hello world", timestamp: Date.now() },
    null,
    2
  );
}

function normalizeRequest(data: any): WorkspaceRequest | null {
  if (!data || typeof data !== "object") return null;
  const headers = Array.isArray(data.headers)
    ? data.headers.map((h: any) => createHeaderRow(h))
    : [createHeaderRow()];
  return {
    id: data.id || genId(),
    name: data.name || "未命名请求",
    method: data.method && methods.includes(data.method) ? data.method : "GET",
    url: data.url || "",
    body: data.body || "",
    headers,
  };
}

function persistWorkspace() {
  const payload = {
    activeId: activeId.value,
    requests: requests.value,
  };
  try {
    localStorage.setItem(storageKey, JSON.stringify(payload));
  } catch (e) {
    console.warn("保存工作区失败", e);
  }
}

watch(
  () => [requests.value, activeId.value],
  () => {
    persistWorkspace();
  },
  { deep: true }
);

function loadWorkspace() {
  const raw = localStorage.getItem(storageKey);
  if (!raw) return;
  try {
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed.requests)) {
      const list = parsed.requests
        .map((r: any) => normalizeRequest(r))
        .filter(Boolean) as WorkspaceRequest[];
      requests.value = list;
      const firstId = parsed.activeId || list[0]?.id || "";
      if (firstId) selectRequest(firstId);
    }
  } catch (e) {
    console.warn("读取工作区失败", e);
  }
}

function headersToMap(headers: HeaderRow[]) {
  return headers.reduce<Record<string, string>>((acc, h) => {
    if (h.enabled && h.key.trim()) {
      acc[h.key.trim()] = h.value;
    }
    return acc;
  }, {});
}

async function sendActiveRequest() {
  if (!activeRequest.value) return;
  if (!activeRequest.value.url.trim()) {
    sendError.value = "请填写完整的 URL";
    return;
  }
  sending.value = true;
  sendError.value = "";
  latestResponse.value = null;

  const payload = {
    method: activeRequest.value.method,
    url: activeRequest.value.url,
    body: activeRequest.value.body,
    header: headersToMap(activeRequest.value.headers),
  };

  const started = performance.now();
  try {
    const res = await sendHttpTest(payload);
    const duration = Math.round(performance.now() - started);
    const merged: AugmentedResponse = {
      ...res.data,
      durationMs: duration,
    };
    latestResponse.value = merged;
    activeRequest.value.lastResponse = merged;
  } catch (err: any) {
    sendError.value =
      err?.response?.data?.error || err?.message || "请求失败，请稍后再试";
  } finally {
    sending.value = false;
  }
}

function exportWorkspace() {
  const payload = {
    version: 1,
    activeId: activeId.value,
    requests: requests.value,
  };
  const blob = new Blob([JSON.stringify(payload, null, 2)], {
    type: "application/json",
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `devdesk-httptest-${Date.now()}.json`;
  a.click();
  URL.revokeObjectURL(url);
}

function handleImportFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  importError.value = "";
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => {
    try {
      const parsed = JSON.parse(String(reader.result || "{}"));
      if (!Array.isArray(parsed.requests)) {
        throw new Error("File format error: requests is missing");
      }
      const list = parsed.requests
        .map((r: any) => normalizeRequest(r))
        .filter(Boolean) as WorkspaceRequest[];
      if (!list.length) throw new Error("No request data found");
      requests.value = list;
      const firstId = parsed.activeId || list[0]?.id;
      if (firstId) {
        selectRequest(firstId);
      } else {
        activeId.value = "";
        latestResponse.value = null;
      }
    } catch (err: any) {
      importError.value = err?.message || "Import failed, please check JSON";
    } finally {
      if (fileInput.value) fileInput.value.value = "";
    }
  };
  reader.readAsText(file);
}

onMounted(() => {
  loadWorkspace();
  if (!requests.value.length) {
    addRequest();
  }
});
</script>

<style scoped>
.httptest-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero {
  background: linear-gradient(120deg, #0ea5e9 0%, #2563eb 50%, #111827 100%);
  color: #f8fafc;
  border-radius: 14px;
  padding: 22px 20px;
}

.eyebrow {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 12px;
  opacity: 0.8;
}

.hero h1 {
  margin: 6px 0;
  font-size: 24px;
}

.hero .sub {
  margin: 0 0 12px;
  max-width: 800px;
  line-height: 1.5;
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.workspace {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 16px;
}

.request-list {
  background: var(--sidebar);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 600px;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.pill {
  display: inline-block;
  background: rgba(37, 99, 235, 0.12);
  color: #1d4ed8;
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.list-sub {
  margin: 4px 0 0;
  color: #6b7280;
  font-size: 12px;
}

.list-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
}

.req-item {
  border: 1px solid var(--border);
  border-radius: 10px;
  background: #fff;
  padding: 10px;
  transition: border-color 0.2s, box-shadow 0.2s;
  display: flex;
  justify-content: space-between;
  gap: 8px;
}

:global([data-theme="dark"]) .req-item {
  background: #0f172a;
}

.req-item.active {
  border-color: #2563eb;
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.18);
}

.req-main {
  display: flex;
  gap: 10px;
  flex: 1;
  cursor: pointer;
}

.req-actions {
  display: flex;
  gap: 4px;
  align-items: flex-start;
}

.req-text {
  flex: 1;
  min-width: 0;
}

.req-name {
  width: 100%;
  border: none;
  background: transparent;
  font-weight: 600;
  font-size: 14px;
}

.req-name:focus {
  outline: none;
}

.req-url {
  margin: 2px 0 0;
  color: #6b7280;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.method {
  min-width: 54px;
  text-align: center;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.2px;
  color: #0f172a;
  background: #e2e8f0;
}

.method[data-method="POST"],
.method[data-method="PUT"],
.method[data-method="PATCH"] {
  background: #fee2e2;
  color: #991b1b;
}

.method[data-method="DELETE"] {
  background: #ffe4e6;
  color: #be123c;
}

.method[data-method="GET"] {
  background: #dcfce7;
  color: #166534;
}

.request-area {
  background: #fff;
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 14px;
}

:global([data-theme="dark"]) .request-area,
:global([data-theme="dark"]) .request-list {
  background: #0f172a;
}

.request-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.request-row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.request-row.secondary {
  justify-content: space-between;
  align-items: flex-end;
}

.stacked {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  font-size: 12px;
  color: #6b7280;
}

.stacked input {
  padding: 8px;
  border: 1px solid var(--border);
  border-radius: 8px;
}

.row-actions {
  display: flex;
  gap: 8px;
}

.method-select {
  width: 120px;
  padding: 10px;
  border-radius: 8px;
  border: 1px solid var(--border);
  font-weight: 700;
}

.url-input {
  flex: 1;
  padding: 10px;
  border-radius: 8px;
  border: 1px solid var(--border);
}

.panel {
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.95);
}

:global([data-theme="dark"]) .panel {
  background: #0b1222;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
}

.hint {
  margin: 4px 0 0;
  color: #9ca3af;
  font-size: 12px;
}

.headers {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 10px;
}

.header-row {
  display: grid;
  grid-template-columns: 24px 1fr 1fr 32px;
  gap: 6px;
  align-items: center;
}

.header-row input[type="text"],
.header-row input[type="checkbox"] {
  padding: 8px;
  border: 1px solid var(--border);
  border-radius: 8px;
}

.header-row input[type="checkbox"] {
  width: 20px;
  height: 20px;
}

.body-input {
  width: 100%;
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 10px;
  font-family: "JetBrains Mono", Consolas, monospace;
  background: #0b1222;
  color: #e5e7eb;
  resize: vertical;
}

.response-panel {
  background: #0b172b;
  color: #e5e7eb;
  border-color: #172238;
}

.resp-meta {
  display: flex;
  gap: 10px;
  align-items: center;
}

.response-content {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 12px;
  margin-top: 12px;
}

.resp-block {
  border: 1px solid #1f2c45;
  border-radius: 10px;
  padding: 10px;
  background: #0f1c33;
}

.resp-title {
  font-weight: 600;
  margin-bottom: 6px;
}

.resp-body {
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  font-family: "JetBrains Mono", Consolas, monospace;
  font-size: 13px;
}

.resp-header {
  display: flex;
  gap: 6px;
  font-size: 13px;
  line-height: 1.5;
}

.placeholder {
  margin-top: 10px;
  color: #cbd5e1;
}

.empty-state {
  border: 1px dashed var(--border);
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  color: #6b7280;
}

.btn {
  border: none;
  cursor: pointer;
  padding: 10px 14px;
  border-radius: 8px;
  font-weight: 600;
  transition: transform 0.1s ease, box-shadow 0.2s ease, opacity 0.2s;
}

.btn.small {
  padding: 6px 10px;
}

.btn.primary {
  background: #2563eb;
  color: #fff;
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.35);
}

.btn.primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  box-shadow: none;
}

.btn.ghost {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text);
}

.hero .btn.ghost {
  border-color: rgba(255, 255, 255, 0.5);
  color: #f8fafc;
}

.import-btn {
  position: relative;
}

.import-btn input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.icon-btn {
  border: 1px solid var(--border);
  background: #fff;
  color: #0f172a;
  border-radius: 8px;
  width: 32px;
  height: 32px;
  cursor: pointer;
}

:global([data-theme="dark"]) .icon-btn {
  background: #0f172a;
  color: #e5e7eb;
}

.icon-btn.danger {
  color: #dc2626;
}

.badge {
  padding: 4px 8px;
  border-radius: 8px;
  font-weight: 700;
}

.badge-green {
  background: #dcfce7;
  color: #166534;
}

.badge-warn {
  background: #fef9c3;
  color: #854d0e;
}

.badge-red {
  background: #fee2e2;
  color: #991b1b;
}

.muted {
  color: #cbd5e1;
  font-size: 13px;
}

.error-text {
  color: #fca5a5;
  margin: 8px 0;
}

.error-banner {
  margin: 0;
  padding: 10px 12px;
  border-radius: 10px;
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fecdd3;
}

@media (max-width: 1080px) {
  .workspace {
    grid-template-columns: 1fr;
  }

  .request-list {
    min-height: auto;
  }
}

@media (max-width: 640px) {
  .request-row {
    flex-direction: column;
    align-items: stretch;
  }

  .method-select {
    width: 100%;
  }

  .row-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .response-content {
    grid-template-columns: 1fr;
  }
}
</style>
