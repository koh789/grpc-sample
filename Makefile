.PHONY: help
help: ## help 表示 `make help` でタスクの一覧を確認できます
	@echo "------- タスク一覧 ------"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36mmake %-20s\033[0m %s\n", $$1, $$2}'

.PHONY: protoc
protoc: ## protoc
	 protoc --go_out=./pkg/pb --go_opt=paths=source_relative --go-grpc_out=./pkg/pb  --go-grpc_opt=paths=source_relative  proto/*


