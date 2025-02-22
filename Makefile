NAMESPACE = gpsd
DEPLOYMENT = gpsd-map-mgmt
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
TAG ?= latest  # If no tag is provided, default to 'latest'

# Use `make develop` for local testing
develop: helm-uninstall build-image push-image helm

docker: build-image push-image

build-image:
	docker build -f docker/Dockerfile --tag $(IMAGE_NAME):$(TAG) --platform linux/amd64 .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	docker run -p 7000:7000 $(IMAGE_NAME):$(TAG)

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

helm:
	helm upgrade --install demo ./helm --set image.tag=$(TAG) --namespace $(NAMESPACE)

helm-uninstall:
	helm uninstall demo -n $(NAMESPACE) 

clean:
	kubectl delete all --all -n $(NAMESPACE)  || true
	kubectl delete namespace $(NAMESPACE)  || true
	sleep 2