#!/usr/bin/env bash
set -euo pipefail

# -------- Redis --------
echo "[1/4] 启动 Redis..."
if ! redis-cli ping >/dev/null 2>&1; then
  sudo service redis-server start
fi
redis-cli ping

# -------- 环境变量 --------
echo "[2/4] 设置 Redis 环境变量..."
export REDIS_ADDR=127.0.0.1:6379
# export REDIS_PASSWORD=xxxxx
# export REDIS_DB=0

# -------- 杀掉占用端口的进程（可选）--------
echo "[3/4] 清理占用端口的进程..."
for p in 8080 8081 8082 8083; do
  pid=$(sudo lsof -t -i:$p || true)
  if [[ -n "${pid}" ]]; then
    sudo kill -9 $pid || true
  fi
done

# -------- 启动后端 --------
echo "[4/4] 启动后端（网关 + CodeShare 实例）..."
cd backend
(go run cmd/main.go) &          GW_PID=$!
(go run cmd/CodeShare/main.go 8081) & CS1_PID=$!
(go run cmd/CodeShare/main.go 8082) & CS2_PID=$!
(go run cmd/CodeShare/main.go 8083) & CS3_PID=$!

# -------- 启动前端（Vite）--------
(
  cd ../frontend
  npm install
  npm run dev
) & FE_PID=$!

echo
echo "✅ 已启动："
echo "  - 网关：     http://localhost:8080"
echo "  - 前端：     http://localhost:5173"
echo "  - CodeShare：8081 / 8082 / 8083"
echo
echo "按 Ctrl+C 可停止所有服务。"

cleanup() {
  echo
  echo "⏹ 停止所有服务..."
  kill $GW_PID $CS1_PID $CS2_PID $CS3_PID $FE_PID 2>/dev/null || true
  wait 2>/dev/null || true
  echo "Done."
}
trap cleanup INT TERM

wait

# 运行方式
# chmod +x start.sh
# ./start.sh (在项目根目录下运行，即Moment目录)
# Ctrl+C 组合键可停止所有服务