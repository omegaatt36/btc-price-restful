# API document

## Register

URL : `/user/register`

method : `POST`

Body : json

```josn
{
	"UserName":"username",
	"Password":"password"
}
```

### success

status code : 200

```json
{
  "msg": "",
  "res": "success",
  "data": null
}
```

### failed

#### case : user exist

status code : 400

```json
{
  "msg": "user exist",
  "res": "fail",
  "data": null
}
```

## Login

URL : `/user/register`

method : `POST`

Body : json

```josn
{
	"UserName":"username",
	"Password":"password"
}
```

### success

status code : 200

```json
{
  "msg": "",
  "res": "success",
  "data": {
    "token": `JWT`
  }
}
```

### failed

#### case : password error

status code : 500

```json
{
  "msg": "password error",
  "res": "fail",
  "data": null
}
```

### case : user not exist

status code : 400

```json
{
  "msg": "user not exist",
  "res": "fail",
  "data": null
}
```

## GetServiceMap

URL : `/getServiceMap`

method : `GET`

parameter : none

header :

```
{
    "Authorization" : `JWT`
}
```

### success

status code : 200

```json
{
    "msg": "",
    "res": "success",
    "data": [
        "CoinMarketCap",
        "CoinGecko",
        "CoinDesk"
    ]
}
```

### failed

#### case : no auth or wrong auth

status code : 401

```json
{
    "msg": "not authorized",
    "res": "fail",
    "data": null
}
```

## GetServiceMap

URL : `/getLatestPrice`

method : `GET`

parameter : source name separated by commas

e.g. `/getLatestPrice/CoinGecko,CoinMarketCap`

header : 

```
{
    "Authorization" : `JWT`
}
```

### success

status code : 200

```json
{
    "msg": "",
    "res": "success",
    "data": {
        "CoinGecko": {
            "source": "CoinGecko",
            "price": 5206.51,
            "timestamp": "2020-03-17T21:50:28+08:00"
        },
        "CoinMarketCap": {
            "source": "CoinMarketCap",
            "price": 5223.24143292,
            "timestamp": "2020-03-17T13:51:39.000Z"
        }
    }
}
```

### failed

#### case : no auth or wrong auth

status code : 401

```json
{
    "msg": "not authorized",
    "res": "fail",
    "data": null
}
```

#### case : no active service

status code : 500

```json
{
    "msg": "no active service",
    "res": "fail",
    "data": null
}
```

## GetServiceMap

URL : `/getLatestAllPrice`

method : `GET`

parameter : none

header : 

```
{
    "Authorization" : `JWT`
}
```

### success

status code : 200

```json
{
    "msg": "",
    "res": "success",
    "data": {
        "CoinDesk": {
            "source": "CoinDesk",
            "price": 5220.7567,
            "timestamp": "2020-03-17T13:52:00+00:00"
        },
        "CoinGecko": {
            "source": "CoinGecko",
            "price": 5208.02,
            "timestamp": "2020-03-17T21:46:07+08:00"
        },
        "CoinMarketCap": {
            "source": "CoinMarketCap",
            "price": 5223.24143292,
            "timestamp": "2020-03-17T13:51:39.000Z"
        }
    }
}
```

### failed

#### case : no auth or wrong auth

status code : 401

```json
{
    "msg": "not authorized",
    "res": "fail",
    "data": null
}
```

#### case : no active service

status code : 500

```json
{
    "msg": "no active service",
    "res": "fail",
    "data": null
}
```
