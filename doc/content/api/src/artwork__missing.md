[Back to index](/)

# /api/v1/artwork/missing

Fetch the list of missing artwork for your account.

## HTTP method

GET

## Request header

Requires authentication header. See [Authentication](/#authentication).

## Response status code

| Code | Description |
|------|-------------|
|200 | success |
|401 | missing authentication credentials |
|403 | bad authentication credentials |
|405 | HTTP method not allowed |
|500 | various internal server errors |

## Response content-type

application/json

## Response body

An object whose keys are the artwork hashes corresponding to missing artwork.
The values are always `true`.

```
{ [hash: ArtworkHash]: true }
```

See [Types](/types).
