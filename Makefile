export MONGO_URL:=mongodb://localhost:27017
export SERVER_URL:=:3000
export MONGO_COLLECTION=BAIRESDEV

compose:
	docker-compose -f docker-compose.yml up

build:
	go build -o ./cmd/dist/goga ./cmd/api.go

run: build
	./cmd/dist/goga

test:
	go test ./pkg/... -count=1 -coverprofile ./cover.out && \
    go tool cover -func ./cover.out && \
    rm ./cover.out
