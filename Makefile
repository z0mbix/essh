.DEFAULT_GOAL := help
.SHELL := bash

start-instance: ## Stop test instance
	@aws ec2 start-instances --instance-id $(AWS_INSTANCE_ID)

stop-instance: ## Stop test instance
	@aws ec2 stop-instances --instance-id $(AWS_INSTANCE_ID)

instance-status: ## Get test instance status
	@aws ec2 describe-instances \
		--output json \
		--instance-id $(AWS_INSTANCE_ID) \
		| jq -rC '.Reservations[0].Instances[0].State.Name'

release: ## Create a new Github release with goreleaser
	goreleaser --rm-dist

help: ## See all the Makefile targets
	@awk 'BEGIN {FS = ":.*##"; \
		printf "Usage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: start-instance stop-instance instance-status release help
