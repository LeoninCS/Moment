# Moment Vue 前端

一个用于对接 **Moment / CodeShare** 服务的 Vue 3 + Vite 前端。

## 功能
- 顶部栏：左侧 `Moment` 返回首页，右侧 GitHub 图标链接到项目仓库
- 首页：功能卡片网格，目前开放 **代码分享**
- 代码分享：填写作者、语言、保存时间、代码内容，上传后根据返回的 `hash` 跳转到 `/code/:hash`
- 代码查看：根据 `hash` 只读展示代码、作者、语言、销毁时间（如有），使用 Prism 高亮

## 运行
```bash
npm i
npm run dev
```

默认开发环境网关地址为 `http://localhost:8080`（见 `.env.development`）。
如有跨域，请在后端 / 网关开启 CORS。

## 环境变量
- `VITE_API_BASE`：API 网关基础地址（默认 `http://localhost:8080`）

## 与后端的接口约定
- 上传：`POST /CodeShare/upload`，body: `{ author, language, content, ttl }`，返回 `{ message, hash, url }`
- 查询：`GET /CodeShare/code/{hash}`，返回 `{ author, language, content, hash, destroy_time }`

## 构建
```bash
npm run build
npm run preview
```
