#!/bin/bash

# Claude Code 配置管理器 - 构建并打包成 DMG
# 使用方法: ./build-dmg.sh

set -e  # 遇到错误时退出

# 配置变量
APP_NAME="Claude Code 配置管理器"
APP_BUNDLE="claude-code-config.app"
VERSION="1.0.0"
DMG_NAME="Claude-Code-Config-Manager-${VERSION}"
BUILD_DIR="./build/bin"
TEMP_DMG_DIR="./build/dmg-temp"
DIST_DIR="./build/dist"

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

echo_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

echo_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查依赖
check_dependencies() {
    echo_info "检查依赖..."
    
    # 检查 wails
    if ! command -v wails &> /dev/null; then
        echo_error "Wails CLI 未安装. 请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    fi
    
    # 检查 create-dmg
    if ! command -v create-dmg &> /dev/null; then
        echo_error "create-dmg 未安装. 请运行: brew install create-dmg"
        exit 1
    fi
    
    echo_success "依赖检查完成"
}

# 清理旧的构建文件
clean_build() {
    echo_info "清理旧的构建文件..."
    rm -rf "${BUILD_DIR}/${APP_BUNDLE}"
    rm -rf "${TEMP_DMG_DIR}"
    rm -rf "${DIST_DIR}"
    # 清理可能遗留的临时挂载点
    if [ -d "/Volumes/Claude Code 配置管理器" ]; then
        echo_warning "发现遗留的挂载点，尝试清理..."
        diskutil unmount "/Volumes/Claude Code 配置管理器" 2>/dev/null || true
    fi
    mkdir -p "${DIST_DIR}"
    echo_success "清理完成"
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

# 创建临时DMG目录
prepare_dmg_contents() {
    echo_info "准备 DMG 内容..."
    
    # 完全清理临时目录
    rm -rf "${TEMP_DMG_DIR}"
    mkdir -p "${TEMP_DMG_DIR}"
    
    # 复制应用程序到临时目录
    cp -R "${BUILD_DIR}/${APP_BUNDLE}" "${TEMP_DMG_DIR}/"
    
    # 复制许可文件（如果存在）
    if [ -f "LICENSE" ]; then
        cp LICENSE "${TEMP_DMG_DIR}/"
    fi
    
    if [ -f "README.md" ]; then
        cp README.md "${TEMP_DMG_DIR}/"
    fi
    
    echo_success "DMG 内容准备完成"
}

# 创建 DMG
create_dmg_package() {
    echo_info "创建 DMG 包..."
    
    # 设置 DMG 选项
    DMG_OPTIONS=(
        --volname "${APP_NAME}"
        --window-pos 200 120
        --window-size 600 400
        --icon-size 100
        --icon "${APP_BUNDLE}" 300 200
        --hide-extension "${APP_BUNDLE}"
        --hdiutil-quiet
    )
    
    # 如果存在图标文件，添加图标选项
    if [ -f "${BUILD_DIR}/${APP_BUNDLE}/Contents/Resources/iconfile.icns" ]; then
        echo_info "使用应用程序图标"
        DMG_OPTIONS+=(--volicon "${BUILD_DIR}/${APP_BUNDLE}/Contents/Resources/iconfile.icns")
    else
        echo_warning "未找到应用程序图标，使用默认图标"
    fi
    
    # 创建 DMG
    create-dmg \
        "${DMG_OPTIONS[@]}" \
        "${DIST_DIR}/${DMG_NAME}.dmg" \
        "${TEMP_DMG_DIR}/"
    
    if [ $? -eq 0 ]; then
        echo_success "DMG 创建成功: ${DIST_DIR}/${DMG_NAME}.dmg"
    else
        echo_error "DMG 创建失败"
        exit 1
    fi
}

# 清理临时文件
cleanup() {
    echo_info "清理临时文件..."
    rm -rf "${TEMP_DMG_DIR}"
    echo_success "清理完成"
}

# 显示结果
show_results() {
    echo ""
    echo_success "🎉 构建完成!"
    echo ""
    echo_info "构建结果:"
    echo "  应用程序: ${BUILD_DIR}/${APP_BUNDLE}"
    echo "  DMG 包: ${DIST_DIR}/${DMG_NAME}.dmg"
    echo ""
    
    # 显示文件大小
    if [ -f "${DIST_DIR}/${DMG_NAME}.dmg" ]; then
        DMG_SIZE=$(du -h "${DIST_DIR}/${DMG_NAME}.dmg" | cut -f1)
        echo_info "DMG 大小: ${DMG_SIZE}"
    fi
    
    echo_info "你可以运行以下命令打开构建结果:"
    echo "  open ${BUILD_DIR}/${APP_BUNDLE}"
    echo "  open ${DIST_DIR}"
}

# 主函数
main() {
    echo ""
    echo_info "🚀 开始构建 ${APP_NAME}"
    echo_info "版本: ${VERSION}"
    echo ""
    
    check_dependencies
    clean_build
    build_app
    prepare_dmg_contents
    create_dmg_package
    cleanup
    show_results
    
    echo ""
    echo_success "✨ 全部完成!"
}

# 执行主函数
main "$@"