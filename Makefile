.PHONY: deps clean build test local run iam-config config toml ready deploy delete

deps: # install dependencies (modules)
	go get -v -t -d ./...

clean: # clean up working copy (i.e. delete compiled binaries)
	@rm -rf .aws-sam
	@go clean

build: # compile binaries and prep deployment templates
	sam build

test: # run tests
	go test -v ./...

local run: build test # run api locally
	sam local start-api

iam-config: # one-off IAM configuration stack
	aws cloudformation deploy --stack-name sample-iam-config --template-file ./cfn/iam-config.yaml --capabilities CAPABILITY_IAM

config: # one-off SAM guided configuration
	sam deploy --guided --no-execute-changeset

samconfig.toml:
	@ \
	echo "ERROR: Missing samconfig.toml. SOLUTION: Run [make config]." >&2; \
	exit 1;

# Supported replacement variables for samconfig.toml
AWS_REGION ?= '{AWS_REGION}'
AWS_SAM_S3_HASH ?= '{AWS_SAM_S3_HASH}'
AWS_SAM_S3_BUCKET ?= '{AWS_SAM_S3_BUCKET}'

toml: samconfig.toml # substitute replacement variables in samconfig.toml
	@ \
	if grep '{' samconfig.toml > /dev/null; then \
		sed -i.bak -e 's/{AWS_REGION}/$(AWS_REGION)/g' -e 's/{AWS_SAM_S3_HASH}/$(AWS_SAM_S3_HASH)/g' -e 's/{AWS_SAM_S3_BUCKET}/$(AWS_SAM_S3_BUCKET)/g' ./samconfig.toml; \
	fi

ready: samconfig.toml # check if configuration is complete
	@ \
	if ! grep 's3_bucket' samconfig.toml > /dev/null; then \
		echo "ERROR: Missing s3_bucket parameter in samconfig.toml. SOLUTION: Run [make config]." >&2; \
		exit 2; \
	elif grep '{' samconfig.toml > /dev/null; then \
		echo "ERROR: Replacement variables in samconfig.toml. SOLUTION: Run [make toml]." >&2; \
		exit 3; \
	fi

deploy: ready # deploy CloudFormation stack
	sam deploy --no-fail-on-empty-changeset

delete: # delete CloudFormation stack
	aws cloudformation delete-stack --stack-name sample-serverless-app
	aws cloudformation delete-stack --stack-name sample-iam-config
