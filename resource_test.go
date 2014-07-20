package gopencils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.Equal(t, api.parseUrl(), "https://test-url.com/", "Parsed Url Should match")
	api.SetQuery(map[string]string{"key": "value"})
	assert.Equal(t, api.parseUrl(), "https://test-url.com/?key=value", "Parsed QueryString Url Should match")
	assert.Equal(t, api.Url, "", "Base Url Should be empty")
	assert.Equal(t, api.Api.Base, "https://test-url.com/", "Base Url Should be empty")
}

func TestResource_auth(t *testing.T) {
	api := Api("https://test-url.com", &BasicAuth{"username", "password"})
	assert.Equal(t, api.Api.BasicAuth.Username, "username", "Username should match")
	assert.Equal(t, api.Api.BasicAuth.Password, "password", "Password should match")
}

func TestResource_get(t *testing.T) {
	api := Api("https://api.github.com")
	// Users endpoint
	users := api.Res("users")

	usernames := []string{"bndr", "torvalds"}

	for _, username := range usernames {
		// Create a new pointer to response Struct
		r := new(respStruct)
		// Get user with id i into the newly created response struct
		_, err := users.Id(username, r).Get()
		if err != nil {
			t.Log("Error Getting Data from Github API")
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
	api := Api("http://httpbin.org")
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(binStruct)
	api.Res("post", r).Post(payload)
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestResource_update(t *testing.T) {
	api := Api("https://httpbin.org")
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(binStruct)
	api.Res("put", r).Put(payload)
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestResource_delete(t *testing.T) {
	api := Api("https://httpbin.org")
	r := new(binStruct)
	api.Id("delete", r).Delete()
	assert.Equal(t, r.Url, "https://httpbin.org/delete", "Url must match")
}

func TestResource_id(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Res("users").Id("test").Url, "users/test", "Url should match")
	assert.Equal(t, api.Res("users").Id("test").parseUrl(), "https://test-url.com/users/test", "Url should match")
	assert.Equal(t, api.Res("users").Id(123).Res("items").Id(111).Url, "users/123/items/111", "Multilevel Url should match")
}
