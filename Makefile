IMAGE_NAME = gpsd/gpsd-map-mgmt
TAG ?= latest  # If no tag is provided, default to 'latest'

build-image:
	docker build -f docker/Dockerfile -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	docker run -d -p 7000:7000 --name test_container $(IMAGE_NAME):$(TAG)

clean-image:
	docker stop test_container || true
	docker rm test_container || true

