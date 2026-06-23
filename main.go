package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const version = "1.1.0"

// ---- config ----

type Config struct {
	BaseURL   string `json:"baseUrl"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginPath string `json:"loginPath"`
}

func loadConfig() *Config {
	paths := []string{
		"yapi-mcp.config.json",
		filepath.Join(userHomeDir(), ".yapi-mcp", "config.json"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			var cfg Config
			if json.Unmarshal(data, &cfg) == nil {
				return &cfg
			}
		}
	}
	return nil
}

func userHomeDir() string {
	if v := os.Getenv("HOME"); v != "" {
		return v
	}
	return os.Getenv("USERPROFILE")
}

// ---- tool helpers ----

func resultJSON(v any) *mcp.CallToolResult {
	text, _ := json.MarshalIndent(v, "", "  ")
	return mcp.NewToolResultText(string(text))
}

func errorResult(msg string) *mcp.CallToolResult {
	return mcp.NewToolResultError(msg)
}

// ---- main ----

func main() {
	fileCfg := loadConfig()

	baseURL := os.Getenv("YAPI_BASE_URL")
	email := os.Getenv("YAPI_EMAIL")
	password := os.Getenv("YAPI_PASSWORD")
	loginPath := os.Getenv("YAPI_LOGIN_PATH")

	if baseURL == "" && fileCfg != nil {
		baseURL = fileCfg.BaseURL
	}
	if email == "" && fileCfg != nil {
		email = fileCfg.Email
	}
	if password == "" && fileCfg != nil {
		password = fileCfg.Password
	}
	if loginPath == "" && fileCfg != nil {
		loginPath = fileCfg.LoginPath
	}

	if baseURL == "" {
		fmt.Fprintln(os.Stderr, "Error: YAPI_BASE_URL is required")
		os.Exit(1)
	}
	if email == "" || password == "" {
		fmt.Fprintln(os.Stderr, "Error: YAPI_EMAIL and YAPI_PASSWORD are required")
		os.Exit(1)
	}

	baseURL = strings.TrimRight(baseURL, "/")

	client := NewYapiClient(YapiConfig{
		BaseURL:   baseURL,
		Email:     email,
		Password:  password,
		LoginPath: loginPath,
	})

	buildInfo := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		buildInfo = info.Main.Version
	}

	s := server.NewMCPServer(
		"yapi-mcp",
		version,
		server.WithInstructions("YAPI MCP Server - 支持 Cookie 认证，查询项目、接口列表、接口分类、接口详情，新增和修改接口"),
	)

	// 0. get version
	s.AddTool(mcp.NewTool("yapi-get-version",
		mcp.WithDescription("获取 YAPI MCP Server 版本信息"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return resultJSON(map[string]any{
			"name":    "yapi-mcp",
			"version": version,
			"build":   buildInfo,
		}), nil
	})

	// 1. get project
	s.AddTool(mcp.NewTool("yapi-get-project",
		mcp.WithDescription("获取 YAPI 项目信息"),
		mcp.WithNumber("projectId",
			mcp.Required(),
			mcp.Description("项目ID"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID := int(req.GetInt("projectId", 0))
		if projectID <= 0 {
			return errorResult("projectId is required"), nil
		}
		project, err := client.GetProject(projectID)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(project), nil
	})

	// 2. get project list
	s.AddTool(mcp.NewTool("yapi-get-project-list",
		mcp.WithDescription("获取 YAPI 项目列表"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projects, err := client.GetProjectList()
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(projects), nil
	})

	// 3. get categories
	s.AddTool(mcp.NewTool("yapi-get-categories",
		mcp.WithDescription("获取项目的接口分类列表"),
		mcp.WithNumber("projectId",
			mcp.Required(),
			mcp.Description("项目ID"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID := int(req.GetInt("projectId", 0))
		if projectID <= 0 {
			return errorResult("projectId is required"), nil
		}
		categories, err := client.GetCategories(projectID)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(categories), nil
	})

	// 4. get interfaces by category
	s.AddTool(mcp.NewTool("yapi-get-interfaces-by-category",
		mcp.WithDescription("获取指定分类下的接口列表"),
		mcp.WithNumber("catid",
			mcp.Required(),
			mcp.Description("分类ID"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		catID := int(req.GetInt("catid", 0))
		if catID <= 0 {
			return errorResult("catid is required"), nil
		}
		category, err := client.GetInterfacesByCategory(catID)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(category), nil
	})

	// 5. get all interfaces by project
	s.AddTool(mcp.NewTool("yapi-get-interfaces",
		mcp.WithDescription("获取项目的所有接口列表"),
		mcp.WithNumber("projectId",
			mcp.Required(),
			mcp.Description("项目ID"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID := int(req.GetInt("projectId", 0))
		if projectID <= 0 {
			return errorResult("projectId is required"), nil
		}
		ifaces, err := client.GetInterfacesByProject(projectID)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(ifaces), nil
	})

	// 6. get interface detail
	s.AddTool(mcp.NewTool("yapi-get-interface",
		mcp.WithDescription("获取接口详情，支持格式：数字ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("接口ID"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		idStr, err := req.RequireString("id")
		if err != nil {
			return errorResult("id is required"), nil
		}
		var id int
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			return errorResult("id must be a number"), nil
		}
		iface, err := client.GetInterface(id)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(iface), nil
	})

	// 7. create interface
	s.AddTool(mcp.NewTool("yapi-create-interface",
		mcp.WithDescription("新增 YAPI 接口"),
		mcp.WithString("title", mcp.Required(), mcp.Description("接口标题")),
		mcp.WithString("path", mcp.Required(), mcp.Description("接口路径")),
		mcp.WithString("method", mcp.Required(), mcp.Description("请求方法")),
		mcp.WithNumber("catid", mcp.Required(), mcp.Description("分类ID")),
		mcp.WithNumber("project_id", mcp.Required(), mcp.Description("项目ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := client.CreateInterface(args)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(result), nil
	})

	// 8. update interface
	s.AddTool(mcp.NewTool("yapi-update-interface",
		mcp.WithDescription("修改 YAPI 接口"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("接口ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := client.UpdateInterface(args)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(result), nil
	})

	// 9. create category
	s.AddTool(mcp.NewTool("yapi-create-category",
		mcp.WithDescription("新增 YAPI 接口分类"),
		mcp.WithString("name", mcp.Required(), mcp.Description("分类名称")),
		mcp.WithNumber("project_id", mcp.Required(), mcp.Description("项目ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := client.CreateCategory(args)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return resultJSON(result), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
