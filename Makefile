IMAGE_REPO=docker-unj-repo.softplan.com.br/unj/benchmark-app
TAG=latest

#Build the binary
build: pre-build
	@echo "Building benchmark-app"
	GOOS_VAL=$(shell go env GOOS) GOARCH_VAL=$(shell go env GOARCH) go build -a -installsuffix cgo -o .//build/benchmark-app ./*.go

#Build the image
container-image:
	@echo "Building docker image"
	@docker build --build-arg GOOS_VAL=$(shell go env GOOS) --build-arg GOARCH_VAL=$(shell go env GOARCH) -t $(IMAGE_REPO) -f Dockerfile --no-cache .
	@echo "Docker image build successfully"

#Pre-build checks
pre-build:
	@echo "Checking system information"
	@if [ -z "$(shell go env GOOS)" ] || [ -z "$(shell go env GOARCH)" ] ; then echo 'ERROR: Could not determine the system architecture.' && exit 1 ; fi

#Tag images
tag-image: 
	@echo 'Tagging docker image'
	@docker tag $(IMAGE_REPO) $(IMAGE_REPO):$(TAG)

#Docker push image
publish:
	@echo "Pushing docker image to repository"
	@docker login docker-unj-repo.softplan.com.br
	@docker push $(IMAGE_REPO):$(TAG)