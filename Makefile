include .mk/test.mk
include .mk/gen.mk
include make.env

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


