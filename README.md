---
date: 2020-03-15
tag:
  - golang
  - restful
  - mongo
author: Raiven Kao
location: Taipei
---

# BTC/USE price RESTful server based on golang

## demand analysis

### required

- based on golnag
- The API interface you provide can be any of the following：RESTful、json rpc、gRPC
  - choose restful
- At least two sources
- When a source is unavailable, the result of its last successful ask is returned
- Use git to manage source code
- Write readme.md and describe what features, features, and TODO which have been implemented

### optional

- Traffic limits, including the number of times your server queries the source, and the number of times the user queries you
- Good testing, annotations, and git commit
- An additional websocket interface is provided to automatically send the latest information whenever the market changes
- Users can choose to use an automatic source determination, the latest data source, or manually specify the source of the data
- Package it as a Dockerfile, docker-compose file, or kubernetes yaml
- There is a simple front-end or cli program that displays the results
- The API you provide has a corresponding file, such as a swagger, or simply a markdown file
- Other features not listed but that you thought would be cool to implement

### uaecase diagram

![](uml/usecase/usecase.png)

### sequence diagram

#### register

![](uml/sequence/register.png)

#### login

![](uml/sequence/login.png)

#### user get latest price

![](uml/sequence/get_latest_price.png)

#### server get remote price

![](uml/sequence/get_remote_price.png)

### sequence diagram with caching mechanism

If time permits, an caching mechanism will be added redis based

#### login with redis

![](uml/sequence/login_redis.png)

#### user get latest price with redis

![](uml/sequence/get_latest_price_redis.png)

## intro

### http router

Use [`gorilla/mux`](https://github.com/gorilla/mux) as http router.

A struct `route` for registering

```go
type route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []route
```

make a fuction `init` to initialize routes.

```go
func init() {
  register("GET", "/", handler, nil)
  ...
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, route{method, pattern, handler, middleware})
}
```

now when you create a new `mux`, should use this handler like this.

```go
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		if route.Middleware == nil {
			r.HandleFunc(route.Pattern, route.Handler).Methods(route.Method)
		} else {
			r.Handle(route.Pattern, route.Middleware(route.Handler)).Methods(route.Method)
		}
	}
	return r
}
```

Also can write a test file and use [`appleboy/gofight`](https://github.com/appleboy/gofight) to make sure your router work normally.

```go
func TestGetDefault(t *testing.T) {
	route := gofight.New()
	route.GET("/").
		Run(routes.NewRouter(), func(res gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "btc-price-restful", res.Body.String())
			assert.Equal(t, http.StatusOK, res.Code)
		})
}
```

### database

All of result data which will be responsed by romote API server have simple structure. Just key-value JSON. So I choose NoSQL database official [`mongodb`](https://github.com/mongodb/mongo-go-driver) package to make storage.

You can use context to ensure the gorutine context will end at the same time. If the mongodb can't connect, the program will fatal error and exit.

```go
ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
// ensure cancle function will be executed
defer cancle() 
client, err := mongo.Connect(ctx, clientOptions())
if err != nil {
	l.Fatal(err) // l for logrus
}
db.SetClint(client)
```

the clientOptions is for docker, can set parameter for different env in dockerfile or docker-compose file.

```go
func clientOptions() *options.ClientOptions {
	host := "db"
	if os.Getenv("profile") != "prod" {
		host = "localhost"
	}
	return options.Client().ApplyURI(
		"mongodb://" + host + ":27017",
	)
}
```