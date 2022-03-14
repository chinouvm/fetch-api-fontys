# fhict-api-fetch

Fetch een lijst van docenten met verschillende filters!

## Usage

Create a config.json file with this code:

```JSON
{
  "email": {
    "from": "Source Email",
    "to": "Destination Email",
    "smtpPassword": "Email Password",
    "mailserver": "Email Server",
    "mailport": "Email Port"
  },
  "api": {
    "address": "https://api.fhict.nl/people",
    "authToken": "Authentication Token from api.fhict.nl"
  }
}
```

## Example

```JSON
{
  "email": {
    "from": "mygmail@gmail.com",
    "to": "myotheremail@email.com",
    "smtpPassword": "password123",
    "mailserver": "smtp.gmail.com",
    "mailport": "587"
  },
  "api": {
    "address": "https://api.fhict.nl/people",
    "authToken": "TOKEN"
  }
}
```
