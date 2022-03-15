VERSION ?= $(shell git rev-parse --short HEAD)

test:
	go test -v -parallel=2 -timeout=10s -shuffle=on ./...

docs:
	go run golang.org/x/tools/cmd/godoc -http=:6060

publish:
	git tag ${VERSION}
	git push origin ${VERSION}
	GOPROXY=proxy.golang.org go list -m github.com/zpatrick/cfg@${VERSION}

.PHONY: test, docs, publish