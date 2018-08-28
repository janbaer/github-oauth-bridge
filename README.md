# Oauth bridge api

This Project contains my source code for my bridge to the Github Oauth. The reason why I wrote this
project is, that using the Github Oauth process is not possible from a web client.

The reason for it is, that the second call to Github in the Oauth workflow
(https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/#web-application-flow) 
is a Post request and it's not supporting CORS which makes it impossible to send this request from a
webclient.

The webapi needs a configuration file which contains all necessary configuration values that are
required to process the oauth login. Since they also have to contain the **secret** id, the config file
is not part of this Git repository.

But here's a sample config file:

```
[
  {
    "redirectUrl": "http://localhost:3000/#/login",
    "clientId": "12345",
    "clientSecretId": "23424fhsdfsd014141"
  }
]
```

To use the bridge you have to send a request to the following url {host}/login?clientId=12345

