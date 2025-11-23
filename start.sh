#!/usr/bin/env bash
set -euo pipefail

# -------- 杀掉占用端口的进程（可选）--------
echo "清理占用端口的进程..."
for p in 8080; do
  pid=$(sudo lsof -t -i:$p || true)
  if [[ -n "${pid}" ]]; then
    sudo kill -9 $pid || true
  fi
done

# -------- 启动后端 --------
echo "启动后端..."
cd backend
(go run cmd/main.go) &     GW_PID=$!


# -------- 启动前端（Vite）--------
(
  cd ../frontend
  npm install
  npm run dev -- --host 0.0.0.0 --port 5173
) & FE_PID=$!

echo
echo "✅ 已启动："
echo "  - 后端：     http://localhost:8080"
echo "  - 前端：     http://localhost:5173"
echo
echo "按 Ctrl+C 可停止所有服务。"

cleanup() {
  echo
  echo "⏹ 停止所有服务..."
  kill $GW_PID  $FE_PID 2>/dev/null || true
  wait 2>/dev/null || true
  echo "Done."
}
trap cleanup INT TERM

wait

# 运行方式
# chmod +x start.sh
# ./start.sh (在项目根目录下运行，即Moment目录)
# Ctrl+C 组合键可停止所有服务