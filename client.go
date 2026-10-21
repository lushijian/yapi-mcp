package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type YapiConfig struct {
	BaseURL   string `json:"baseUrl"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginPath string `json:"loginPath"`
}

type yapiLoginResponse struct {
	UID      int    `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type YapiProject struct {
	ID       int    `json:"_id"`
	Name     string `json:"name"`
	Basepath string `json:"basepath"`
	Desc     string `json:"desc"`
	Members  []struct {
		UID      int    `json:"uid"`
		Role     string `json:"role"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"members"`
}

type YapiCategory struct {
	ID        int             `json:"_id"`
	Name      string          `json:"name"`
	Desc      string          `json:"desc"`
	ProjectID int             `json:"project_id"`
	List      []YapiInterface `json:"list"`
}

type YapiGroup struct {
	ID int `json:"_id"`
}

type YapiInterface struct {
	ID                  int      `json:"_id"`
	Title               string   `json:"title"`
	Path                string   `json:"path"`
	Method              string   `json:"method"`
	ProjectID           int      `json:"project_id"`
	CatID               int      `json:"catid"`
	UID                 int      `json:"uid"`
	AddTime             int64    `json:"add_time"`
	UpTime              int64    `json:"up_time"`
	Desc                string   `json:"desc"`
	Status              string   `json:"status"`
	ReqBodyType         string   `json:"req_body_type"`
	ReqBodyOther        string   `json:"req_body_other"`
	ResBody             string   `json:"res_body"`
	ResBodyType         string   `json:"res_body_type"`
	ReqBodyIsJsonSchema bool     `json:"req_body_is_json_schema"`
	ResBodyIsJsonSchema bool     `json:"res_body_is_json_schema"`
	ReqHeaders          []any    `json:"req_headers"`
	ReqQuery            []any    `json:"req_query"`
	ReqParams           []any    `json:"req_params"`
	ReqBodyForm         []any    `json:"req_body_form"`
	Tag                 []string `json:"tag"`
	Type                string   `json:"type"`
	APIOpened           bool     `json:"api_opened"`
	EditUID             int      `json:"edit_uid"`
	Index               int      `json:"index"`
}

type YapiInterfaceDetail struct {
	YapiInterface
	Username  string `json:"username"`
	Editor    string `json:"editor,omitempty"`
	QueryPath any    `json:"query_path,omitempty"`
}

// YapiGlobalSearchResult 是 YApi 全局搜索返回的项目、分组和接口结果。
type YapiGlobalSearchResult struct {
	Project   []map[string]any `json:"project"`
	Group     []map[string]any `json:"group"`
	Interface []map[string]any `json:"interface"`
}

type yapiApiResponse[T any] struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    T      `json:"data"`
}

type YapiClient struct {
	config  YapiConfig
	client  *http.Client
	cookies []*http.Cookie
	uid     int
}

func NewYapiClient(config YapiConfig) *YapiClient {
	if config.LoginPath == "" {
		config.LoginPath = "/api/user/login"
	}
	return &YapiClient{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *YapiClient) login() error {
	body := map[string]string{
		"email":    c.config.Email,
		"password": c.config.Password,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.config.BaseURL+c.config.LoginPath, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// save set-cookie headers
	c.cookies = resp.Cookies()

	var apiResp yapiApiResponse[yapiLoginResponse]
	if err := json.Unmarshal(respBody, &apiResp); err == nil && apiResp.Errcode == 0 {
		c.uid = apiResp.Data.UID
		return nil
	}
	if apiResp.Errcode != 0 {
		return fmt.Errorf("YAPI login error: %s (code: %d)", apiResp.Errmsg, apiResp.Errcode)
	}

	// check if we got yapi cookies
	for _, cookie := range c.cookies {
		if strings.HasPrefix(cookie.Name, "_yapi_") {
			return nil
		}
	}
	return fmt.Errorf("YAPI login failed: no yapi cookies received")
}

func (c *YapiClient) ensureAuth() error {
	if len(c.cookies) > 0 {
		return nil
	}
	return c.login()
}

func (c *YapiClient) request(path string, method string, body any) ([]byte, error) {
	if err := c.ensureAuth(); err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.config.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// check for api error
	var apiErr struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Errcode != 0 {
		return nil, fmt.Errorf("YAPI error: %s (code: %d)", apiErr.Errmsg, apiErr.Errcode)
	}

	return respBody, nil
}

func (c *YapiClient) GetProject(projectID int) (*YapiProject, error) {
	body, err := c.request(fmt.Sprintf("/api/project/get?id=%d", projectID), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[YapiProject]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &resp.Data, nil
}

func (c *YapiClient) GetGroups() ([]YapiGroup, error) {
	body, err := c.request("/api/group/list", "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[[]YapiGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (c *YapiClient) getProjectListByGroup(groupID int) ([]YapiProject, error) {
	body, err := c.request(fmt.Sprintf("/api/project/list?group_id=%d", groupID), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[struct {
		List []YapiProject `json:"list"`
	}]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data.List, nil
}

// GetProjectList 获取指定分组或当前用户可见的全部项目。
func (c *YapiClient) GetProjectList(groupIDs ...int) ([]YapiProject, error) {
	if len(groupIDs) > 0 && groupIDs[0] > 0 {
		return c.getProjectListByGroup(groupIDs[0])
	}

	groups, err := c.GetGroups()
	if err != nil {
		return nil, err
	}
	projects := make([]YapiProject, 0)
	for _, group := range groups {
		if group.ID <= 0 {
			continue
		}
		groupProjects, err := c.getProjectListByGroup(group.ID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, groupProjects...)
	}
	return projects, nil
}

// GlobalSearch 调用 YApi 网页端使用的全局搜索接口。
func (c *YapiClient) GlobalSearch(keyword string) (*YapiGlobalSearchResult, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("keyword is required")
	}

	body, err := c.request("/api/project/search?q="+url.QueryEscape(keyword), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[YapiGlobalSearchResult]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &resp.Data, nil
}

func (c *YapiClient) GetCategories(projectID int) ([]YapiCategory, error) {
	body, err := c.request(fmt.Sprintf("/api/interface/getCatMenu?project_id=%d", projectID), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[[]YapiCategory]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (c *YapiClient) GetInterfacesByCategory(catID int) (*YapiCategory, error) {
	body, err := c.request(fmt.Sprintf("/api/interface/list_cat?catid=%d", catID), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[YapiCategory]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &resp.Data, nil
}

func (c *YapiClient) GetInterfacesByProject(projectID int) ([]YapiInterfaceDetail, error) {
	body, err := c.request(fmt.Sprintf("/api/interface/list?project_id=%d", projectID), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[[]YapiInterfaceDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (c *YapiClient) GetInterface(id int) (*YapiInterfaceDetail, error) {
	body, err := c.request(fmt.Sprintf("/api/interface/get?id=%d", id), "GET", nil)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[YapiInterfaceDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &resp.Data, nil
}

func (c *YapiClient) CreateInterface(params map[string]any) (*YapiInterfaceDetail, error) {
	params["uid"] = c.uid
	body, err := c.request("/api/interface/add", "POST", params)
	if err != nil {
		return nil, err
	}
	var resp yapiApiResponse[YapiInterfaceDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &resp.Data, nil
}

func (c *YapiClient) UpdateInterface(params map[string]any) (*YapiInterfaceDetail, error) {
	body, err := c.request("/api/interface/up", "POST", params)
	if err != nil {
		return nil, err
	}
	// /api/interface/up returns {"errcode":0,"data":{"n":1,...}} (MongoDB result)
	var resp yapiApiResponse[json.RawMessage]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &YapiInterfaceDetail{}, nil
}

func (c *YapiClient) CreateCategory(params map[string]any) (*YapiCategory, error) {
	body, err := c.request("/api/interface/add_cat", "POST", params)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Errcode int          `json:"errcode"`
		Errmsg  string       `json:"errmsg"`
		Data    YapiCategory `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &resp.Data, nil
}
