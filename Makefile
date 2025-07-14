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

TESTS ?= tests/two_shell_bgp_tests
TEMPLATE ?= template/bgp_template.yaml

# TESTS ?= tests/two_shell_ospf_tests
# TEMPLATE ?= template/ospf_template.yaml

# TESTS ?= tests/two_shell_ng_bgp_tests
# TEMPLATE ?= template/bgp_template.yaml

# TESTS ?= tests/two_shell_ng_ospf_tests
# TEMPLATE ?= template/ospf_template.yaml

# TESTS ?= tests/multi_shell_bgp_tests
# TEMPLATE ?= template/bgp_template.yaml

# TESTS ?= tests/multi_shell_ospf_tests
# TEMPLATE ?= template/ospf_template.yaml

.PHONY: s1-tle s2-tle s3-tle s1 s2 s3 ground1 ground2 ground3 two-shell multi-shell

s1-tle:
	go run ./main.go tle -N 60 -P 10 -E 10 -I 85 -H 1175 -F 0 -D 100 -o ${TESTS}/s1.tle -n S1

s2-tle:
	go run ./main.go tle -N 108 -P 9 -E 9 -I 55 -H 1150 -F 0 -D 100 -o ${TESTS}/s2.tle -n S2

s3-tle:
	go run ./main.go tle -N 648 -P 36 -E 36 -I 89 -H 1100 -F 0 -D 100 -o ${TESTS}/s3.tle -n S3

s1: s1-tle
	go run ./main.go generate -N 60 -P 10 -E 10 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	--offset 0 -F 0 -i ${TESTS}/s1.tle -o ${TESTS}/s1.yaml -t ${TEMPLATE} -r 1 -n S1

s2: s2-tle
	go run ./main.go generate -N 108 -P 9 -E 9 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 \
	--offset 0 -F 0 -i ${TESTS}/s2.tle -o ${TESTS}/s2.yaml -t ${TEMPLATE} -r 2 -n S2

s3: s3-tle
	go run ./main.go generate -N 648 -P 36 -E 36 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 \
	--offset 0 -F 0 -i ${TESTS}/s3.tle -o ${TESTS}/s3.yaml -t ${TEMPLATE} -r 3 -n S3

ground1:
	go run ./main.go ground -N 20 --longitude 113.2 --latitude 19.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	-o ${TESTS}/ground1.yaml -t ${TEMPLATE} -r 1 -n Ground1

ground2:
	go run ./main.go ground -N 40 --longitude 113.2 --latitude 15.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 \
	-o ${TESTS}/ground2.yaml -t ${TEMPLATE} -r 2 -n Ground2
	
ground3:
	go run ./main.go ground -N 140 --longitude 113.2 --latitude 11.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 \
	-o ${TESTS}/ground3.yaml -t ${TEMPLATE} -r 3 -n Ground3

two-shell: s1 s2 ground1 ground2
	go run ./main.go merge -i ./${TESTS} -o test.yaml

two-shell-ng: s1 s2
	go run ./main.go merge -i ./${TESTS} -o test.yaml

multi-shell: s1 s2 s3 ground1 ground2 ground3
	go run ./main.go merge -i ./${TESTS} -o test.yaml

multi-shell-ng: s1 s2 s3
	go run ./main.go merge -i ./${TESTS} -o test.yaml