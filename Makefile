MODULE := github.com/lhopki01/kubernetes-job-ui

test:
	go test -race ./...

test-cover:
	go test -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

lint: lint-go lint-javascript

lint-go:
	golangci-lint run

lint-javascript:
	cd website && ./node_modules/.bin/eslint src/**/**.jsx

build-binaryfs:
	cd website && rm -rf build && npm run build
	go-bindata -pkg site -o internal/site/bindata.go -prefix "website/build" website/build/...

release: build-binaryfs
	git tag -a $$VERSION
	git push origin $$VERSION
	goreleaser --rm-dist
