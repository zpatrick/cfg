VERSION ?= $(shell git rev-parse --short HEAD)

deps:
	go install -v golang.org/x/tools/cmd/godoc@latest

test:
	go test -v -parallel=2 -timeout=10s -shuffle=on ./...

docs:
	godoc -http=:6060

publish:
	git tag ${VERSION}
	git push origin ${VERSION}
	GOPROXY=proxy.golang.org go list -m github.com/zpatrick/cfg@${VERSION}

.PHONY: deps, test, docs, publish