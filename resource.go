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

// Gopencils is a Golang REST Client with which you can easily consume REST API's. Supported Response formats: JSON
package gopencils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var ErrCantUseAsQuery = errors.New("can't use options[0] as Query")

// Resource is basically an url relative to given API Baseurl.
type Resource struct {
	Api         *ApiStruct
	Url         string
	id          string
	QueryValues url.Values
	Payload     io.Reader
	Headers     http.Header
	Response    interface{}
	Raw         *http.Response
}

// Creates a new Resource.
func (r *Resource) Res(options ...interface{}) *Resource {
	if len(options) > 0 {
		var url string
		if len(r.Url) > 0 {
			url = r.Url + "/" + options[0].(string)
		} else {
			url = options[0].(string)
		}

		newR := &Resource{Url: url, Api: r.Api, Headers: http.Header{}}

		if len(options) > 1 {
			newR.Response = options[1]
		}

		return newR
	}
	return r
}

// Same as Res() Method, but returns a Resource with url resource/:id
func (r *Resource) Id(options ...interface{}) *Resource {
	if len(options) > 0 {
		id := ""
		switch v := options[0].(type) {
		default:
			id = v.(string)
		case int:
			id = strconv.Itoa(v)
		}
		url := r.Url + "/" + id
		newR := &Resource{id: id, Url: url, Api: r.Api, Headers: http.Header{}, Response: &r.Response}

		if len(options) > 1 {
			newR.Response = options[1]
		}
		return newR
	}
	return r
}

// Sets QueryValues for current Resource
func (r *Resource) SetQuery(querystring map[string]string) *Resource {
	for k, v := range querystring {
		r.QueryValues.Set(k, v)
	}
	return r
}

// Performs a GET request on given Resource
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Get(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		if qry, ok := options[0].(map[string]string); ok {
			r.SetQuery(qry)
		} else {
			return nil, ErrCantUseAsQuery
		}

	}
	return r.do("GET")
}

// Performs a HEAD request on given Resource
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Head(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		if qry, ok := options[0].(map[string]string); ok {
			r.SetQuery(qry)
		} else {
			return nil, ErrCantUseAsQuery
		}
	}
	return r.do("HEAD")
}

// Performs a PUT request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (r *Resource) Put(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("PUT")
}

// Performs a POST request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (r *Resource) Post(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("POST")
}

// Performs a Delete request on given Resource.
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Delete(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		if qry, ok := options[0].(map[string]string); ok {
			r.SetQuery(qry)
		} else {
			return nil, ErrCantUseAsQuery
		}
	}
	return r.do("DELETE")
}

// Performs a Delete request on given Resource.
// Accepts map[string]string as parameter, will be used as querystring.
func (r *Resource) Options(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		if qry, ok := options[0].(map[string]string); ok {
			r.SetQuery(qry)
		} else {
			return nil, ErrCantUseAsQuery
		}
	}
	return r.do("OPTIONS")
}

// Performs a PATCH request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (r *Resource) Patch(options ...interface{}) (*Resource, error) {
	if len(options) > 0 {
		r.Payload = r.SetPayload(options[0])
	}
	return r.do("PATCH")
}

// Main method, opens the connection, sets basic auth, applies headers,
// parses response json.
func (r *Resource) do(method string) (*Resource, error) {
	r.Api.BaseUrl.Path = r.Url
	r.Api.BaseUrl.RawQuery = r.QueryValues.Encode()
	url := r.Api.BaseUrl.String()
	req, err := http.NewRequest(method, url, r.Payload)
	if err != nil {
		return r, err
	}

	if r.Api.BasicAuth != nil {
		req.SetBasicAuth(r.Api.BasicAuth.Username, r.Api.BasicAuth.Password)
	}

	if r.Headers != nil {
		for k, _ := range r.Headers {
			req.Header.Set(k, r.Headers.Get(k))
		}
	}

	resp, err := r.Api.Client.Do(req)
	if err != nil {
		return r, err
	}

	r.Raw = resp

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(r.Response)
	if err != nil {
		return r, err
	}

	return r, nil
}

// Sets Payload for current Resource
func (r *Resource) SetPayload(args interface{}) io.Reader {
	var b []byte
	b, _ = json.Marshal(args)
	r.SetHeader("Content-Type", "application/json")
	return bytes.NewBuffer(b)
}

// Sets Headers
func (r *Resource) SetHeader(key string, value string) {
	r.Headers.Add(key, value)
}

// Overwrites the client that will be used for requests.
// For example if you want to use your own client with OAuth2
func (r *Resource) SetClient(c *http.Client) {
	r.Api.Client = c
}
