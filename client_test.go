package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGlobalSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/user/login":
			w.Header().Set("Set-Cookie", "_yapi_token=test-token; Path=/")
			_, _ = fmt.Fprint(w, `{"errcode":0,"errmsg":"成功","data":{"uid":7}}`)
		case "/api/project/search":
			if got := r.URL.Query().Get("q"); got != "订单 支付" {
				t.Errorf("search keyword = %q, want %q", got, "订单 支付")
			}
			_, _ = fmt.Fprint(w, `{"errcode":0,"errmsg":"ok","data":{"project":[],"group":[],"interface":[{"_id":12,"title":"支付订单"}]}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewYapiClient(YapiConfig{
		BaseURL:  server.URL,
		Email:    "user@example.com",
		Password: "password",
	})
	result, err := client.GlobalSearch(" 订单 支付 ")
	if err != nil {
		t.Fatalf("GlobalSearch() error = %v", err)
	}
	if len(result.Interface) != 1 || result.Interface[0]["title"] != "支付订单" {
		t.Fatalf("unexpected search result: %#v", result.Interface)
	}
}

func TestGlobalSearchRequiresKeyword(t *testing.T) {
	client := NewYapiClient(YapiConfig{})
	if _, err := client.GlobalSearch("   "); err == nil {
		t.Fatal("GlobalSearch() with empty keyword should fail")
	}
}

func TestGetProjectListWithoutGroupID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/user/login":
			w.Header().Set("Set-Cookie", "_yapi_token=test-token; Path=/")
			_, _ = fmt.Fprint(w, `{"errcode":0,"data":{"uid":7}}`)
		case "/api/group/list":
			_, _ = fmt.Fprint(w, `{"errcode":0,"data":[{"_id":3},{"_id":4}]}`)
		case "/api/project/list":
			_, _ = fmt.Fprintf(w, `{"errcode":0,"data":{"list":[{"_id":%s,"name":"项目%s"}]}}`, r.URL.Query().Get("group_id"), r.URL.Query().Get("group_id"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewYapiClient(YapiConfig{BaseURL: server.URL, Email: "user", Password: "password"})
	projects, err := client.GetProjectList()
	if err != nil {
		t.Fatalf("GetProjectList() error = %v", err)
	}
	if len(projects) != 2 || projects[0].ID != 3 || projects[1].ID != 4 {
		t.Fatalf("unexpected projects: %#v", projects)
	}
}
