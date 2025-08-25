# Claude Code 配置管理器

一个优雅的桌面应用程序，用于管理 Claude Code 的配置文件和 API 密钥。基于 Wails 框架构建，提供直观的用户界面来管理您的 AI 助手设置。

## ✨ 特性

- 🎨 **现代化界面**: 基于 Vue 3 和 Vite 构建的响应式用户界面
- 🔐 **安全管理**: 安全地存储和管理 Anthropic API 密钥
- 🌐 **多环境配置**: 支持不同的 Base URL 配置
- 📁 **配置文件管理**: 创建、保存和切换多个配置文件
- 💫 **实时预览**: 实时查看配置状态和文件位置
- 🎯 **用户友好**: 直观的界面设计，支持密码显示/隐藏
- 🔍 **状态监控**: 实时显示配置加载状态
- 💾 **自动保存**: 智能保存配置到正确位置

## 🏗️ 技术栈

### 后端
- **Go**: 主要开发语言
- **Wails v2**: 桌面应用框架
- **原生文件系统**: 直接访问用户配置文件

### 前端
- **Vue 3**: 现代化前端框架
- **Vite**: 快速构建工具
- **响应式设计**: 支持各种屏幕尺寸
- **CSS3**: 现代化样式设计

## 📦 安装与运行

### 环境要求

- Go 1.18+ 
- Node.js 16+
- npm 或 yarn

### 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 安装前端依赖
cd frontend
npm install
```

### 开发模式

```bash
# 启动开发服务器
wails dev
```

开发模式会启动 Vite 开发服务器，提供快速的热重载功能。您还可以在浏览器中访问 http://localhost:34115 进行开发调试。

### 构建生产版本

```bash
# 构建可分发的生产版本
wails build
```

## 📁 项目结构

```
claude-code-config/
├── main.go              # 应用程序入口
├── app.go               # 应用程序核心逻辑
├── config_service.go    # 配置文件服务
├── profile_service.go   # 配置文件管理服务
├── models.go            # 数据模型定义
├── wails.json          # Wails 项目配置
├── frontend/           # 前端代码
│   ├── src/
│   │   ├── App.vue              # 主应用组件
│   │   ├── components/
│   │   │   └── ClaudeConfigManager.vue  # 配置管理器
│   │   ├── assets/              # 静态资源
│   │   ├── main.js             # 前端入口
│   │   └── style.css          # 全局样式
│   ├── dist/                  # 构建输出
│   ├── package.json           # 前端依赖
│   └── vite.config.js         # Vite 配置
└── build/                    # 构建资源
```

## 🚀 功能说明

### 配置管理
- **加载配置**: 从 `~/.claude/settings.json` 读取现有配置
- **保存配置**: 将配置保存到正确的文件位置
- **实时验证**: 验证配置格式和必填字段
- **密码保护**: 支持 API 密钥的显示/隐藏切换

### 配置文件管理
- **创建配置文件**: 保存当前配置为可重用的配置文件
- **管理配置文件**: 查看、编辑和删除保存的配置文件
- **快速切换**: 一键应用不同的配置文件
- **元数据支持**: 支持配置文件描述和时间戳

### 用户界面
- **双视图模式**: 当前配置和配置文件管理两个视图
- **响应式设计**: 适配不同屏幕尺寸
- **优雅动画**: 平滑的过渡效果和微交互
- **状态指示**: 清晰的配置状态和错误提示

## ⚙️ 配置文件格式

应用程序管理标准的 Claude Code 配置文件：

```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "sk-ant-api03-xxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.anthropic.com"
  },
  "permissions": {
    "allow": [],
    "deny": []
  }
}
```

## 🔧 开发指南

### 添加新功能

1. **后端功能**: 在 `app.go` 中添加新的方法
2. **数据模型**: 在 `models.go` 中定义数据结构
3. **前端组件**: 在 `components/` 目录下创建新组件
4. **样式设计**: 遵循现有的设计系统

### 代码风格

- **Go**: 遵循 Go 官方代码风格
- **Vue**: 使用 Composition API 和 `<script setup>` 语法
- **CSS**: 使用现代化的 CSS 特性和动画

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 📄 许可证

本项目采用 MIT 许可证。详情请参见 LICENSE 文件。

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的桌面应用框架
- [Vue.js](https://vuejs.org/) - 现代化的前端框架
- [Claude](https://claude.ai/) - 强大的 AI 助手

---

**Made with ❤️ for Claude Code**