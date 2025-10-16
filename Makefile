include .mk/test.mk

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

TLE_FILE_DIR = tests/tle
KUBENET_FILE_DIR = tests/kubenet
OUTPUT_TEST_DIR = tests/output
TEMPLATE_DIR = template

.PHONY: tle-file-dir
tle-file-dir:
	mkdir -p $(TLE_FILE_DIR)

$(OUTPUT_TEST_DIR):
	mkdir -p $(OUTPUT_TEST_DIR)

.PHONY: s1-tle s2-tle s3-tle s1 s2 s3 ground1 ground2 ground3 two-shell multi-shell

# include on demand for your tests
include .mk/multi_shell_mixed_bgp_ospf_test.mk
# include .mk/multi_shell_no_routing_test.mk
# include .mk/multi_shell_ospf_test.mk
# include .mk/multi_shell_wxwx_test.mk
# include .mk/physic_test.mk
