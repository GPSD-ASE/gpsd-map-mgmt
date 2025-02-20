IMAGE_NAME = gpsd/gpsd-map-mgmt
TAG ?= latest  # If no tag is provided, default to 'latest'

build-image:
	docker build -f docker/Dockerfile -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)
