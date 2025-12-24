# DuiDuiMao（兑兑猫）

> 一个基于LinuxDo Connect登录的公益CDK兑换平台

## 项目架构

> 仅供参考

project/
├── Docs/
│   └── API/
│       └── 相关API文档.md
│   └── 杂七杂八文档或参考项目
│
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── model/
│
├── web/
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── assets/
│   │   └── main.js
│   ├── dist/                    # 前端构建产物
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
│
├── config.example.yaml
├── Dockerfile
├── docker-compose.yml
├── docker-compose.dev.yml
└── go.mod

## 项目包发布

1. 使用GitHub Actions打包成Docker镜像
2. 一个Docker即可运行项目内容

### GitHub Actions

1. 在GitHub Actions运行的时候构建前端，本地构建的不上传云端
2. 做好代码质量检查

### Docker部署注意事项

1. 需要提供两个Yml文件，一个使用云端构建好的镜像，一个是供开发者本地构建的镜像

## 项目技术栈

1. 后端:GO语言
2. 前端：Vue + Shadcn-vue

## 环境变量配置参数

> 此处指的是docker部署Yml里面的env

1. adminname：管理员账户名，未配置使用`root`
2. adminpassword：管理员密码，未配置使用`rootpassword`
3. port：端口，默认为3001