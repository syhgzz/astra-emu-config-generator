PHYSIC_DIR = ${KUBENET_FILE_DIR}/physic_tests
PHYSIC_TEMPLATE = ${TEMPLATE_DIR}/physic_template.yaml

physic_dir:
	mkdir -p ${PHYSIC_DIR}

physic-one-tle: tle-file-dir
	go run ./main.go tle -N 100 -P 10 -E 10 -I 85 -H 500 -F 0 -D 100 -o ${TLE_FILE_DIR}/phy1.tle -n S1

physic-one-shell: physic-one-tle $(OUTPUT_TEST_DIR)
	go run ./main.go generate -N 100 -P 10 -E 10 \
	--accessIntfNum 16 --accessFadingmodel fading_model1 --accessLinkmodel Satellite --accessTopologymodel cone1 \
	--offset 0 -F 0 -i ${TLE_FILE_DIR}/phy1.tle -o ${OUTPUT_TEST_DIR}/phy1.yaml -t ${PHYSIC_TEMPLATE} -r 1 -n S1