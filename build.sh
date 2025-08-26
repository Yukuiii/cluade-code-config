#!/bin/bash

# Claude Code 配置管理器 - 快速构建脚本 (仅构建，不打包DMG)
# 使用方法: ./build.sh

set -e  # 遇到错误时退出

# 配置变量
APP_NAME="Claude Code 配置管理器"
APP_BUNDLE="claude-code-config.app"
BUILD_DIR="./build/bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

echo_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查 Wails
check_wails() {
    echo_info "检查 Wails CLI..."
    
    if ! command -v wails &> /dev/null; then
        echo_error "Wails CLI 未安装. 请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    fi
    
    echo_success "Wails CLI 检查完成"
}

# 构建应用程序
build_app() {
    echo_info "构建应用程序..."
    wails build
    
    if [ ! -d "${BUILD_DIR}/${APP_BUNDLE}" ]; then
        echo_error "构建失败: 应用程序包未找到"
        exit 1
    fi
    
    echo_success "应用程序构建完成"
}

# 打开应用程序
open_app() {
    echo_info "打开应用程序..."
    open "${BUILD_DIR}/${APP_BUNDLE}"
}

# 主函数
main() {
    echo ""
    echo_info "🚀 快速构建 ${APP_NAME}"
    echo ""
    
    check_wails
    build_app
    
    echo ""
    echo_success "✨ 构建完成!"
    echo_info "应用程序位置: ${BUILD_DIR}/${APP_BUNDLE}"
    echo ""
    
    # 询问是否打开应用程序
    read -p "是否立即打开应用程序? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open_app
    fi
}

# 执行主函数
main "$@"