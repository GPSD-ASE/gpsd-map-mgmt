build-image:
	docker build -f docker/Dockerfile -t gpsd/gpsd-map-mgmt:v1 .

push-image:
	docker push gpsd/gpsd-map-mgmt:v1

run-image:
	docker run -p 7000:7000 -it gpsd/gpsd-map-mgmt:v1

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

build:
	kubectl create namespace gpsd || true

setup:
	kubectl apply -f deployments/map-mgmt-deployment.yaml
	kubectl apply -f services/map-mgmt-service.yaml

all: build-image push-image build setup

clean:
	kubectl delete all --all -n gpsd || true

	kubectl delete namespace gpsd || true
	sleep 2