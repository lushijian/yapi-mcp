# YAPI MCP Server

一个用于 YAPI 接口管理平台的 MCP (Model Context Protocol) 服务器，支持查询项目、接口列表、接口分类、接口详情，以及新增和修改接口等功能。

## 功能特性

- 🔍 **项目查询** - 获取项目信息和项目列表
- 📂 **分类管理** - 查询接口分类、新增分类
- 📋 **接口管理** - 查询接口列表、接口详情、新增接口、修改接口
- 🔐 **多种认证方式** - 支持 Cookie 认证（邮箱密码登录）和 Token 认证

## 安装

```bash
# 克隆项目
git clone https://github.com/W-Abel-jia/yapi-interface-mcp.git
cd yapi-interface-mcp

# 安装依赖
npm install

# 构建
npm run build
```

## 配置

### 方式一：配置文件（推荐）

在项目根目录创建 `yapi-interface-mcp.config.json` 文件：

```json
{
  "baseUrl": "http://your-yapi-server:3000",
  "email": "your-email@example.com",
  "password": "your-password"
}
```

或者使用 Token 认证：

```json
{
  "baseUrl": "http://your-yapi-server:3000",
  "token": "your-yapi-token"
}
```

### 方式二：环境变量

| 环境变量 | 说明 | 必填 |
|---------|------|------|
| `YAPI_BASE_URL` | YAPI 服务器地址 | ✅ |
| `YAPI_EMAIL` | 登录邮箱 | 二选一 |
| `YAPI_PASSWORD` | 登录密码 | 二选一 |
| `YAPI_TOKEN` | YAPI Token | 二选一 |

> 💡 认证方式优先级：Token > Cookie（邮箱密码登录）

## 使用方式

### 在 Claude Desktop 中使用

编辑 Claude Desktop 配置文件：
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

添加以下配置：

```json
{
  "mcpServers": {
    "yapi": {
      "command": "node",
      "args": ["/path/to/yapi-interface-mcp/dist/index.js"],
      "env": {
        "YAPI_BASE_URL": "http://your-yapi-server:3000",
        "YAPI_EMAIL": "your-email@example.com",
        "YAPI_PASSWORD": "your-password"
      }
    }
  }
}
```

重启 Claude Desktop 即可使用。

### 在其他 MCP 客户端中使用

如 Cursor、Windsurf、Cline 等，参考各客户端的 MCP 配置方式：

- **命令**: `node /path/to/yapi-interface-mcp/dist/index.js`
- **环境变量**: 配置 `YAPI_BASE_URL` 和认证信息

### 使用 MCP Inspector 测试

```bash
cd yapi-interface-mcp
npx @modelcontextprotocol/inspector node dist/index.js
```

## 支持的工具

| 工具名称 | 功能描述 | 必需参数 |
|---------|---------|---------|
| `yapi-get-project` | 获取项目信息 | `projectId` |
| `yapi-get-project-list` | 获取项目列表 | - |
| `yapi-get-categories` | 获取项目的接口分类列表 | `projectId` |
| `yapi-get-interfaces-by-category` | 获取指定分类下的接口列表 | `catid` |
| `yapi-get-interfaces` | 获取项目的所有接口列表 | `projectId` |
| `yapi-get-interface` | 获取接口详情 | `id` |
| `yapi-create-interface` | 新增接口 | `title`, `path`, `method`, `catid` |
| `yapi-update-interface` | 修改接口 | `id` |
| `yapi-create-category` | 新增分类 | `name`, `project_id` |

### 工具参数说明

#### yapi-get-project
```
projectId: number - 项目ID
```

#### yapi-get-categories
```
projectId: number - 项目ID
```

#### yapi-get-interfaces-by-category
```
catid: number - 分类ID
```

#### yapi-get-interfaces
```
projectId: number - 项目ID
```

#### yapi-get-interface
```
id: string - 接口ID（数字格式）
```

#### yapi-create-interface
```
title: string - 接口标题
path: string - 接口路径，如 /api/user
method: "GET" | "POST" | "PUT" | "DELETE" | "HEAD" | "OPTIONS" | "PATCH"
catid: number - 分类ID
desc?: string - 接口描述
status?: "done" | "designing" | "dev" | "undone"
req_body_type?: "form" | "json" | "file" | "raw" | "text"
req_body_other?: string - 请求体 JSON Schema
req_body_form?: Array<{name, type, required, desc}>
req_headers?: Array<{name, value?, required, desc?}>
req_params?: Array<{name, desc?}>
req_query?: Array<{name, required, desc?}>
res_body_type?: "json" | "raw" | "text"
res_body?: string - 响应体 JSON Schema
```

#### yapi-update-interface
```
id: number - 接口ID
...其他参数同 yapi-create-interface（均为可选）
```

#### yapi-create-category
```
name: string - 分类名称
project_id: number - 项目ID
desc?: string - 分类描述
```

## 开发

```bash
# 安装依赖
npm install

# 开发模式（监听文件变化）
npm run dev

# 构建
npm run build

# 运行
npm start
```

## 注意事项

1. **安全性** - 配置文件 `yapi-interface-mcp.config.json` 包含敏感信息，请勿提交到版本控制系统。已添加到 `.gitignore`。
2. **认证方式** - 推荐使用邮箱密码登录（Cookie 认证），YAPI 的 Token 认证可能有权限限制。
3. **网络访问** - 确保 MCP 服务器能够访问 YAPI 服务器。

## License

MIT
