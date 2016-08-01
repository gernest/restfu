# restfu
Simple RESTFul app in Go

# prerequisites
It is tested against the latest version of go go1.7 but it should work just fine
with go1.2+.

Uses gorilla mux for routing so install it like this

```shell
go get github.com/gorilla/mux
```

# Installation

You can use go get 

```shell
go get github.com/gernest/restfu
```

Then start the server like this
```shell
restfu -port 8080
```

If no flag for port is specified then port 8080 is used by default.

# CAVEATS

-  The products are stored in memory. Restatring the server loses all the
  previous records.
-  Records are stored in a map , although it is safe for concurrent use it uses
  locks to ensure safety , hence it has all the costs of locking.

- Reocrds returned by ` GET /products.json` are not ordered( In any sence ,
  neither id nor time created as we dont track the time creation)
- For updates , the product object recived from the request replaces the
  available product as is, meaning no missing fields are accounted for.

# Examples

## POST /products

Creates a new product record.

``` shell
curl curl -v -H "Content-Type:application/json" -X POST http://localhost:8080/products  -d '{"name":"smart", "desc": "smart templates"}
```

## POST/products/{id}

Updates a product specified by id

## GET /products.json

Returns a list of all products

## GET /products/{id},json

Returns a single product by id



