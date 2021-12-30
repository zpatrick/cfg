test:
	go test -shuffle=on -v ./...

doc:
	go run golang.org/x/tools/cmd/godoc -http=:6060

.PHONY: test, doc