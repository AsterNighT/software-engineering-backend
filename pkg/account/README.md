# Create

Request

```json
{
    "id": 1234567,  // uint
    "email": "^\\w+@\\w+[.\\w+]+$",  // String
    "type": "admin",  // string, in {"patient", "doctor", "admin"}
    "name": "Scott",
    "passwd": "unhashed passwd"  // use account.HashPassword(), discuss later
}
```

Return:

```go
c.JSON(http.StatusOK, api.Return("Created", echo.Map{
    "account":      account,
    "cookie_token": token,
}))
```

# Log in

Request:

```json
{
    "email": "^\\w+@\\w+[.\\w+]+$",  // String
    "passwd": "unhashed passwd"  // use account.HashPassword(), discuss later
}
```

Return:

```go
c.JSON(http.StatusOK, api.Return("Logged in", echo.Map{
    "account":      account,
    "cookie_token": token,
}))
```

# Log out

Operation:

```go
tokenCookie, _ := c.Get("tokenCookie").(*http.Cookie)

tokenCookie.Value = ""
tokenCookie.Expires = time.Unix(0, 0)
```

Return:

```go
c.JSON(http.StatusOK, api.Return("Account logged out", nil))
```

# Modify password

Request:

```json
{
    "email": "^\\w+@\\w+[.\\w+]+$",  // String
    "passwd": "unhashed passwd",
    "newpasswd": "unhashed passwd"
}
```

Return:

```go
c.JSON(http.StatusOK, api.Return("Successfully modified", nil))
```

