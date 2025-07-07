#!/bin/bash

# 测试前端构建脚本

set -e

echo "🚀 开始测试前端构建..."

# 进入web目录
cd web

# 清理旧的构建
echo "🧹 清理旧的构建..."
rm -rf build
mkdir -p build

# 构建default
echo "🔨 构建 default..."
cd default
npm install --legacy-peer-deps
DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ../../VERSION) npm run build
cd ..

# 构建berry  
echo "🔨 构建 berry..."
cd berry
npm install --legacy-peer-deps
DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ../../VERSION) npm run build
cd ..

# 构建air
echo "🔨 构建 air..."
cd air
npm install --legacy-peer-deps  
DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ../../VERSION) npm run build
cd ..

# 检查构建结果
echo "📋 构建结果："
ls -la build/

echo "✅ 前端构建完成！"
