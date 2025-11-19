# ========= 构建阶段：打包 Vite =========
FROM node:20-alpine AS build

WORKDIR /app

# 先拷贝依赖描述文件
COPY package*.json ./

# 安装依赖（CI 环境推荐用 npm ci）
RUN npm ci

# 拷贝剩下的前端代码
COPY . .

# 打包
RUN npm run build

# ========= 运行阶段：nginx 提供静态文件 + 反向代理 =========
FROM nginx:alpine

# 拷贝我们自己的 nginx 配置，覆盖默认配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 拷贝打包好的静态资源到 nginx 默认静态目录
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
