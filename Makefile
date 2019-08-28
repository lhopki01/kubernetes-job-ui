MODULE := github.com/lhopki01/kubernetes-job-ui

test:
	go test -race ./...

test-cover:
	go test -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

lint:
	golangci-lint run

release:
	git tag -a $$VERSION
	git push origin $$VERSION
	goreleaser --rm-dist

build:
	go-bindata -o bindata/bindata.go -pkg bindata website/build/...
	CGO_ENABLED=0 go build -ldflags "-X $(MODULE)/cmd.Version=$$VERSION"
