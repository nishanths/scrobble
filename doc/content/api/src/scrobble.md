[Back to index](/)

# /api/v1/scrobble

Scrobble your listening history. The list of scrobbled items in the request
is used to _entirely replace_ your old listening history.

## HTTP method

POST

## Request header

Requires authentication header. See [Authentication](/#authentication).

## Request content-type

application/json

## Request body

```
MediaItem[]
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

