# go-gin-demo-service

基于 [Gin](https://github.com/gin-gonic/gin) 框架的 Go 语言 Web API 示例服务，包含用户管理、登录认证和文件下载功能。

## 环境要求

- Go 1.26+
- 无需数据库，数据存储在内存中

## 快速开始

```bash
# 克隆项目
git clone <仓库地址>
cd go-gin-demo-service

# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

服务启动后访问：

- 本机：`http://localhost:8080`
- 局域网：`http://<你的IP>:8080`

## 项目结构

```
go-gin-demo-service/
├── main.go              # 入口文件，路由注册
├── handler/
│   ├── user.go          # 用户 CRUD 处理
│   ├── auth.go          # 登录认证处理
│   └── download.go      # 文件下载处理
├── model/
│   ├── user.go          # User 模型
│   └── response.go      # 统一响应结构与状态码
├── static/              # 静态资源目录（可下载文件）
├── go.mod
└── go.sum
```

## 统一响应格式

所有接口返回统一的 JSON 结构：

```json
{
  "success": true,
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1714800000
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| success | bool | 请求是否成功 |
| code | int | 业务状态码（0=成功，400=参数错误，401=未授权，404=未找到，500=服务器错误） |
| message | string | 提示信息 |
| data | any | 返回数据 |
| timestamp | int64 | Unix 时间戳 |

---

## 接口文档

### 健康检查

```
GET /health
```

**响应示例：**

```json
{
  "status": "ok"
}
```

---

### 登录

```
POST /login
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例：**

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

**成功响应：**

```json
{
  "success": true,
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "fake-token-admin",
    "username": "admin"
  },
  "timestamp": 1714800000
}
```

> 默认账号：`admin` / `123456`

---

### 用户管理

#### 获取用户列表

```
GET /users
```

```bash
curl http://localhost:8080/users
```

**响应示例：**

```json
{
  "success": true,
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      { "id": 1, "name": "张三", "age": 25 },
      { "id": 2, "name": "李四", "age": 30 }
    ],
    "total": 2
  },
  "timestamp": 1714800000
}
```

#### 获取单个用户

```
GET /users/:id
```

```bash
curl http://localhost:8080/users/1
```

**响应示例：**

```json
{
  "success": true,
  "code": 0,
  "message": "success",
  "data": { "id": 1, "name": "张三", "age": 25 },
  "timestamp": 1714800000
}
```

#### 创建用户

```
POST /users
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 姓名 |
| age | int | 是 | 年龄（0-150） |

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"王五","age":28}'
```

**响应示例：**

```json
{
  "success": true,
  "code": 0,
  "message": "创建成功",
  "data": { "id": 3, "name": "王五", "age": 28 },
  "timestamp": 1714800000
}
```

#### 更新用户

```
PUT /users/:id
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 姓名 |
| age | int | 是 | 年龄（0-150） |

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"张三丰","age":100}'
```

**响应示例：**

```json
{
  "success": true,
  "code": 0,
  "message": "更新成功",
  "data": { "id": 1, "name": "张三丰", "age": 100 },
  "timestamp": 1714800000
}
```

#### 删除用户

```
DELETE /users/:id
```

```bash
curl -X DELETE http://localhost:8080/users/1
```

**响应示例：**

```json
{
  "success": true,
  "code": 0,
  "message": "删除成功",
  "data": null,
  "timestamp": 1714800000
}
```

---

### 文件下载

#### 获取文件列表

```
GET /files
```

```bash
curl http://localhost:8080/files
```

**响应示例：**

```json
{
  "success": true,
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      { "name": "test.txt", "size": 1024, "url": "/download/test.txt" }
    ],
    "total": 1
  },
  "timestamp": 1714800000
}
```

#### 下载文件

```
GET /download/:filename
```

```bash
# 在浏览器中直接访问
http://localhost:8080/download/test.txt

# 或使用 curl 下载
curl -O http://localhost:8080/download/test.txt
```

> 可下载的文件放置在项目根目录的 `static/` 文件夹中。

---

## 注意事项

- 本项目为演示用途，数据存储在内存中，服务重启后数据会丢失
- 登录接口使用硬编码账号密码，Token 为模拟值，未实现真实的 JWT 认证
- 用户模型中 `age` 字段校验范围为 0-150
- 文件下载接口已做路径穿越防护，不允许文件名中包含 `..`
