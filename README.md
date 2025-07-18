# pos-dns-server

- This is a Go-based server that acts like a DNS router for company-specific database redirection.
- Each client provides a company name, and the server matches it to the registered IP/database endpoint.
- Initial binding between client and company is done using a passkey.

### Run the server
`make run`
or 
`go run cmd/server/main.go`


