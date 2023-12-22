# Makefile

# Set project ID and image name
GCP_PROJECT_ID=growth-squad-396607
IMAGE=cart-backend
LATEST_TAG=latest
GCR_HOSTNAME=asia.gcr.io

.PHONY: docker-compose push-bot push-service configure-docker deploy build-bot run-bot

run-bot:
	docker run -it --rm --network host --name $(BOT_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(LATEST_TAG)
run-service:
	docker run -it --rm --network host --name $(SERVICE_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(SERVICE_IMAGE):$(LATEST_TAG)
run-db:
	docker run -it --rm --name postgres -e POSTGRES_PASSWORD=user -e POSTGRES_USER=user -p 5432:5432 postgres:14.9-alpine


docker-compose:
	docker-compose -f docker-compose.yaml up

build-image:
	docker build --build-arg ENV=$(env) -f dockerfile -t $(IMAGE) .
	docker tag $(IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE):$(LATEST_TAG)


# Push the Docker image to Google Container Registry
push-service: configure-docker build-image
	@echo "Pushing Service Docker image to GCR..."
	docker tag $(IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE):$(LATEST_TAG)
	docker push $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE):$(LATEST_TAG)

# Authenticate Docker to GCR
configure-docker:
	@echo "Configuring Docker for GCR..."
	gcloud auth configure-docker $(GCR_HOSTNAME)

# Deploy the Docker image to a GCP VM
deploy:
	@echo "Deploying Docker image to GCP VM..."
	gcloud auth configure-docker $(GCR_HOSTNAME) && \
	docker pull $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE):$(LATEST_TAG) && \
	docker stop $(IMAGE) || true && \
	docker rm $(IMAGE) || true && \
	docker run -d --name $(IMAGE) -v certs:/var/www/.cache -p 80:80 -p 443:443 $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE):$(LATEST_TAG)