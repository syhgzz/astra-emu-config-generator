NN ?= kubenet/$(TEST)
DATA ?= $(OUTPUT_TEST_DIR)/$(TEST).yaml

SERVER ?= http://localhost:9090
# SERVER ?= http://10.0.0.10:9090
# SERVER ?= http://10.0.0.31:31002

.PHONY: exp-apply exp-delete exp-start exp-stop

check-test:
	@if [ -z "$(TEST)" ]; then \
		echo "❌ ERROR: TEST environment variable is not set!"; \
		exit 1; \
	fi
	@echo "✅ TEST is $(TEST)"

exp-apply: check-test
	@echo "🎯 Applying experiment $(NN) from $(DATA)..."
	@curl -X POST -H "Content-Type: multipart/form-data" \
    -F "file=@${DATA}" \
    "$(SERVER)/manage/network_experiment/apply?nn=${NN}"

exp-delete: check-test
	@echo "🗑️ Deleting experiment $(NN)..."
	@curl -X POST "$(SERVER)/manage/network_experiment/delete?nn=${NN}"

exp-start: check-test
	@echo "🚀 Starting experiment $(NN)..."
	@curl -X POST "$(SERVER)/manage/scene/start?nn=${NN}"

exp-stop: check-test
	@echo "🛑 Stopping experiment $(NN)..."
	@curl -X POST "$(SERVER)/manage/scene/stop?nn=${NN}"

exp-reset: check-test
	@echo "🔄 Resetting experiment $(NN)..."
	@curl -X POST "$(SERVER)/manage/scene/reset?nn=${NN}"

exp-flow-start: check-test
	@echo "🚦 Starting flow for experiment $(NN)..."
	@curl -X POST "$(SERVER)/manage/flow/start?nn=${NN}"

exp-flow-stop: check-test
	@echo "🚦 Stopping flow for experiment $(NN)..."
	@curl -X POST "$(SERVER)/manage/flow/stop?nn=${NN}"