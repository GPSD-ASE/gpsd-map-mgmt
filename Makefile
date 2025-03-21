NAMESPACE = gpsd
DEPLOYMENT = gpsd-map-mgmt
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
TAG ?= latest  # If no tag is provided, default to 'latest'
LOCAL_CHART_NAME = helm
REMOTE_CHART_REPOSITORY = gpsd-ase.github.io
SERVICE_NAME = $(DEPLOYMENT)
TAG_VERSION=$(shell echo $(TAG) | sed 's/^v//')

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

deploy-gh-pages: gh-pages-publish helm-repo-update

gh-pages-publish:
	@echo "Publishing Helm chart for $(SERVICE_NAME) to GitHub Pages..."
	rm -rf /tmp/gpsd-* /tmp/index.yaml
	helm package ./$(LOCAL_CHART_NAME) -d /tmp 
	helm repo index /tmp --url https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/ --merge /tmp/index.yaml
	git fetch origin gh-pages
	git checkout gh-pages
	ls -l /tmp
	cp /tmp/$(SERVICE_NAME)-0.0.15.tgz /tmp/index.yaml .
	git add .
	git commit -m "fix: commit to update Github Pages"
	git push origin gh-pages -f
	sleep 5
	curl -k https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/index.yaml

helm-repo-update:
	@echo "Adding and updating Helm repo for $(SERVICE_NAME)..."
	helm repo add $(SERVICE_NAME) https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/
	helm repo update
	helm repo list