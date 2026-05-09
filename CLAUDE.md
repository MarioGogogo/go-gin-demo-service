# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# 启动服务
go run main.go

# 安装/更新依赖
go mod tidy

# 构建二进制
go build -o server main.go

# 检查代码问题
go vet ./...
```

## Project Architecture

Gin-based Web API service with PostgreSQL backend. No ORM — uses raw SQL via `database/sql` + `pgx`.

### Directory Layout

```
main.go                 # 入口：初始化数据库、注册路由、启动 HTTP 服务
database/database.go    # 数据库初始化、建表、迁移、种子数据
handler/                # Gin handler 函数（按领域拆分）
model/                  # 数据模型 + 统一响应结构
static/                 # 静态资源（可下载文件、chunk bundle）
uploads/modules/        # 模块上传存储目录
web/index.html          # 管理后台前端（单页 HTML）
```

### Key Design Decisions

- **PostgreSQL 直连** — 通过 `database.DB` 全局变量访问 `*sql.DB`，所有 handler 直接写 SQL，无 ORM
- **统一响应格式** — 所有 API 返回 `{"success":bool, "code":int, "message":string, "data":any, "timestamp":int64}`，通过 `model.Ok()`/`model.Fail()` 构建
- **模拟登录** — 硬编码 `admin`/`123456`，返回 fake token，无真实 JWT
- **文件上传** — 支持 `app`（覆盖更新）和 `hotel-module`（同名覆盖）两种类型，自动记录版本历史

### Routes

| 分组 | 路径 | 说明 |
|------|------|------|
| `/users` | GET/POST/PUT/DELETE | 用户 CRUD |
| `/login` | POST | 登录认证 |
| `/files`, `/download/:filename` | GET | 文件列表 & 下载（static/目录）|
| `/app/update` | POST | App 更新检查 |
| `/api/modules` | GET/POST/PUT/DELETE | 模块管理 + 版本历史 + 下载 |
| `/api/chunk-modules` | GET | 分包模块列表 |
| `/admin` | GET | 管理后台页面 |
| `/health` | GET | 健康检查 |

### Database Tables

- `users` — id, name, age
- `modules` — id, name, type, version, file_name, file_path, file_size, changelog, code, download_url, timestamps
- `version_histories` — module_id, version, file_name, file_size, changelog, created_at
- `app_updates` — has_update, version_code, version_name, download_url, changelog, force_update, file_size

### Important Notes

- 服务启动时自动建表和插入种子数据（默认用户 + App 更新记录）
- 下载接口有路径穿越防护（拒绝含 `..` 的文件名）
- 数据库连接信息硬编码在 `database/database.go` 中
- 前端管理页面是单个 HTML 文件 `web/index.html`
