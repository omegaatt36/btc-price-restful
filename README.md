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

- [ ] Traffic limits, including the number of times your server queries the source, and the number of times the user queries you

Good testing, annotations, and git commit

An additional websocket interface is provided to automatically send the latest information whenever the market changes

Users can choose to use an automatic source determination, the latest data source, or manually specify the source of the data

Package it as a Dockerfile, docker-compose file, or kubernetes yaml

There is a simple front-end or cli program that displays the results

The API you provide has a corresponding file, such as a swagger, or simply a markdown file

Other features not listed but that you thought would be cool to implement

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