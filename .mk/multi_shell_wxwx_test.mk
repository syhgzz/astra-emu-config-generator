MULTI_SHELL_WXWX_DIR = ${KUBENET_FILE_DIR}/multi_shell_wxwx_tests
MULTI_SHELL_WXWX_TEMPLATE = ${TEMPLATE_DIR}/wxwx_template.yaml

MULTI_SHELL_WXWX_DIR:
	mkdir -p ${MULTI_SHELL_WXWX_DIR}

s1-tle: tle-file-dir
	go run ./main.go generate walker tle -c 1 -N 256 -P 16 -E 16 -I 82.6 -H 400 -F 14 -D 100 -o ${TLE_FILE_DIR}/s1.tle -r wxwx-1

s2-tle: tle-file-dir
	go run ./main.go generate walker tle -c 2 -N 4096 -P 64 -E 64 -I 60 -H 550 -F 30 -D 100 -o ${TLE_FILE_DIR}/s2.tle -r wxwx-1

s3-tle: tle-file-dir
	go run ./main.go generate walker tle -c 3 -N 224 -P 16 -E 16 -I 60 -H 1200 -F 10 -D 100 -o ${TLE_FILE_DIR}/s3.tle -r wxwx-1

s1: s1-tle MULTI_SHELL_WXWX_DIR
	go run ./main.go generate walker kubenet -c 1 -N 256 -P 16 -E 16 -F 14 --cross-sphere true \
	--fake --accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 --networkDeviceModel ndm1 \
	--offset -1 -i ${TLE_FILE_DIR}/s1.tle -o ${MULTI_SHELL_WXWX_DIR}/s1.yaml -t ${MULTI_SHELL_WXWX_TEMPLATE} -r 1
	@cp ${MULTI_SHELL_WXWX_DIR}/s1.yaml ${OUTPUT_TEST_DIR}/s1.yaml

s2: s2-tle MULTI_SHELL_WXWX_DIR
	go run ./main.go generate walker kubenet -c 1 -N 4096 -P 64 -E 64 -F 30 --cross-sphere true \
	--fake --accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 --networkDeviceModel ndm1 \
	--offset -1 -i ${TLE_FILE_DIR}/s2.tle -o ${MULTI_SHELL_WXWX_DIR}/s2.yaml -t ${MULTI_SHELL_WXWX_TEMPLATE} -r 2
	@cp ${MULTI_SHELL_WXWX_DIR}/s2.yaml ${OUTPUT_TEST_DIR}/s2.yaml

s3: s3-tle MULTI_SHELL_WXWX_DIR
	go run ./main.go generate walker kubenet -c 1 -N 224 -P 16 -E 16 -F 10 --cross-sphere true \
	--fake --accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 --networkDeviceModel ndm1 \
	--offset -1 -i ${TLE_FILE_DIR}/s3.tle -o ${MULTI_SHELL_WXWX_DIR}/s3.yaml -t ${MULTI_SHELL_WXWX_TEMPLATE} -r 3
	@cp ${MULTI_SHELL_WXWX_DIR}/s3.yaml ${OUTPUT_TEST_DIR}/s3.yaml

multi-shell: s1 s2 s3 $(OUTPUT_TEST_DIR)
	go run ./main.go merge -i ${MULTI_SHELL_WXWX_DIR} -o $(OUTPUT_TEST_DIR)/wxwx.yaml