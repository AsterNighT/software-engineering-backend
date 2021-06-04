# Create

Request

```json
// Usage
{
    "id": 1234567,  // uint
    "email": "^\\w+@\\w+[.\\w+]+$",  // String
    "type": "admin",  // string, in {"patient", "doctor", "admin"}
    "name": "Scott",
    "passwd": "unhashed passwd"  // use account.HashPassword(), discuss later
}
// demo
localhost:12448/api/account/create
{
    "id": 123456,
    "email": "a@a.com",
    "type": "admin",
    "name": "Scott",
    "passwd": "12345678"
}
```

Return:

```json
{
    "status": "Created",
    "data": {
        "account": {
            "ID": 123456,
            "Email": "a@a.com",
            "Type": "admin",
            "Name": "Scott",
            "Passwd": "$2a$10$YE383brJLvif0Y5Q3QGusOyTfR51eIjB63BioDWsVLyy.Hq4aoV/G"
        },
        "cookie_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2fQ.6k2qzGh3ub4L3ttRKXL4IqGHyOLAc7Wg5i0xto9XgQk"
    }
}
```

# Log in

Request:

```json
localhost:12448/api/account/login
{
    "email": "a@a.com",
    "passwd": "12345678"
}
```

Return:

```json
{
    "status": "Logged in",
    "data": {
        "account": {
            "ID": 123456,
            "Email": "a@a.com",
            "Type": "admin",
            "Name": "Scott",
            "Passwd": "$2a$10$YE383brJLvif0Y5Q3QGusOyTfR51eIjB63BioDWsVLyy.Hq4aoV/G"
        },
        "cookie_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2fQ.6k2qzGh3ub4L3ttRKXL4IqGHyOLAc7Wg5i0xto9XgQk"
    }
}
```

# Log out

Operation:

```json
localhost:12448/api/account/123456/logout
```

Return:

```json
{
    "status": "Account logged out"
}
```

# Modify password

Request:

```json
localhost:12448/api/account/123456/modifypasswd
{
    "email": "a@a.com",
    "passwd": "12345678",
    "newpasswd": "123456789"
}
```

Return:

```json
{
    "status": "Successfully modified"
}
```

