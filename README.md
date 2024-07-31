# ELL

| under development

this is a tool management for static token to comunication service to service

## Motivation
common usage static token
```
[ Service A ] -> req with static token -> [ Service B ]
[ Service C ] -> req with static token -> [ Service B ]
```

if service b is change the auth token, service a and service c must be change manually and posible to downtime.
maybe change manually is not a problem if you have 2 or 3 service, but if you have 10 or more service, it's a problem.

**With Ell**

You just need change the token in UI and all service will get the new token. And ell have feature for renew token.


## How to use

### installation
run server
```bash
$ git clone git@github.com/born2ngopi/ell.git
$ cd ell
$ go run main.go
```

for client go
```bash
$ go get github.com/born2ngopi/ell
```


### Simple Middleware
```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        token := ell.GetToken("service-b")
        auth := r.Header.Get("Authorization")
        if auth != token {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
    })
}
```

### simple client request
```go
func GetUser(ctx context.Context) (User, error) {

    token := ell.GetToken("service-b")
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://service-b/user", nil)
    if err != nil {
        return User{}, err
    }

    req.Header.Set("Authorization", token)
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return User{}, err
    }
    defer res.Body.Close()

    var user User
    if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
        return User{}, err
    }

    return user, nil
}
```
