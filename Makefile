build:
	docker-compose build

go-test-image:
	docker build . -f test.dockerfile -t example-adserver/go-test 

test-unit: go-test-image
	docker run --rm -it -v "$(PWD):/go/src/app" -w "/go/src/app" example-adserver/go-test ./... 

test-e2e: build go-test-image
	docker-compose up -d
	# wait for everything to be ready
	sleep 20
	docker run \
		--env-file="$(PWD)/e2e/test.env" \
		--rm -it \
		-v "$(PWD):/go/src/app" \
		-w "/go/src/app" \
		--network="example-adserver_default" \
		example-adserver/go-test -tags e2e ./e2e
	docker-compose stop

test-cleanup:
	docker-compose kill
	docker-compose rm -f 