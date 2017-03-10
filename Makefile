all:
	go install -ldflags -v  ./

install: all
	@echo

test:
	cd src; CGO_ENABLED=0 go test ./...

clean:
	cd src; go clean -i ./...

style:
	@$(QCHECKSTYLE) src

gofmt:
	find . -name '*.go' | xargs -l1 go fmt
