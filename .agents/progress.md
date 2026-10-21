# 任务进度

- 目标：为 yapi-mcp 增加与 YApi 网页一致的全局搜索，并让项目列表在不传 groupId 时返回当前用户可见的全部项目。
- 已完成：接入 `/api/project/search?q=...` 全局搜索；新增 `yapi-global-search` MCP 工具；项目列表支持可选 groupId，不传时读取分组并汇总项目；README 已更新。
- 验证：`GOWORK=off go test ./...` 通过；使用独立可写缓存执行 `GOWORK=off go build ./...` 通过；`git diff --check` 通过。
- 上下文快照：本次任务已完成。若继续增强，需先确认是否要扩展 YApi 原生全局搜索之外的接口描述、请求参数和响应体全文过滤。
