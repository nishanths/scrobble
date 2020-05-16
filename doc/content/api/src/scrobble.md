[Back to index](/)

# /api/v1/scrobble

## HTTP method

POST

## Request header

Requires authentication header. See [Authentication](/#authentication).

## Request content-type

application/json

## Request body

```
[]MediaItem
```

See [Types](/types).

## Response status code

| Code | Description |
|------|-------------|
|200 | success |
|400 | bad request body |
|401 | missing authentication credentials |
|403 | bad authentication credentials |
|405 | HTTP method not allowed |
|500 | various internal server errors |

