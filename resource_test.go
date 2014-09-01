package gopencils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
type httpbinResponse struct {
	Args    string
	Headers map[string]string
	Url     string
	Json    map[string]interface{}
}

func TestResouceUrl(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Api.BaseUrl.String(), "https://test-url.com",
		"Parsed Url Should match")
	api.SetQuery(map[string]string{"key1": "value1", "key2": "value2"})
	assert.Equal(t, api.QueryValues.Encode(), "key1=value1&key2=value2",
		"Parsed QueryString Url Should match")
	assert.Equal(t, api.Url, "", "Base Url Should be empty")
}

func TestCanUsePathInResourceUrl(t *testing.T) {
	testMux.HandleFunc("/path/to/api/resname/id123",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(rw, `{"Test":"Okay"}`)
		})

	res := Api(testSrv.URL+"/path/to/api", nil)

	var resp struct {
		Test string
	}

	_, err := res.Res("resname").Id("id123", &resp).Get()
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "Okay", resp.Test, "resp shoul be Okay")
}

func TestCanUseAuthForApi(t *testing.T) {
	api := Api("https://test-url.com", &BasicAuth{"username", "password"})
	assert.Equal(t, api.Api.BasicAuth.Username, "username",
		"Username should match")
	assert.Equal(t, api.Api.BasicAuth.Password, "password",
		"Password should match")
}

func TestCanGetResource(t *testing.T) {
	// github stubs
	testMux.HandleFunc("/users/bndr",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(rw, readJson("_tests/github_bndr.json"))
		})
	testMux.HandleFunc("/users/torvalds",
		func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(rw, readJson("_tests/github_torvalds.json"))
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

func TestCanCreateResource(t *testing.T) {
	testMux.HandleFunc("/post",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "POST", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/post", "unexpected Path")
			assert.Equal(t, req.Header.Get("Content-Type"), "application/json",
				"Expected json content type")
			fmt.Fprintln(rw, readJson("_tests/common_response.json"))
		})

	api := Api(testSrv.URL)
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(httpbinResponse)
	api.Res("post", r).Post(payload)
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestCanPutResource(t *testing.T) {
	testMux.HandleFunc("/put",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "PUT", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/put", "unexpected Path")
			assert.Equal(t, req.Header.Get("Content-Type"), "application/json",
				"Expected json content type")
			fmt.Fprintln(rw, readJson("_tests/common_response.json"))
		})

	api := Api(testSrv.URL)
	payload := map[string]interface{}{"Key": "Value1"}
	r := new(httpbinResponse)
	api.Res("put", r).Put(payload)
	assert.Equal(t, r.Json["Key"], "Value1", "Payload must match")
}

func TestCanDeleteResource(t *testing.T) {
	testMux.HandleFunc("/delete",
		func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, "DELETE", "unexpected Method")
			assert.Equal(t, req.URL.Path, "/delete", "unexpected Path")
			fmt.Fprintln(rw, readJson("_tests/delete_response.json"))
		})

	api := Api(testSrv.URL)
	r := new(httpbinResponse)
	api.Id("delete", r).Delete()
	assert.Equal(t, r.Url, "https://httpbin.org/delete", "Url must match")
}

func TestResource_id(t *testing.T) {
	api := Api("https://test-url.com")
	assert.Equal(t, api.Res("users").Id("test").Url, "users/test",
		"Url should match")
	assert.Equal(t, api.Res("users").Id(123).Res("items").Id(111).Url,
		"users/123/items/111", "Multilevel Url should match")
}

func TestDoNotDecodeBodyOnErr(t *testing.T) {
	tests := []int{
		400, 401, 500, 501,
	}

	// will be changed later
	code := 0
	testMux.HandleFunc("/error",
		func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(code)
			fmt.Fprintln(rw, readJson("_tests/error.json"))
		})

	api := Api(testSrv.URL)

	for _, code = range tests {
		resp := make(map[string]interface{})
		api.Id("error", &resp).Get()

		assert.Equal(t, resp, map[string]interface{}{},
			fmt.Sprintf("response should be unparsed: %d", code))
	}
}

func readJson(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(buf)
}
