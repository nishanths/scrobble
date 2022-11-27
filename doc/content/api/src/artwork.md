[Back to index](/)

# /api/v1/artwork

**This endpoint is not functional and will respond with a 410 HTTP status.
Previous documentation is below.**

Upload any missing artwork, previously determined using [`/api/v1/artwork/missing`](/artwork__missing).

## HTTP method

POST

## Request header

Requires authentication header. See [Authentication](/#authentication).

## Request query parameters

| Key | Value | Notes |
|-----|-------|-------|
| format | "GIF" \| "PNG" \| "JPEG" \| "JPEG2000" | the format of the image data |

## Request body

The artwork image data as bytes.

## Response status code

| Code | Description |
|------|-------------|
|200 | success |
|400 | failed to decode artwork image data |
|401 | missing authentication credentials |
|403 | bad authentication credentials |
|405 | HTTP method not allowed |
|500 | various internal server errors |

## Response content-type

application/json

## Response body

The artwork hash as a JSON string.

```
ArtworkHash
```

See [Types](/types).
