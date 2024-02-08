test:
	@echo Running tests...
	go test

memory_test:
	@echo Running memory tests...
	go test --tags memory_test

fmt_check:
	@echo Checking format...
	@test -z $(shell go fmt ./...) || (echo "Code is not formatted according to gofmt. Please run 'go fmt ./...' to fix." && exit 1)