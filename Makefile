
.PHONY: deepcopy

deepcopy:
	which ./bin/deepcopy-gen > /dev/null || go build -o ./bin/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen
	./bin/deepcopy-gen -i ./apis  --bounding-dirs "./apis/" -v=4 --output-base ./ --output-file-base zz_generated.deepcopy

format:
	find . -name \*.go -exec goimports -local github.com/6RiverSystems -w {} \;

lint:
	golangci-lint run

bench:
	go test -bench=./... -race ./...

test:
	go test -covermode=count -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
