PACTICIPANT := "pactflow-example-provider-js-kafka"
GITHUB_REPO := "pactflow/example-provider-js-kafka"
PACT_CHANGED_WEBHOOK_UUID := "0278fe46-09be-4b9d-b037-559d2891f752"
PACT_CLI="docker run --rm -v ${PWD}:${PWD} -e PACT_BROKER_BASE_URL -e PACT_BROKER_TOKEN pactfoundation/pact-cli:latest"

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
# set as if it was on Github Actions.
# Use this for quick feedback when playing around with your workflows.
fake_ci: .env
	CI=true \
	GIT_COMMIT=`git rev-parse --short HEAD`+`date +%s` \
	GIT_BRANCH=`git rev-parse --abbrev-ref HEAD` \
	PACT_BROKER_PUBLISH_VERIFICATION_RESULTS=true \
	make ci

ci_webhook: .env
	npm run test

fake_ci_webhook:
	CI=true \
	GIT_COMMIT=`git rev-parse --short HEAD`+`date +%s` \
	GIT_BRANCH=`git rev-parse --abbrev-ref HEAD` \
	PACT_BROKER_PUBLISH_VERIFICATION_RESULTS=true \
	make ci_webhook

## =====================
## Build/test tasks
## =====================

test: .env
	npm run test

## =====================
## Deploy tasks
## =====================

deploy: deploy_app record_deployment

no_deploy:
	@echo "Not deploying as not on master branch"

can_i_deploy: .env
	"${PACT_CLI}" broker can-i-deploy --pacticipant ${PACTICIPANT} --version ${GIT_COMMIT} --to-environment production

deploy_app:
	@echo "Deploying to production"

record_deployment: .env
	@"${PACT_CLI}" broker record-deployment --pacticipant ${PACTICIPANT} --version ${GIT_COMMIT} --environment production

## ======================
## Misc
## ======================

.env:
	touch .env

docker-logs:
	@docker-compose -f kafka-cluster.yml logs -f

docker-rm:
	@docker-compose -f kafka-cluster.yml rm -vfs

docker-stop:
	@docker-compose -f kafka-cluster.yml stop

docker:
	@docker-compose -f kafka-cluster.yml up -d --no-recreate

start: docker
	npm start

create-product:
	curl -v -X POST -H"Content-Type: application/json" localhost:8001/products -d '{"type":"pizza", "name":"food"}'

update-product:
	curl -v -X PUT -H"Content-Type: application/json" localhost:8001/products/1 -d '{"type":"pizza", "name":"food", "version":"v1"}'

delete-product:
	curl -v -X DELETE -H"Content-Type: application/json" localhost:8001/products/1 -d '{"type":"pizza", "name":"food", "version":"v1"}'

.PHONY: docker docker-stop docker-rm docker-logs
