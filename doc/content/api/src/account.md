[Back to index](/)

# /api/v1/account

Fetch information about your account.

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

```
{
  username: string
  private: boolean
}
```
