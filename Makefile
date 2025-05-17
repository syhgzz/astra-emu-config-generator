NN ?= kubenet/example
DATA ?= ./test.yaml
SERVER ?= http://localhost:9090
# SERVER ?= http://10.0.0.31:31002

TARGETS = \
    windows/amd64 \
    linux/amd64 \
    linux/arm64

MAIN = main.go

OUTPUT_DIR = build

# 默认目标
all: $(TARGETS)

# 通用构建规则
$(TARGETS):
	@GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/$(@:/=-)/kubenet-config-generator$(if $(findstring windows,$@),.exe,) $(MAIN)
	@echo "Built $@"

clean:
	rm -rf $(OUTPUT_DIR)

.PHONY: exp-apply exp-delete exp-start exp-stop make-config

exp-apply:
	curl -X POST -H "Content-Type: application/x-yaml" \
		--data-binary @${DATA} \
		"$(SERVER)/control/network_experiment/apply?nn=${NN}"

exp-delete:
	curl -X POST "$(SERVER)/control/network_experiment/delete?nn=${NN}"

exp-start:
	curl -X POST "$(SERVER)/control/scene/start?nn=${NN}"

exp-stop:
	curl -X POST "$(SERVER)/control/scene/stop?nn=${NN}"

make-config:
	go run ./main.go tle -N 6000 -P 100 -I 86 -H 500 -F 0 -E 1 -D 100 -o parts/s1.tle -n S1
	go run ./main.go generate -N 6000 -P 100 -E 1 --accessIntfNum 1 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 --offset 0 -F 0 -i parts/s1.tle -o parts/s1.yaml -n S1
	go run ./main.go tle -N 5400 -P 100 -I 70 -H 600 -F 0.3 -E 2 -D 100 -o parts/s2.tle -n S2
	go run ./main.go generate -N 5400 -P 100 -E 2 --accessIntfNum 1 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 --offset -1 -F 0.3 -i parts/s2.tle -o parts/s2.yaml -n S2
	go run ./main.go tle -N 14400 -P 200 -I 88 -H 1000 -F 0 -E 9 -D 100 -o parts/s3.tle -n S3
	go run ./main.go generate -N 14400 -P 200 -E 9 --accessIntfNum 1 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 --offset 0 -F 0 -i parts/s3.tle -o parts/s3.yaml -n S3
	go run ./main.go merge -i ./parts -o test.yaml