include make.env

TLE_FILE_DIR = output/tle
KUBENET_CONFIG_FILE_DIR = output/kubenet
TEMPLATE_DIR = template
CONFIG_DIR = config

TLE_FILE_NAME = ${SHELL_NAME}.tle
KUBENET_CONFIG_FILE_NAME = ${SHELL_NAME}.yaml

CONFIG_PATH= ${CONFIG_DIR}/${CONFIG_NAME}
TLE_PATH= ${TLE_FILE_DIR}/${TLE_FILE_NAME}
TEMPLATE_PATH= ${TEMPLATE_DIR}/${TEMPLATE_NAME}
KUBENET_CONFIG_PATH= ${KUBENET_CONFIG_FILE_DIR}/${KUBENET_CONFIG_FILE_NAME}
MERGE_OUTPUT_PATH = ${KUBENET_CONFIG_FILE_DIR}/${MERGE_OUTPUT_NAME}

# --- generate multiple shells ---

.PHONY: gen-all tle-gen-all kubenet-config-gen-all

gen-all: kubenet-config-gen-all merge

tle-gen-all: $(SHELLS:%=tle-gen-%)

kubenet-config-gen-all: $(SHELLS:%=kubenet-config-gen-%)

tle-gen-%: tle-file-dir
	go run ./main.go generate walker tle -c $(CONFIG_PATH) -s $* -o $(TLE_FILE_DIR)/s$*.tle

kubenet-config-gen-%: tle-gen-% kubenet-file-dir
	go run ./main.go generate walker kubenet -c $(CONFIG_PATH) -s $* -t $(TEMPLATE_PATH) -i $(TLE_FILE_DIR)/s$*.tle -o $(KUBENET_CONFIG_FILE_DIR)/s$*.yaml

merge:
	go run ./main.go merge -i $(KUBENET_CONFIG_FILE_DIR) -o $(MERGE_OUTPUT_PATH)

# --- generate one shell ---

.PHONY: gen-one tle-gen kubenet-config-gen

gen-one: kubenet-config-gen

tle-gen: tle-file-dir 
	go run ./main.go generate walker tle -c ${CONFIG_PATH} -s ${SHELL_INDEX} -o ${TLE_PATH} 

kubenet-config-gen: tle-gen kubenet-file-dir 
	go run ./main.go generate walker kubenet -c ${CONFIG_PATH} -s ${SHELL_INDEX} -t ${TEMPLATE_PATH} -i ${TLE_PATH} -o ${KUBENET_CONFIG_PATH}

# --- make dir ---

.PHONY: tle-file-dir kubenet-file-dir
tle-file-dir:
	@mkdir -p $(TLE_FILE_DIR)

kubenet-file-dir:
	@mkdir -p $(KUBENET_CONFIG_FILE_DIR)

# --- clean ---

.PHONY: clean-output
clean-output:
	rm -rf $(TLE_FILE_DIR) $(KUBENET_CONFIG_FILE_DIR)