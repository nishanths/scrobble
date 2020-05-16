[Back to index](/)

# /api/v1/scrobbled

Fetch a user's scrobbled albums by artwork color.

## HTTP method

GET

## Request header

The authentication header is required if the requested profile is private;
otherwise it is optional. For private profiles only the API keys of the profile
owner will work.

See [Authentication](/#authentication).

## Request query parameters

The query parameters should be in one of these two forms.

The first form can be used to fetch all scrobbled songs or just
loved scrobbled songs.

| Key | Value | Notes |
|-----|-------|-------|
| username | string | required |
| loved | "true" | optional; whether to return only loved songs; any other value but "true" corresponds to false |
| limit | number | optional |

The second form can be used to fetch a specific song.

| Key | Value | Notes |
|-----|-------|-------|
| username | string | required |
| song | SongIdentifier | required; see [Types](/types) |

## Response status code

| Code | Description |
|------|-------------|
|200 | success |
|400 | missing 'username' query parameter, or incorrect combination of query parameters |
|403 | insufficient authentication credentials to view private account, or bad credentials |
|404 | no user exists for specified username |
|405 | HTTP method not allowed |
|500 | various internal server errors |

## Response content-type

application/json

## Response body

For the first form of request query parameters, the response body structure is:

```
{
  // Total number of matching songs (not restricted by the "limit"
  // query parameter).
  total: number
  // Matching songs.
  songs: SongResponse[]
}
```

For the second form of request query parameters, the response body structure is:

```
{
  // The single song corresponding to the request, or an empty list if
  // there is no such song.
  songs: [SongResponse] | []
}
```

See [Types](/types).
