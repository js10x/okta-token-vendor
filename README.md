## Okta Token Vendor

A simple CLI for getting access tokens for testing Okta integrations. Simply provide the necessary arguments outlined below and the necessary requests will be made to the OKTA issuer to acquire an access token. This command line utility is only intended for development purposes for testing OKTA integrations.

### Compile From Source (default build with flag for removing the debug symbols from the binary)

```powershell
go build -ldflags=-w -o "oktv.exe"
```

### Okta Error Messages

In the case that the Okta issuer rejects your request configuration, an Okta specific error code is returned. The following is an example of this type of error.

```
Error occurred when fetching the SESSION TOKEN: 
Error Received From Okta:
Code: [E0000004]
Summary: [Authentication failed]
```

These codes are logged to the console during execution, you can find the full list of how to interpret these error codes on Okta's website, found here:

[Okta Error Codes](https://developer.okta.com/docs/reference/error-codes/)

### Flows Supported

* **Authorization Code Grant Flow with PKCE** (*Proof Key for Code Exchange*)

### Usage

The CLI will check your environment variables to find values for the following input variables, but you can pass them as flags when invoking the CLI as well. If you pass them in as flags to the CLI, those will take precedence.

* `CLIENT_ID`

* `ISSUER`

* `REDIRECT_URI` 

### All the flags

```powershell
oktv.exe -user "abc" -pw "abc" -cid "client_id" -iss "issuer" -callback "redirect uri" -o "path/to/file/token.txt"
```

#### Example usage (no output file provided):

```powershell
oktv.exe -user "abc" -pw "abc" -iss "https://okta-domain.com/oauth2/0x0" -cid "0x0" -callback "http://localhost:4200/login/callback"
```

#### Example usage (output file provided and using email instead of shortname):

```powershell
oktv.exe -user "abc" -pw "abc" -iss "https://okta-domain.com/oauth2/0x0" -cid "0x0" -callback "http://localhost:4200/login/callback" -o "path/to/file/token.txt"
```

### Help

The following arguments can be passed to the CLI to invoke the help documentation:

* -h

* --h

* -help

* --help