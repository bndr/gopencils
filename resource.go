// Copyright 2014 Vadim Kravcenko
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package gopencils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var queryString map[string]string

type Resource struct {
	Api         *ApiStruct
	Url         string
	id          string
	Querystring map[string]string
	Payload     io.Reader
	Headers     http.Header
	Response    interface{}
	Error       interface{}
	Raw         *http.Response
}

func (r *Resource) Res(options ...interface{}) *Resource {
	if len(options) > 0 {
		url := r.Url + "/" + options[0].(string)

		r.Api.Methods[url] = &Resource{Url: url, Api: r.Api, Headers: http.Header{}}

		if len(options) > 1 {
			r.Api.Methods[url].Response = options[1]
		}

		return r.Api.Methods[url]
	}
	return nil
}

func (r *Resource) Id(value interface{}) *Resource {
	r.id = value.(string)
	r.Url += "/" + r.id
	return r
}

func (r *Resource) SetQuery(querystring map[string]string) *Resource {
	r.Querystring = querystring
	return r
}

func (r *Resource) Get(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Querystring = options[0].(map[string]string)
	}
	return r.do("GET")
}

func (r *Resource) Head(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Querystring = options[0].(map[string]string)
	}
	return r.do("HEAD")
}

func (r *Resource) Put(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("PUT")
}

func (r *Resource) Post(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("POST")
}

func (r *Resource) Delete(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Querystring = options[0].(map[string]string)
	}
	return r.do("DELETE")
}

func (r *Resource) Options(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Querystring = options[0].(map[string]string)
	}
	return r.do("OPTIONS")
}

func (r *Resource) Patch(options ...interface{}) *Resource {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("PATCH")
}

func (r *Resource) do(method string) *Resource {
	url := r.parseUrl()
	req, err := http.NewRequest(method, url, r.Payload)
	if err != nil {
		panic(err)
	}

	if r.Api.BasicAuth != nil {
		req.SetBasicAuth(r.Api.BasicAuth.Username, r.Api.BasicAuth.Username)
	}

	if r.Headers != nil {
		req.Header = r.Headers
	}

	resp, err := r.Api.Client.Do(req)
	if err != nil {
		panic(err)
	}

	r.Raw = resp

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, r.Response)

	if err != nil {
		panic(err)
	}

	return r
}

func (r *Resource) SetPayload(args interface{}) io.Reader {
	var b []byte
	b, _ = json.Marshal(args)
	r.SetHeader("Content-Type", "application/json")
	return bytes.NewBuffer(b)
}

func (r *Resource) SetHeader(key string, value string) {
	r.Headers.Add(key, value)
}

func (r *Resource) SetClient(c *http.Client) {
	r.Api.Client = c
}

func (r Resource) parseUrl() string {
	url := r.Api.Base + r.Url
	separator := "?"
	for k, v := range r.Querystring {
		url += fmt.Sprintf("%s%s=%s", separator, k, v)
		separator = "&"
	}
	return url
}
