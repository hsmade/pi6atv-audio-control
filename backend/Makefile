build: vendor test vet
	 GOARCH=arm GOARM=7 GOOS=linux go build

vendor:
	go mod vendor

test:
	GOARCH=arm GOARM=7 GOOS=linux go test ./...

vet:
	GOARCH=arm GOARM=7 GOOS=linux go vet ./...