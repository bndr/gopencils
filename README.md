# Gopencils - Dynamically consume REST APIs
[![GoDoc](https://godoc.org/github.com/bndr/gopencils?status.svg)](https://godoc.org/github.com/bndr/gopencils)
[![Build Status](https://travis-ci.org/bndr/gopencils.svg?branch=master)](https://travis-ci.org/bndr/gopencils)

## Summary
Gopencils is a REST Client written in go. Easily consume any REST API's. Supported Response formats: JSON

## Install

    go get github.com/bndr/gopencils

## Simple to use

Gopencils was designed to help you easily make requests to REST APIs without much hassle. It supports both Basic-Auth as well as OAuth.

Example Basic-Auth

```go

type UserExample struct {
	Id            string
	Name          string
	Origin        string
	Url           string
	SomeJsonField string
}
// Create Basic Auth
auth := gopencils.BasicAuth{"username", "password"}

// Create New Api with our auth
api := gopencils.Api("http://your-api-url.com/api/", &auth)

// Create a pointer to our response struct
resp := &UserExample{}

// Perform a GET request
// URL Requested: http://your-api-url.com/api/users/1
api.Res("users", resp).Id(1).Get()

// Get Single Item
api.Res("users", resp).Id(1).Get()

// Perform a GET request with Querystring
querystring := map[string]string{"page": "100", "per_page": "1000"}

// URL Requested: http://your-api-url.com/api/users/123/items?page=100&per_page=1000
resource := api.Res("users").Id(123).Res("items", resp).Get(querystring)

// Now resp contains the returned json object
// resource.Raw contains raw http response,

// You can supply Path suffix to the api which will be applied to every url
// e.g /items/id.json
api := gopencils.Api("http://your-api-url.com/api/", &auth, ".json")

```

Example Github Api

```go

package main

import (
	"fmt"
	"github.com/bndr/gopencils"
)

type respStruct struct {
	Login string
	Id    int
	Name  string
}

func main() {
	api := gopencils.Api("https://api.github.com")
	// Users Resource
	users := api.Res("users")

	usernames := []string{"bndr", "torvalds", "coleifer"}

	for _, username := range usernames {
		// Create a new pointer to response Struct
		r := new(respStruct)
		// Get user with id i into the newly created response struct
		_, err := users.Id(username, r).Get()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r)
		}
	}
}
```
More examples in the examples folder.

## Why?
I work a lot with REST APIs and I caught myself writing the same code over and over, so I decided to make a library that would help me (and others) to quickly consume them.

## Is it ready?

It is still in beta. But I would be glad if you could test it on your pet projects. The API will be improved, but no breaking changes are planned.

## Contribute

All Contributions are welcome. The todo list is on the bottom of this README. Feel free to send a pull request.

## License

Apache License 2.0

## TODO
1. Add more Options (Flexibility)
2. Support XML Response
3. Better Error Handling
