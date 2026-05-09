# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

基于 Gin 框架的 Go 后端服务，提供用户管理、模块/分包管理、文件上传下载、App 更新检查等功能。使用 PostgreSQL 数据库（pgx/v5 驱动）。

## 常用命令

```bash
# 启动服务（监听 0.0.0.0:8080）
go run main.go

# 编译
go build -o server main.go
```


# 编译
GO_ENV=production go build -o server main.go

# 运行
GO_ENV=production ./server
或者编译时不嵌入环境变量，运行时再指定（推荐，同一个二进制文件可灵活切换环境）：


# 编译（不区分环境）
go build -o server main.go

//宝塔面板 Linux 服务器一般是 linux/amd64 架构，而你的 Mac 是 darwin/amd64 或 darwin/arm64。需要交叉编译
GOOS=linux GOARCH=amd64 go build -o server main.go

# 本地运行
./server


# 线上运行
chmod +x server
GO_ENV=production ./server


本项目没有测试文件、Makefile、Dockerfile 或 CI/CD 配置。

## 架构

**扁平三层结构**，无 service/repository 层：

- **main.go** — 路由注册，所有路由直接在 `main()` 中通过 Gin router group 定义
- **handler/** — 处理函数，直接内联 SQL 操作全局 `database.DB`（`*sql.DB`）
- **model/** — 数据结构定义，含 JSON/binding tag；`response.go` 提供统一响应封装（`Ok()`/`Fail()`）
- **database/** — 数据库初始化，包含建表、迁移（ALTER TABLE）、种子数据
- **web/** — 管理后台单页应用（`index.html`），通过 `/admin` 路由提供

**路由分组**：`/users`（用户CRUD）、`/login`、`/files`+`/download/:filename`（静态文件）、`/api/modules`（模块管理+上传）、`/api/chunk-modules`、`/app/update`、`/admin`、`/health`

**数据库**：启动时自动执行建表和迁移。四张表：`users`、`modules`、`version_histories`、`app_updates`。

**无中间件**：除 Gin 默认的 Logger+Recovery 外无任何自定义中间件（无认证、无 CORS）。

## 注意事项

- README.md 中"无需数据库、内存存储"的描述已过时，项目实际使用 PostgreSQL
- 数据库连接、服务端口、登录凭据等均硬编码在源码中，无配置文件机制
- 文件上传目录为 `./uploads/modules`，静态文件目录为 `./static`
