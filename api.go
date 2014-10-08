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
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Basic Auth
type BasicAuth struct {
	Username string
	Password string
}

// Main Api Instance.
// No Options yet supported.
type ApiStruct struct {
	BaseUrl    *url.URL
	BasicAuth  *BasicAuth
	Client     *http.Client
	Cookies    *cookiejar.Jar
	PathSuffix string
}

// Create a new API Instance and returns a Resource
// Accepts URL as parameter, and either a Basic Auth or a OAuth2 Client.
func Api(baseUrl string, options ...interface{}) *Resource {
	u, err := url.Parse(baseUrl)
	if err != nil {
		// TODO: don't panic..
		panic("Api() - url.Parse(baseUrl) Error:" + err.Error())
	}

	apiInstance := &ApiStruct{BaseUrl: u, BasicAuth: nil}

	for _, o := range options {
		switch v := o.(type) {
		case *BasicAuth:
			apiInstance.BasicAuth = v
		case *http.Client:
			apiInstance.Client = v
		case string:
			apiInstance.PathSuffix = v
		}
	}

	if apiInstance.Client == nil {
		apiInstance.Cookies, _ = cookiejar.New(nil)

		// Skip verify by default?
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{
			Transport: tr,
			Jar:       apiInstance.Cookies,
		}
		apiInstance.Client = client
	}
	return &Resource{Url: "", Api: apiInstance}
}
