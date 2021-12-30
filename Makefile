test:
	go test -shuffle=on -v ./...

docs:
	go run golang.org/x/tools/cmd/godoc -http=:6060

.PHONY: test, docs