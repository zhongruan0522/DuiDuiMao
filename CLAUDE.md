# CLAUDE.md

> 本文件用于存储 Claude Code 的项目级配置与说明

---

## 项目概述

**DuiDuiMao（兑兑猫）** 是一个基于 LinuxDo Connect 登录的公益 CDK 兑换平台。

### 核心功能
- 用户通过 LinuxDo Connect 进行身份认证
- CDK 兑换功能
- 管理员后台管理

---

## 技术栈

### 后端
- **语言**: Go
- **框架**: Gin (推荐，基于标准 HTTP 路由)
- **数据库**: 待定

### 前端
- **框架**: Vue 3
- **UI 组件**: Shadcn-vue
- **构建工具**: Vite
- **样式**: Tailwind CSS

### 部署
- **容器化**: Docker
- **编排**: Docker Compose
- **CI/CD**: GitHub Actions

---

## 项目架构

```
DuiDuiMao/
├── Docs/                          # 项目文档
│   ├── API/                       # API 文档
│   │   └── 相关API文档.md
│   └── 杂七杂八文档或参考项目
│
├── cmd/                           # 应用入口
│   └── server/
│       └── main.go                # 后端服务入口
│
├── internal/                      # 内部包（不对外暴露）
│   ├── handler/                   # HTTP 处理器（路由层）
│   ├── service/                   # 业务逻辑层
│   └── model/                     # 数据模型层
│
├── web/                           # 前端项目
│   ├── src/
│   │   ├── components/            # Vue 组件
│   │   ├── views/                 # 页面视图
│   │   ├── assets/                # 静态资源
│   │   └── main.js                # 前端入口
│   ├── dist/                      # 前端构建产物（本地不上传）
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
│
├── config.example.yaml            # 配置文件模板
├── Dockerfile                     # Docker 镜像构建文件
├── docker-compose.yml             # 生产环境编排（云端镜像）
├── docker-compose.dev.yml         # 开发环境编排（本地构建）
└── go.mod                         # Go 依赖管理
```

---

## 环境变量配置

### Docker 部署环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `adminname` | 管理员账户名 | `root` |
| `adminpassword` | 管理员密码 | `rootpassword` |
| `port` | 服务监听端口 | `3001` |

### 配置文件
- 复制 `config.example.yaml` 为 `config.yaml` 后修改配置
- `config.yaml` 已被 `.gitignore` 屏蔽，不会提交到仓库

---

## 开发规范

### Go 后端规范
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 编码风格
- 使用 `gofmt` 格式化代码
- 错误处理不使用 `panic`，返回 error 类型
- 包名使用小写单词，不使用下划线或驼峰

### Vue 前端规范
- 组件命名使用 PascalCase（如 `UserProfile.vue`）
- 组合式 API（Composition API）优于选项式 API
- 使用 `<script setup>` 语法糖
- 样式优先使用 Tailwind CSS 工具类

### Git 提交规范
- feat: 新功能
- fix: 修复 Bug
- docs: 文档更新
- style: 代码格式调整
- refactor: 重构（不改变功能）
- test: 测试相关
- chore: 构建/工具链相关

---

## 部署流程

### GitHub Actions
1. 自动构建前端产物（`web/dist/`）
2. 构建后端二进制文件
3. 打包成 Docker 镜像并推送

### Docker 部署
```bash
# 生产环境（使用云端构建的镜像）
docker-compose -f docker-compose.yml up -d

# 开发环境（本地构建镜像）
docker-compose -f docker-compose.dev.yml up -d
```

---

## 注意事项

1. **前端构建**: 本地开发时不需要手动构建前端，GitHub Actions 会自动构建
2. **配置安全**: `config.yaml` 包含敏感信息，已被 `.gitignore` 屏蔽
3. **Docker 编排**: 提供两个 yml 文件，分别用于生产环境和开发环境
