package gopencils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testMux *http.ServeMux
	testSrv *httptest.Server
)

func init() {
	testMux = http.NewServeMux()
	testSrv = httptest.NewServer(testMux)
}

type respStruct struct {
	Login   string
	Id      int
	Name    string
	Message string
}
type binStruct struct {
	Args    string
	Headers map[string]string
	Url     string
	Json    map[string]interface{}
}

func TestResource_url(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Api.BaseUrl.String(), "https://test-url.com", "Parsed Url Should match")
	api.SetQuery(map[string]string{"key1": "value1", "key2": "value2"})
	assert.Equal(t, api.QueryValues.Encode(), "key1=value1&key2=value2", "Parsed QueryString Url Should match")
	assert.Equal(t, api.Url, "", "Base Url Should be empty")
}

func TestResource_auth(t *testing.T) {
	api := Api("https://test-url.com", &BasicAuth{"username", "password"})
	assert.Equal(t, api.Api.BasicAuth.Username, "username", "Username should match")
	assert.Equal(t, api.Api.BasicAuth.Password, "password", "Password should match")
}

func TestResource_get(t *testing.T) {
	// github stubs
	testMux.HandleFunc("/users/bndr", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, `{"login": "bndr","id": 1145456,"avatar_url": "https://avatars.githubusercontent.com/u/1145456?","gravatar_id": "8d05db0b0b8b74d5a0f93d0b1db22909","url": "https://api.github.com/users/bndr","html_url": "https://github.com/bndr","followers_url": "https://api.github.com/users/bndr/followers","following_url": "https://api.github.com/users/bndr/following{/other_user}","gists_url": "https://api.github.com/users/bndr/gists{/gist_id}","starred_url": "https://api.github.com/users/bndr/starred{/owner}{/repo}","subscriptions_url": "https://api.github.com/users/bndr/subscriptions","organizations_url": "https://api.github.com/users/bndr/orgs","repos_url": "https://api.github.com/users/bndr/repos","events_url": "https://api.github.com/users/bndr/events{/privacy}","received_events_url": "https://api.github.com/users/bndr/received_events","type": "User","site_admin": false,"name": "Vadim Kravcenko","company": "","blog": "http://vadimkravcenko.com","location": "Germany","email": "bndrzz@gmail.com","hireable": false,"bio": null,"public_repos": 17,"public_gists": 2,"followers": 13,"following": 0,"created_at": "2011-10-22T19:21:17Z","updated_at": "2014-07-20T10:24:25Z"}`)
	})
	testMux.HandleFunc("/users/torvalds", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, `{"login": "torvalds","id": 1024025,"avatar_url": "https://avatars.githubusercontent.com/u/1024025?","gravatar_id": "fb47627bc8c0bcdb36321edfbf02e916","url": "https://api.github.com/users/torvalds","html_url": "https://github.com/torvalds","followers_url": "https://api.github.com/users/torvalds/followers","following_url": "https://api.github.com/users/torvalds/following{/other_user}","gists_url": "https://api.github.com/users/torvalds/gists{/gist_id}","starred_url": "https://api.github.com/users/torvalds/starred{/owner}{/repo}","subscriptions_url": "https://api.github.com/users/torvalds/subscriptions","organizations_url": "https://api.github.com/users/torvalds/orgs","repos_url": "https://api.github.com/users/torvalds/repos","events_url": "https://api.github.com/users/torvalds/events{/privacy}","received_events_url": "https://api.github.com/users/torvalds/received_events","type": "User","site_admin": false,"name": "Linus Torvalds","company": "Linux Foundation","blog": null,"location": "Portland, OR","email": null,"hireable": false,"bio": null,"public_repos": 2,"public_gists": 0,"followers": 17605,"following": 0,"created_at": "2011-09-03T15:26:22Z","updated_at": "2014-07-20T10:27:30Z"}`)
	})

	api := Api(testSrv.URL)
	// Users endpoint
	users := api.Res("users")

	usernames := []string{"bndr", "torvalds"}

	for _, username := range usernames {
		// Create a new pointer to response Struct
		r := new(respStruct)
		// Get user with id i into the newly created response struct
		_, err := users.Id(username, r).Get()
		if err != nil {
			t.Log("Error Getting Data from Test API\nErr:", err)
		} else {
			assert.Equal(t, r.Message, "", "Error message must be empty")
			assert.Equal(t, r.Login, username, "Username should be equal")
		}
	}
	resp := new(respStruct)
	api.Res("users", resp).Id("bndr").Get()
	assert.Equal(t, resp.Login, "bndr")
	resp2 := new(respStruct)
	api.Res("users").Id("bndr", resp2).Get()
	assert.Equal(t, resp2.Login, "bndr")
}

func TestResource_create(t *testing.T) {
	testMux.HandleFunc("/post", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "POST", "unexpected Method")
		assert.Equal(t, req.URL.Path, "/post", "unexpected Path")
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json", "Expected json content type")
		fmt.Fprintln(rw, `{"args": {},"data": "{\"Key\":\"Value1\"}","files": {},"form": {},"headers": {"Accept": "*/*",  "Connection": "close",  "Content-Length": "16",  "Content-Type": "application/json",  "Host": "httpbin.org",  "User-Agent": "curl/7.37.0",  "X-Request-Id": "6268bee8-2ea0-4144-802a-6166fe18d84f"},"json": {"Key": "Value1"},"origin": "95.91.230.168","url": "https://httpbin.org/post"}`)
	})

	api := Api(testSrv.URL)
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(binStruct)
	api.Res("post", r).Post(payload)
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestResource_update(t *testing.T) {
	testMux.HandleFunc("/put", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "PUT", "unexpected Method")
		assert.Equal(t, req.URL.Path, "/put", "unexpected Path")
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json", "Expected json content type")
		fmt.Fprintln(rw, `{"args": {},"data": "{\"Key\":\"Value1\"}","files": {},"form": {},"headers": {"Accept": "*/*",  "Connection": "close",  "Content-Length": "16",  "Content-Type": "application/json",  "Host": "httpbin.org",  "User-Agent": "curl/7.37.0",  "X-Request-Id": "6268bee8-2ea0-4144-802a-6166fe18d84f"},"json": {"Key": "Value1"},"origin": "95.91.230.168","url": "https://httpbin.org/post"}`)
	})

	api := Api(testSrv.URL)
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(binStruct)
	api.Res("put", r).Put(payload)
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestResource_delete(t *testing.T) {
	testMux.HandleFunc("/delete", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "DELETE", "unexpected Method")
		assert.Equal(t, req.URL.Path, "/delete", "unexpected Path")
		fmt.Fprintln(rw, `{"args": {},"data": "","files": {},"form": {},"headers": {"Accept": "*/*",  "Connection": "close",  "Host": "httpbin.org",  "User-Agent": "curl/7.37.0",  "X-Request-Id": "b29e2435-926f-4fb4-bd1d-ec1b179e1523"},"json": null,"origin": "95.91.230.168","url": "https://httpbin.org/delete"}`)
	})

	api := Api(testSrv.URL)
	r := new(binStruct)
	api.Id("delete", r).Delete()
	assert.Equal(t, r.Url, "https://httpbin.org/delete", "Url must match")
}

func TestResource_id(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Res("users").Id("test").Url, "users/test", "Url should match")
	assert.Equal(t, api.Res("users").Id(123).Res("items").Id(111).Url, "users/123/items/111", "Multilevel Url should match")
}
