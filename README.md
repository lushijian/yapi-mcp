# yapi-mcp

YAPI 接口管理平台的 MCP (Model Context Protocol) 服务器，Go 实现。支持标准登录和 LDAP 登录。

## 功能特性

- 项目查询、项目列表
- 接口分类查询、新增分类
- 接口列表、接口详情、新增接口、修改接口
- Cookie 认证，支持标准登录和 LDAP 登录

## 快速开始

### 方式一：go run（类似 npx，无需安装）

```bash
go run github.com/lushijian/yapi-mcp@latest
```

### 方式二：go install

```bash
go install github.com/lushijian/yapi-mcp@latest
yapi-mcp
```

## 配置

支持环境变量和配置文件 `yapi-mcp.config.json`：

| 环境变量 | 说明 | 必填 |
|---------|------|------|
| `YAPI_BASE_URL` | YAPI 服务器地址 | ✅ |
| `YAPI_EMAIL` | 登录邮箱 | ✅ |
| `YAPI_PASSWORD` | 登录密码 | ✅ |
| `YAPI_LOGIN_PATH` | 登录路径，默认 `/api/user/login` | ❌ |

### LDAP 登录示例

```json
{
  "baseUrl": "http://yapi.maizuo.com",
  "email": "batty",
  "password": "your-password",
  "loginPath": "/api/user/login_by_ldap"
}
```

### 标准登录示例

```json
{
  "baseUrl": "https://yapi.sdyxmall.com",
  "email": "user@example.com",
  "password": "your-password"
}
```

## MCP 客户端配置

```json
{
  "mcpServers": {
    "yapi": {
      "command": "go",
      "args": ["run", "github.com/lushijian/yapi-mcp@latest"],
      "env": {
        "YAPI_BASE_URL": "http://your-yapi-server:3000",
        "YAPI_EMAIL": "your-email@example.com",
        "YAPI_PASSWORD": "your-password"
      }
    }
  }
}
```

### LDAP 登录配置

```json
{
  "mcpServers": {
    "yapi": {
      "command": "go",
      "args": ["run", "github.com/lushijian/yapi-mcp@latest"],
      "env": {
        "YAPI_BASE_URL": "http://yapi.maizuo.com",
        "YAPI_EMAIL": "batty",
        "YAPI_PASSWORD": "your-password",
        "YAPI_LOGIN_PATH": "/api/user/login_by_ldap"
      }
    }
  }
}
```

## 支持的工具

| 工具名称 | 功能描述 | 必需参数 |
|---------|---------|---------|
| `yapi-get-version` | 获取版本信息 | - |
| `yapi-get-project` | 获取项目信息 | `projectId` |
| `yapi-get-project-list` | 获取项目列表 | - |
| `yapi-get-categories` | 获取接口分类列表 | `projectId` |
| `yapi-get-interfaces-by-category` | 获取分类下接口列表 | `catid` |
| `yapi-get-interfaces` | 获取项目所有接口 | `projectId` |
| `yapi-get-interface` | 获取接口详情 | `id` |
| `yapi-create-interface` | 新增接口 | `title`, `path`, `method`, `catid`, `project_id` |
| `yapi-update-interface` | 修改接口 | `id` |
| `yapi-create-category` | 新增分类 | `name`, `project_id` |

## License

MIT
