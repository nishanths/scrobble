# /api/v1/artwork/color

Fetch a user's scrobbled albums by artwork color.

The response is a list of `SongResponse` objects. Each song in the response
corresponds to any one of the songs in matching albums.

### HTTP method

GET

### Request header

Requires authentication header. See [Authentication](/index.html).

### Request query parameters


	username: string
	color:    "red" | "orange" | "brown" | "yellow" | "green" | "blue" | "violet" | "pink" | "black" | "gray" | "white"
	limit:    number // optional


### Response status code

| Code | Description |
|------|-------------|
|200 | success |
|400 | missing or invalid 'color' query parameter, or missing 'username' query parameter |
|403 | insufficient credentials to view private account |
|404 | no user found for specified username |
|405 | HTTP method not allowed |
|500 | various internal server errors |

## Response content-type

application/json

## Response body

```
SongResponse[]
```

See [Types](/types.html).
