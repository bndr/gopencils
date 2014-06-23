# Gopencils - Dynamically consume REST APIs

## Summary
Easily consume any REST API's. Supported Response formats: JSON and XML(todo)

## Simple to use

Gopencils was designed to help you easily make requests to REST APIs without much hassle. It supports both Basic-Auth as well as OAuth.

Example Basic-Auth

```go
// Create Basic Auth
auth := gopencils.BasicAuth{"username", "password"}

// Create New Api with our auth
api := gopencils.Api("http://your-api-url.com/api/", &auth)
// Create a pointer to our response struct
resp := &respStruct{}

// Perform a GET request
// URL Requested: http://your-api-url.com/api/users/1
api.Res("users", resp).Id(1).Get()
```
More examples in the examples folder.

## Why?
I work a lot with REST APIs and I caught myself writing the same code over and over, so I decided to make a library that would help me (and others) to quickly consume them.

## Is it ready?

It is still in beta. But I would be glad if you could test it on your pet projects. The API will be improved, but no braking changes are planned. 

## Contribute

All Contributions are welcome. The todo list is on the bottom of this README. Feel free to send a pull request.

## License

Apache License 2.0

## TODO
0. Add more Options (Flexibility)
1. Support XML Response
2. Better Error Handling
