PACTICIPANT := "pactflow-example-consumer-go-kafka"
GITHUB_REPO := "pactflow/example-provider-js-kafka"
PACT_CLI="docker run --rm -v ${PWD}:${PWD} -e PACT_BROKER_BASE_URL -e PACT_BROKER_TOKEN pactfoundation/pact-cli:latest"
PACT_GO_VERSION=2.0.8
PACT_DOWNLOAD_DIR=/tmp
ifeq ($(OS),Windows_NT)
	PACT_DOWNLOAD_DIR=$$TMP
endif

# Only deploy from master
ifeq ($(GIT_BRANCH),master)
	DEPLOY_TARGET=deploy
else
	DEPLOY_TARGET=no_deploy
endif

all: test

## ====================
## CI tasks
## ====================

ci: test can_i_deploy $(DEPLOY_TARGET)

# Run the ci target from a developer machine with the environment variables
# set as if it was on GitHub Actions.
# Use this for quick feedback when playing around with your workflows.
fake_ci: .env
	CI=true \
	GIT_COMMIT=`git rev-parse --short HEAD`+`date +%s` \
	GIT_BRANCH=`git rev-parse --abbrev-ref HEAD` \
	PACT_BROKER_PUBLISH_VERIFICATION_RESULTS=true \
	make ci

## =====================
## Build/test tasks
## =====================

run:
	go run .

test: .env
	go test -v -count=1 .

## =====================
## Deploy tasks
## =====================

deploy: deploy_app record_deployment

no_deploy:
	@echo "Not deploying as not on master branch"

can_i_deploy: .env
	@echo "can_i_deploy"
	@"${PACT_CLI}" broker can-i-deploy \
	  --pacticipant ${PACTICIPANT} \
	  --version ${GIT_COMMIT} \
	  --to-environment production \
	  --retry-while-unknown 0 \
	  --retry-interval 10

deploy_app:
	@echo "Deploying to prod"

record_deployment: .env
	@"${PACT_CLI}" broker record-deployment --pacticipant ${PACTICIPANT} --version ${GIT_COMMIT} --environment production

## ======================
## Misc
## ======================

.env:
	touch .env

install:
	go install github.com/pact-foundation/pact-go/v2@v$(PACT_GO_VERSION)
	pact-go -l DEBUG install --libDir $(PACT_DOWNLOAD_DIR);

install_cli:
	@if [ ! -d pact/bin ]; then\
		echo "--- 🐿 Installing Pact CLI dependencies"; \
		cd /tmp; \
		curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash -x; \
		echo "/tmp/pact/bin" >> ${GITHUB_PATH}; \
  fi