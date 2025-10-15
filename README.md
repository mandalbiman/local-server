# local-server


## Local setup
1. Install go in your system
Please visit this link https://go.dev/doc/install

2. Run `go mod tidy` inside this repo folder to download all the golang dependencies.

3. Run `go run cmd/main.go refresh` to load the user database

4. Run the server `go run cmd/main.go`

5.After this go to postman and in the address bar write `http://localhost:9090/api/v1` you will get `hello world` as response.

