MULTI_SHELL_OSPF_DIR = ${KUBENET_FILE_DIR}/multi_shell_ospf_tests
MULTI_SHELL_OSPF_TEMPLATE = ${TEMPLATE_DIR}/ospf_template.yaml

multi_shell_ospf_dir:
	mkdir -p ${MULTI_SHELL_OSPF_DIR}

s1-tle: tle-file-dir
	go run ./main.go tle -N 60 -P 10 -E 10 -I 85 -H 1175 -F 0 -D 100 -o {TLE_FILE_DIR}/s1.tle -n S1

s2-tle: tle-file-dir
	go run ./main.go tle -N 108 -P 9 -E 9 -I 55 -H 1150 -F 0 -D 100 -o {TLE_FILE_DIR}/s2.tle -n S2

s3-tle: tle-file-dir
	go run ./main.go tle -N 648 -P 36 -E 36 -I 89 -H 1100 -F 0 -D 100 -o {TLE_FILE_DIR}/s3.tle -n S3

s1: s1-tle multi_shell_ospf_dir
	go run ./main.go generate -N 60 -P 10 -E 10 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	--offset 0 -F 0 -i {TLE_FILE_DIR}/s1.tle -o {KUBENET_FILE_DIR}/multi_shell_ospf_tests/s1.yaml -t {TEMPLATE} -r 1 -n S1

s2: s2-tle multi_shell_ospf_dir
	go run ./main.go generate -N 108 -P 9 -E 9 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 \
	--offset 0 -F 0 -i {TLE_FILE_DIR}/s2.tle -o {KUBENET_FILE_DIR}/multi_shell_ospf_tests/s2.yaml -t {TEMPLATE} -r 2 -n S2

s3: s3-tle multi_shell_ospf_dir
	go run ./main.go generate -N 648 -P 36 -E 36 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 \
	--offset 0 -F 0 -i {TLE_FILE_DIR}/s3.tle -o {KUBENET_FILE_DIR}/multi_shell_ospf_tests/s3.yaml -t {TEMPLATE} -r 3 -n S3

ground1: multi_shell_ospf_dir
	go run ./main.go ground -N 20 --longitude 113.2 --latitude 19.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	-o {KUBENET_FILE_DIR}/multi_shell_ospf_tests/ground1.yaml -t {TEMPLATE} -r 1 -n Ground1

ground2: multi_shell_ospf_dir
	go run ./main.go ground -N 40 --longitude 113.2 --latitude 15.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone2 \
	-o {KUBENET_FILE_DIR}/multi_shell_ospf_tests/ground2.yaml -t {TEMPLATE} -r 2 -n Ground2
	
ground3: multi_shell_ospf_dir
	go run ./main.go ground -N 140 --longitude 113.2 --latitude 11.7 \
	--accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone3 \
	-o {KUBENET_FILE_DIR}/multi_shell_ospf_tests/ground3.yaml -t {TEMPLATE} -r 3 -n Ground3

two-shell: s1 s2 ground1 ground2 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_OSPF_DIR} -o $(OUTPUT_TEST_DIR)/two_shell_ospf_test.yaml

two-shell-ng: s1 s2 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_OSPF_DIR} -o $(OUTPUT_TEST_DIR)/two_shell_ng_ospf_test.yaml

multi-shell: s1 s2 s3 ground1 ground2 ground3 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_OSPF_DIR} -o $(OUTPUT_TEST_DIR)/multi_shell_ospf_test.yaml

multi-shell-ng: s1 s2 s3 output-test-dir
	go run ./main.go merge -i ${MULTI_SHELL_OSPF_DIR} -o $(OUTPUT_TEST_DIR)/multi_shell_ng_ospf_test.yaml