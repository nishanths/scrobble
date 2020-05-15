# API documentation

Welcome to the documentation for v1 of Scrobble's HTTP API.

## Base URL

The base URL for the API is https://selective-scrobble.appspot.com.

## Authentication

Most endpoints listed below require authentication. Authentication is done by
including the following HTTP header in requests.

```
X-Scrobble-API-Key: <your API key>
```

The API key can be found on the [dashboard](https://scrobble.littleroot.org/dashboard/api-key) when signed in.

## Types

The request and response structures for some endpoints refer to types such as
`Song`. These types are described in the [Types](/types) page.

## Endpoints

* [`/api/v1/account`](/account)
* [`/api/v1/account/delete`](/account__delete)
* [`/api/v1/scrobbled`](/scrobbled)
* [`/api/v1/artwork/color`](/artwork__color)
* [`/api/v1/scrobble`](/scrobble)
* [`/api/v1/artwork`](/artwork)
* [`/api/v1/artwork/missing`](/artwork__missing)
