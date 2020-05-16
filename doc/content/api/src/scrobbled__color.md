[Back to index](/)

# /api/v1/scrobbled/color

Fetch a user's scrobbled albums by artwork color.

## HTTP method

GET

## Request header

The authentication header is required if the requested profile is private;
otherwise it is optional. For private profiles only the API keys of the profile
owner will work.

See [Authentication](/#authentication).

## Request query parameters

| Key | Value | Notes |
|-----|-------|-------|
| username | string | required |
| color | "red" \| "orange" \| "brown" \| "yellow" \| "green" \| "blue" \| "violet" \| "pink" \| "black" \| "gray" \| "white" | required |
| limit | number | optional |

## Response status code

| Code | Description |
|------|-------------|
|200 | success |
|400 | missing or invalid 'color' query parameter, or missing 'username' query parameter |
|403 | insufficient authentication credentials to view private account, or bad credentials |
|404 | no user exists for specified username |
|405 | HTTP method not allowed |
|500 | various internal server errors |

## Response content-type

application/json

## Response body

```
SongResponse[]
```
The response is a list of `SongResponse` objects. Each song in the response
corresponds any one of the songs in albums with matching artwork.

See [Types](/types).
