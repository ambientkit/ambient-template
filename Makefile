# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make env
# * make storage
# * make run-env

# Load the environment variables.
-include .env

.PHONY: default
default: run

################################################################################
# Setup app
################################################################################

.PHONY: env
env:
	@echo Generating .env file.
	@cp testdata/.env .env

.PHONY: privatekey
privatekey:
	# Session key to encrypt the cookie store.
	@GOBIN=$(shell pwd)/bin go install github.com/ambientkit/plugin/pkg/uuid/cmd/privatekey
	@./bin/privatekey

# Pass in ARGS.
# https://stackoverflow.com/a/14061796
ifeq (passhash,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

.PHONY: passhash
passhash:
	# Password hash that is base64 encoded.
	@GOBIN=$(shell pwd)/bin go install github.com/ambientkit/plugin/pkg/passhash/cmd/passhash
	@./bin/passhash ${ARGS}

.PHONY: storage
storage:
	@echo Creating session and site storage files locally.
	cp testdata/storage/session.bin storage/session.bin
	cp testdata/storage/site.bin storage/site.bin

################################################################################
# Build App
################################################################################

.PHONY: run-env
run-env:
	@echo Starting local server with .env.
	AMB_DOTENV=true go run cmd/myapp/main.go

.PHONY: run-local
run-local:
	@echo Starting local server without .env.
	go run cmd/myapp/main.go

.PHONY: buildrun
buildrun: build run

# Build the docker image.
.PHONY: build
build:
	docker build -t ${APPNAME}:$(shell cat VERSION.md) .

# Run and exec into the docker container.
.PHONY: exec
exec:
	docker run --rm --env-file .env -v ${shell pwd}/storage:/app/storage -it ${APPNAME}:$(shell cat VERSION.md) sh

# Run a new docker container.
.PHONY: run
run:
	docker run --rm --env-file .env -v ${shell pwd}/storage:/app/storage -p 8080:8080 ${APPNAME}:$(shell cat VERSION.md)

################################################################################
# Tools
################################################################################

# Update Go dependencies.
.PHONY: update
update:
	go get -u -f -d ./...
	go mod tidy

# Pass in ARGS.
# https://stackoverflow.com/a/14061796
ifeq (update-ambient,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

# Update Ambient dependency.
.PHONY: update-ambient
update-ambient:
	go get github.com/ambientkit/ambient@${ARGS}
	go mod tidy -compat=1.17

# Pass in ARGS.
# https://stackoverflow.com/a/14061796
ifeq (update-plugin,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

# Update Ambient plugin dependency.
.PHONY: update-plugin
update-plugin:
	go get github.com/ambientkit/plugin@${ARGS}
	go mod tidy -compat=1.17

# Install swagger to project bin folder to allow generating a Swagger spec from code.
.PHONY: swagger-install
swagger-install:
	curl -o ./bin/swagger -L "https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_darwin_amd64"
	chmod +x ./bin/swagger

# Generate swagger spec and replace with correct version.
.PHONY: swagger
swagger:
	./bin/swagger generate spec -o swagger.json cmd/myapp/main.go
	sed "s/VERSIONPLACEHOLDER/${shell cat VERSION.md}/g" swagger.json > tmp.json && mv tmp.json swagger.json

# Serve the swagger spec.
.PHONY: swagger-serve
swagger-serve:
	./bin/swagger serve -F swagger swagger.json

# Install air to project bin folder to allow live rebuilding so you can make code changes quickly.
.PHONY: air-install
air-install:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

# Run air.
.PHONY: air
air:
	./bin/air

################################################################################
# Deploy app to Google Cloud
################################################################################

.PHONY: gcp-init
gcp-init:
	@echo Creating the initial files in Google Cloud Storage.
	gsutil mb -p $(AMB_GCP_PROJECT_ID) -l ${AMB_GCP_REGION} -c Standard gs://${AMB_GCP_BUCKET}
	gsutil versioning set on gs://${AMB_GCP_BUCKET}
	gsutil cp storage/initial/site.bin gs://${AMB_GCP_BUCKET}/storage/site.bin
	gsutil cp storage/initial/session.bin gs://${AMB_GCP_BUCKET}/storage/session.bin

.PHONY: gcp-deploy
gcp-deploy:
	@echo Deploying to Google Cloud Run.
	gcloud builds submit --tag gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE}
	gcloud run deploy --image gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE} \
		--platform managed \
		--allow-unauthenticated \
		--region ${AMB_GCP_REGION} ${AMB_GCP_CLOUDRUN_NAME} \
		--update-env-vars AMB_SESSION_KEY=${AMB_SESSION_KEY} \
		--update-env-vars AMB_PASSWORD_HASH=${AMB_PASSWORD_HASH} \
		--update-env-vars AMB_GCP_PROJECT_ID=${AMB_GCP_PROJECT_ID} \
		--update-env-vars AMB_GCP_BUCKET=${AMB_GCP_BUCKET}

.PHONY: gcp-delete
gcp-delete:
	@echo Removing files from Google Cloud.
	-gcloud run services delete --platform managed --region ${AMB_GCP_REGION} ${AMB_GCP_CLOUDRUN_NAME}
	-gsutil -m rm -r -f gs://${AMB_GCP_BUCKET}

################################################################################
# Deploy app to AWS
################################################################################

.PHONY: aws-init
aws-init:
	@echo Creating the initial files in AWS S3.
ifeq "${AWS_REGION}" "us-east-1"
	aws s3api create-bucket --bucket ${AMB_AWS_BUCKET}
else
	aws s3api create-bucket --bucket ${AMB_AWS_BUCKET} --create-bucket-configuration '{"LocationConstraint": "${AWS_REGION}"}'
endif
	aws s3api put-public-access-block --bucket ${AMB_AWS_BUCKET} --public-access-block-configuration '{"BlockPublicAcls": true,"IgnorePublicAcls": true,"BlockPublicPolicy": true,"RestrictPublicBuckets": true}'
	aws s3 cp storage/initial/site.bin s3://${AMB_AWS_BUCKET}/storage/site.bin
	aws s3 cp storage/initial/session.bin s3://${AMB_AWS_BUCKET}/storage/session.bin

.PHONY: aws-deploy
aws-deploy:
	@echo Deploying to AWS App Runner.
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AMB_AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com
	-aws ecr create-repository --repository-name ${AMB_GCP_IMAGE}
	docker build -t ${AMB_AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AMB_GCP_IMAGE}:${AMB_APP_VERSION} .
	docker push ${AMB_AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AMB_GCP_IMAGE}:${AMB_APP_VERSION}
	-aws cloudformation create-stack --stack-name ${AMB_GCP_CLOUDRUN_NAME} \
		--template-body file://deploy/aws-apprunner.json --capabilities CAPABILITY_IAM \
		--parameters ParameterKey=ParameterSessionKey,ParameterValue=${AMB_SESSION_KEY} \
		ParameterKey=ParameterPasswordHash,ParameterValue=${AMB_PASSWORD_HASH} \
		ParameterKey=ParameterAWSS3Bucket,ParameterValue=${AMB_AWS_BUCKET} \
		ParameterKey=ParameterAWSECRName,ParameterValue=${AMB_GCP_IMAGE} \
		ParameterKey=ParameterAppVersion,ParameterValue=${AMB_APP_VERSION}
	-aws cloudformation update-stack --stack-name ${AMB_GCP_CLOUDRUN_NAME} \
		--template-body file://deploy/aws-apprunner.json --capabilities CAPABILITY_IAM \
		--parameters ParameterKey=ParameterSessionKey,ParameterValue=${AMB_SESSION_KEY} \
		ParameterKey=ParameterPasswordHash,ParameterValue=${AMB_PASSWORD_HASH} \
		ParameterKey=ParameterAWSS3Bucket,ParameterValue=${AMB_AWS_BUCKET} \
		ParameterKey=ParameterAWSECRName,ParameterValue=${AMB_GCP_IMAGE} \
		ParameterKey=ParameterAppVersion,ParameterValue=${AMB_APP_VERSION}

.PHONY: aws-delete
aws-delete:
	@echo Removing files from AWS.
	-aws cloudformation delete-stack --stack-name ${AMB_GCP_CLOUDRUN_NAME}
	-aws ecr delete-repository --repository-name ${AMB_GCP_IMAGE} --force
	-aws s3 rm s3://${AMB_AWS_BUCKET} --recursive
	-aws s3api delete-bucket --bucket ${AMB_AWS_BUCKET}

################################################################################
# Deploy app to Azure
################################################################################

.PHONY: azcli-start
azcli-start:
	@echo Starting Azure CLI in docker container.
	# Run docker in the background
	docker run -d -t --name azurecli -v $(shell pwd):/root mcr.microsoft.com/azure-cli

.PHONY: azcli-stop
azcli-stop:
	@echo Stopping Azure CLI in docker container.
	docker rm -f azurecli

.PHONY: az-init
az-init:
	@echo Creating the initial files in Azure storage.
	./deploy/azure-init.sh

.PHONY: az-deploy
az-deploy:
	@echo Deploying to Azure Functions.
	./deploy/azure-deploy.sh

.PHONY: az-delete
az-delete:
	@echo Removing Function and storage account from Azure.
#	-az functionapp delete --resource-group ${AMB_AZURE_RESOURCE_GROUP} --name ${AMB_AZURE_FUNCTION_NAME}
#	-az storage account delete --resource-group ${AMB_AZURE_RESOURCE_GROUP} --name ${AZURE_STORAGE_ACCOUNT}
	-az group delete --name ${AMB_AZURE_RESOURCE_GROUP} --yes