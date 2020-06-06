.PHONY: unit-tests
go-test-image:
	docker build . -f test.dockerfile -t example-adserver/go-test 

unit-tests: go-test-image
	docker run --rm -it -v "$(PWD):/go/src/app" -w "/go/src/app" example-adserver/go-test go test ./...