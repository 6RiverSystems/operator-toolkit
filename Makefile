
.PHONY: deepcopy

deepcopy:
	which ./bin/deepcopy-gen > /dev/null || go build -o ./bin/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen
	./bin/deepcopy-gen -i ./apis  --bounding-dirs "./apis/" -v=4 --output-base ./ --output-file-base zz_generated.deepcopy
