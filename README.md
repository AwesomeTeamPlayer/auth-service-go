#auth-service-go

## Available methods:

### App.CreateAccount
Params:
* email

Returns `true` on success, `false` otherwise.

Creates `users.created.auth` event:
```json
{
  "email_address": "address@email.com"
}
```

### App.Register
Params:
* email
* password

Returns `true` on success, `false` otherwise.

Creates `users.registered` event:
```json
{
  "email_address": "address@email.com"
}
```

### App.Login
Params:
* email
* password
* label

Returns session key string on success, empty string otherwise.

Creates `users.logged.in` event:
```json
{
  "email_address": "address@email.com"
}
```

### App.Logout
Params:
* email
* session key

Returns `true` if session existed before, `false` otherwise.

Creates `users.logged.out` event:
```json
{
  "email_address": "address@email.com"
}
```

### App.GetEmails
Params:
* page (start counting from 0)
* limit (at least 1)

Returns emails in ascent order:
```json
{
  "results": [ 
    {
      "emailAddress": "address@email.com",
      "hasPassword": true
    },
    //  ...
  ],
  "countAll": 123
}
```

### App.GetLoggedUsers
Params:
* page (start counting from 0)
* limit (at least 1)

Returns emails in ascent order:
```json
{
  "results": [
    "address@email.com",
    //  ...
  ],
  "countAll": 123
}
```

### App.GetSessions
Params:
* email

It returns users list of existing sessions:
```json
[
   {
      "email": "address@email.com",
      "key":"abc123...",
      "label": "Some text"
   },
   // ...
]
```
