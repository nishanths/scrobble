# Scrobble API documentation

Welcome to the documentation for v1 of Scrobble's HTTP API. The API primarily uses
JSON. The API provides ways to fetch account details, fetch a user's scrobbled
songs, and scrobble new songs, among other things.

## Base URL

The base URL for the API is https://selective-scrobble.appspot.com.

## Authentication

Most endpoints listed below require authentication (see each endpoint for details).
Authentication is done by including the following HTTP header in requests.

```
X-Scrobble-API-Key: <your API key>
```

Your API key can be found in the [dashboard](https://scrobbl.es/dashboard/api-key) when signed in.

## Types

The request and response bodies for some endpoints refer to named types such as
`Song`. These types are described in the [Types](/types) page.

## Endpoints

* [`/api/v1/account`](/account)
* [`/api/v1/account/delete`](/account__delete)
* [`/api/v1/scrobbled`](/scrobbled)
* [`/api/v1/scrobbled/color`](/scrobbled__color)
* [`/api/v1/scrobble`](/scrobble)
* [`/api/v1/artwork`](/artwork)
* [`/api/v1/artwork/missing`](/artwork__missing)
