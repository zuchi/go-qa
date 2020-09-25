compose:
	docker-compose -f docker-compose.yml up

test:
	go test ./pkg/... -count=1 -coverprofile ./cover.out && \
    go tool cover -func ./cover.out && \
    rm ./cover.out
