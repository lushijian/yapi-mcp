#!/usr/bin/env node

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { existsSync, readFileSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { YapiClient } from "./client.js";

interface Config {
  baseUrl: string;
  token?: string;
  email?: string;
  password?: string;
}

// 尝试从配置文件读取配置
function loadConfig(): Config | null {
  // 尝试多个可能的配置文件路径
  const possiblePaths = [
    // 当前工作目录
    join(process.cwd(), "yapi-interface-mcp.config.json"),
    // 脚本所在目录 (dist/index.js -> ../yapi-interface-mcp.config.json)
    join(dirname(fileURLToPath(import.meta.url)), "..", "yapi-interface-mcp.config.json"),
    // 用户主目录
    join(process.env.HOME || "", ".yapi-interface-mcp", "config.json"),
  ];

  for (const configPath of possiblePaths) {
    if (existsSync(configPath)) {
      try {
        const content = readFileSync(configPath, "utf-8");
        console.error(`[yapi-interface-mcp] Loaded config from: ${configPath}`);
        return JSON.parse(content);
      } catch (e) {
        console.error(`[yapi-interface-mcp] Warning: Failed to parse ${configPath}:`, e);
      }
    }
  }
  return null;
}

// 从环境变量或配置文件获取配置
const fileConfig = loadConfig();
const YAPI_BASE_URL = process.env.YAPI_BASE_URL || fileConfig?.baseUrl;
const YAPI_TOKEN = process.env.YAPI_TOKEN || fileConfig?.token;
const YAPI_EMAIL = process.env.YAPI_EMAIL || fileConfig?.email;
const YAPI_PASSWORD = process.env.YAPI_PASSWORD || fileConfig?.password;

if (!YAPI_BASE_URL) {
  console.error("Error: YAPI_BASE_URL is required. Set it via environment variable or yapi-mcp.config.json");
  process.exit(1);
}

if (!YAPI_TOKEN && !(YAPI_EMAIL && YAPI_PASSWORD)) {
  console.error("Error: Either YAPI_TOKEN or both YAPI_EMAIL and YAPI_PASSWORD are required");
  console.error("You can set them via environment variables or create a yapi-mcp.config.json file");
  process.exit(1);
}

// 创建 YAPI 客户端
const yapiClient = new YapiClient({
  baseUrl: YAPI_BASE_URL,
  token: YAPI_TOKEN,
  email: YAPI_EMAIL,
  password: YAPI_PASSWORD,
});

// 创建 MCP 服务器
const server = new McpServer({
  name: "yapi-mcp",
  version: "1.0.0",
});

// ==================== 工具定义 ====================

// 1. 获取项目信息
server.tool(
  "yapi-get-project",
  "获取 YAPI 项目信息",
  {
    projectId: z.number().describe("项目ID"),
  },
  async ({ projectId }) => {
    try {
      const project = await yapiClient.getProject(projectId);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(project, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 2. 获取项目列表
server.tool(
  "yapi-get-project-list",
  "获取 YAPI 项目列表",
  {},
  async () => {
    try {
      const projects = await yapiClient.getProjectList();
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(projects, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 3. 获取接口分类列表
server.tool(
  "yapi-get-categories",
  "获取项目的接口分类列表",
  {
    projectId: z.number().describe("项目ID"),
  },
  async ({ projectId }) => {
    try {
      const categories = await yapiClient.getCategories(projectId);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(categories, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 4. 获取分类下的接口列表
server.tool(
  "yapi-get-interfaces-by-category",
  "获取指定分类下的接口列表",
  {
    catid: z.number().describe("分类ID"),
  },
  async ({ catid }) => {
    try {
      const category = await yapiClient.getInterfacesByCategory(catid);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(category, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 5. 获取项目的所有接口列表
server.tool(
  "yapi-get-interfaces",
  "获取项目的所有接口列表",
  {
    projectId: z.number().describe("项目ID"),
  },
  async ({ projectId }) => {
    try {
      const interfaces = await yapiClient.getInterfacesByProject(projectId);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(interfaces, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 6. 获取接口详情
server.tool(
  "yapi-get-interface",
  "获取接口详情，支持格式：数字ID",
  {
    id: z.string().describe("接口ID，格式：数字ID"),
  },
  async ({ id }) => {
    try {
      const interfaceDetail = await yapiClient.getInterface(parseInt(id, 10));
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(interfaceDetail, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 7. 新增接口
server.tool(
  "yapi-create-interface",
  "新增 YAPI 接口",
  {
    title: z.string().describe("接口标题"),
    path: z.string().describe("接口路径，如 /api/user"),
    method: z.enum(["GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"]).describe("请求方法"),
    catid: z.number().describe("分类ID"),
    desc: z.string().optional().describe("接口描述"),
    status: z.enum(["done", "designing", "dev", "undone"]).optional().describe("接口状态"),
    req_body_type: z.enum(["form", "json", "file", "raw", "text"]).optional().describe("请求体类型"),
    req_body_other: z.string().optional().describe("请求体 JSON Schema"),
    req_body_form: z.array(
      z.object({
        name: z.string().describe("字段名"),
        type: z.enum(["text", "file"]).describe("字段类型"),
        required: z.enum(["1", "0"]).describe("是否必填"),
        desc: z.string().optional().describe("字段描述"),
      })
    ).optional().describe("表单字段列表"),
    req_headers: z.array(
      z.object({
        name: z.string().describe("请求头名称"),
        value: z.string().optional().describe("请求头值"),
        required: z.enum(["1", "0"]).describe("是否必填"),
        desc: z.string().optional().describe("描述"),
      })
    ).optional().describe("请求头列表"),
    req_params: z.array(
      z.object({
        name: z.string().describe("参数名"),
        desc: z.string().optional().describe("参数描述"),
      })
    ).optional().describe("路径参数列表"),
    req_query: z.array(
      z.object({
        name: z.string().describe("参数名"),
        required: z.enum(["1", "0"]).describe("是否必填"),
        desc: z.string().optional().describe("参数描述"),
      })
    ).optional().describe("Query 参数列表"),
    res_body_type: z.enum(["json", "raw", "text"]).optional().describe("响应体类型"),
    res_body: z.string().optional().describe("响应体 JSON Schema"),
  },
  async (params) => {
    try {
      const result = await yapiClient.createInterface(params);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 8. 修改接口
server.tool(
  "yapi-update-interface",
  "修改 YAPI 接口",
  {
    id: z.number().describe("接口ID"),
    title: z.string().optional().describe("接口标题"),
    path: z.string().optional().describe("接口路径"),
    method: z.enum(["GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"]).optional().describe("请求方法"),
    catid: z.number().optional().describe("分类ID"),
    desc: z.string().optional().describe("接口描述"),
    status: z.enum(["done", "designing", "dev", "undone"]).optional().describe("接口状态"),
    req_body_type: z.enum(["form", "json", "file", "raw", "text"]).optional().describe("请求体类型"),
    req_body_other: z.string().optional().describe("请求体 JSON Schema"),
    req_body_form: z.array(
      z.object({
        name: z.string().describe("字段名"),
        type: z.enum(["text", "file"]).describe("字段类型"),
        required: z.enum(["1", "0"]).describe("是否必填"),
        desc: z.string().optional().describe("字段描述"),
      })
    ).optional().describe("表单字段列表"),
    req_headers: z.array(
      z.object({
        name: z.string().describe("请求头名称"),
        value: z.string().optional().describe("请求头值"),
        required: z.enum(["1", "0"]).describe("是否必填"),
        desc: z.string().optional().describe("描述"),
      })
    ).optional().describe("请求头列表"),
    req_params: z.array(
      z.object({
        name: z.string().describe("参数名"),
        desc: z.string().optional().describe("参数描述"),
      })
    ).optional().describe("路径参数列表"),
    req_query: z.array(
      z.object({
        name: z.string().describe("参数名"),
        required: z.enum(["1", "0"]).describe("是否必填"),
        desc: z.string().optional().describe("参数描述"),
      })
    ).optional().describe("Query 参数列表"),
    res_body_type: z.enum(["json", "raw", "text"]).optional().describe("响应体类型"),
    res_body: z.string().optional().describe("响应体 JSON Schema"),
  },
  async (params) => {
    try {
      const result = await yapiClient.updateInterface(params);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 9. 新增分类
server.tool(
  "yapi-create-category",
  "新增 YAPI 接口分类",
  {
    name: z.string().describe("分类名称"),
    project_id: z.number().describe("项目ID"),
    desc: z.string().optional().describe("分类描述"),
  },
  async ({ name, project_id, desc }) => {
    try {
      const result = await yapiClient.createCategory({
        name,
        project_id,
        desc,
      });
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      return {
        content: [
          {
            type: "text" as const,
            text: JSON.stringify({ error: errorMessage }, null, 2),
          },
        ],
        isError: true,
      };
    }
  }
);

// 启动服务器
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("YAPI MCP Server running on stdio");
}

main().catch((error) => {
  console.error("Fatal error in main():", error);
  process.exit(1);
});
