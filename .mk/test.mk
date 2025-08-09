NN ?= kubenet/$(TEST)
DATA ?= ./tests/output/$(TEST)

SERVER ?= http://localhost:9090
# SERVER ?= http://10.0.0.31:31002

.PHONY: exp-apply exp-delete exp-start exp-stop

check-test:
	@if [ -z "$(TEST)" ]; then \
		echo "❌ ERROR: TEST environment variable is not set!"; \
		exit 1; \
	fi
	@echo "✅ TEST is $(TEST)"

exp-apply:
	curl -X POST -H "Content-Type: multipart/form-data" \
    -F "file=@${DATA}" \
    "$(SERVER)/manage/network_experiment/apply?nn=${NN}"

exp-delete:
	curl -X POST "$(SERVER)/manage/network_experiment/delete?nn=${NN}"

exp-start:
	curl -X POST "$(SERVER)/manage/scene/start?nn=${NN}"

exp-stop:
	curl -X POST "$(SERVER)/manage/scene/stop?nn=${NN}"

exp-reset:
	curl -X POST "$(SERVER)/manage/scene/reset?nn=${NN}"

exp-flow-start:
	curl -X POST "$(SERVER)/manage/flow/start?nn=${NN}"

exp-flow-stop:
	curl -X POST "$(SERVER)/manage/flow/stop?nn=${NN}"