MULTI_SHELL_BGP_DIR = ${KUBENET_FILE_DIR}/multi_shell_bgp_tests
MULTI_SHELL_BGP_TEMPLATE = ${TEMPLATE_DIR}/bgp_template.yaml

multi_shell_bgp_dir:
	mkdir -p ${MULTI_SHELL_BGP_DIR}

s1-tle: tle-file-dir
	go run ./main.go tle -N 800 -P 20 -E 20 -I 60 -H 1000 -F 1 -D 100 -o ${TLE_FILE_DIR}/s1.tle -n S1

s2-tle: tle-file-dir
	go run ./main.go tle -N 560 -P 35 -E 35 -I 55 -H 1000 -F 0 -D 100 -o ${TLE_FILE_DIR}/s2.tle -n S2

s3-tle: tle-file-dir
	go run ./main.go tle -N 648 -P 36 -E 36 -I 89 -H 1500 -F 0 -D 100 -o ${TLE_FILE_DIR}/s3.tle -n S3

s1: s1-tle multi_shell_bgp_dir
	go run ./main.go generate -N 800 -P 20 -E 20 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	--fake true --offset -1 -i ${TLE_FILE_DIR}/s1.tle -o ${MULTI_SHELL_BGP_DIR}/s1.yaml -t ${MULTI_SHELL_BGP_TEMPLATE} -r 1 -n S1
	@cp ${MULTI_SHELL_BGP_DIR}/s1.yaml ${OUTPUT_TEST_DIR}/s1.yaml

s2: s2-tle multi_shell_bgp_dir
	go run ./main.go generate -N 560 -P 35 -E 35 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 \
	--offset 0 -i ${TLE_FILE_DIR}/s2.tle -o ${MULTI_SHELL_BGP_DIR}/s2.yaml -t ${MULTI_SHELL_BGP_TEMPLATE} -r 2 -n S2
	@cp ${MULTI_SHELL_BGP_DIR}/s2.yaml ${OUTPUT_TEST_DIR}/s2.yaml

s3: s3-tle multi_shell_bgp_dir
	go run ./main.go generate -N 648 -P 36 -E 36 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 \
	--offset 0 -i ${TLE_FILE_DIR}/s3.tle -o ${MULTI_SHELL_BGP_DIR}/s3.yaml -t ${MULTI_SHELL_BGP_TEMPLATE} -r 3 -n S3
	@cp ${MULTI_SHELL_BGP_DIR}/s3.yaml ${OUTPUT_TEST_DIR}/s3.yaml

ground1: multi_shell_bgp_dir
	go run ./main.go ground -N 20 --longitude 100.2 --latitude 19.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	-o ${MULTI_SHELL_BGP_DIR}/ground1.yaml -t ${MULTI_SHELL_BGP_TEMPLATE} -r 1 -n Ground1

ground2: multi_shell_bgp_dir
	go run ./main.go ground -N 40 --longitude 130.2 --latitude 15.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 \
	-o ${MULTI_SHELL_BGP_DIR}/ground2.yaml -t ${MULTI_SHELL_BGP_TEMPLATE} -r 2 -n Ground2
	
ground3: multi_shell_bgp_dir
	go run ./main.go ground -N 140 --longitude 160.2 --latitude 11.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 \
	-o ${MULTI_SHELL_BGP_DIR}/ground3.yaml -t ${MULTI_SHELL_BGP_TEMPLATE} -r 3 -n Ground3

two-shell: s1 s2 ground1 ground2 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_BGP_DIR} -o $(OUTPUT_TEST_DIR)/two_shell_bgp_test.yaml

two-shell-ng: s1 s2 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_BGP_DIR} -o $(OUTPUT_TEST_DIR)/two_shell_ng_bgp_test.yaml

multi-shell: s1 s2 s3 ground1 ground2 ground3 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_BGP_DIR} -o $(OUTPUT_TEST_DIR)/multi-shell.yaml

multi-shell-ng: s1 s2 s3 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_BGP_DIR} -o $(OUTPUT_TEST_DIR)/multi-shell.yaml