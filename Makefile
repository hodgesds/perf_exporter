DATE := $(shell date --iso-8601=seconds)

build/perf_exporter: test | vendor
	@go build -o ./build/perf_exporter \
		-ldflags "-X main.BuildDate=${DATE}" \
		./perf_exporter

go.mod:
	@GO111MODULE=on go mod tidy

go.sum: | go.mod
	@GO111MODULE=on go mod verify

vendor: | go.sum
	@GO111MODULE=on go mod vendor 

.PHONEY: install
install: | vendor
	@go install -v -ldflags "-X main.BuildDate=${DATE} ./perf_exporter

test: | vendor
	@go test -v -race -cover ./...

.PHONEY: clean
clean:
	rm -rf build vendor
