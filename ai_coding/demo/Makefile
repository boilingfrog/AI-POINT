.PHONY: lint vet test run build

# lint：用 gofmt 作为格式门禁（golangci-lint/golint 可能未安装）。
# 有未格式化文件就列出并以非零退出，作为真正的门禁而非提示。
lint:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "未格式化的文件:"; echo "$$out"; exit 1; fi

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

run:
	go run ./cmd/server
