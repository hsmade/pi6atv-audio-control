build: vendor test vet
	 GOARCH=arm GOARM=7 GOOS=linux go build cmd/backend.go

vendor:
	go mod vendor

test:
	go test ./...

vet:
	go vet ./...