.PHONY: webui build run test vet docker clean

# Build the React web console into internal/webui/dist (embedded by Go).
webui:
	cd webui && npm install && npm run build

# Build the server binary (requires the web console to be built first).
build:
	go build -o kivo ./cmd/kivo

# Convenience: full local build.
full: webui build

run: full
	./kivo

test:
	go test ./...

vet:
	go vet ./...

docker:
	docker build -t kivo .

clean:
	rm -f kivo kivo-bench kivo-test
	rm -rf internal/webui/dist webui/node_modules