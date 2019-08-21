## Twistlock Cloud Discovery on GCP User Guide

### Overview
- Cloud Discovery provides point in time enumeration of all the cloud native platform services, such as container registries, managed Kubernetes platforms, and serverless services used across your cloud providers, accounts, and regions. It’s a powerful tool for audit and security practitioners that want a simple way to discover all the 'unknown unknowns' across environments without having to manually login to multiple provider consoles, click through many pages, and manually export the data.
- Cloud Discovery connects to cloud providers' native platform APIs to discover services and their metadata and requires only read permissions. Cloud Discovery also has a network discovery option that uses port scanning to sweep IP ranges and discover cloud native infrastructure and apps, such as Docker Registries and Kubernetes API servers, with weak settings or authentication. This is useful to discover 'self-installed' cloud native components not provided as a service by a cloud provider, such as a Docker Registry running on an EC2 instance. Cloud Discovery is provided as a simple Docker container image that can be run anywhere and works well for both interactive use and automation.
- Cloud Discovery is another open source contribution provided by Twistlock.
- [GCP Marketplace](https://console.cloud.google.com/marketplace/details/twistlock/cloud-discovery)

### One-time setup
- No special setup is necessary to use Cloud Discovery other that having a container runtime or orchestrator that can run the Cloud Discovery Docker container image.

### Installation
- Simply pull the Cloud Discovery container image to a machine with a container runtime and run the container. For example:

```sh
docker run -d --name cloud-discovery --restart=always \
-e BASIC_AUTH_USERNAME=admin -e BASIC_AUTH_PASSWORD=pass -e PORT=9083 -p 9083:9083  twistlock/cloud-discovery
```

### Basic Usage
#### Scan and list all GCP assets
```sh
SERVICE_ACCOUNT=$(cat <service_account_secret> | base64 | tr -d '\n')
curl -k -v -u admin:pass --raw --data '{"credentials": [{"secret":"'${SERVICE_ACCOUNT}'", "provider":"gcp"}]}' https://localhost:9083/discover
```
Output
```sh
Type        Region            ID
GKE         us-central1-a     cluster-1
GKE         us-central1-a     cluster-2
GCR         gcr.io            registry-1
GCR         gcr.io            registry-2
Functions   us-central1       function-1
```

#### Scan all GCP assets and show full metadata for each of them
```sh
SERVICE_ACCOUNT=$(cat <service_account_secret> | base64 | tr -d '\n')
curl -k -v -u admin:pass --raw --data '{"credentials": [{"secret":"'${SERVICE_ACCOUNT}'", "provider":"gcp"}]}' https://localhost:9083/discover?format=json
```

### Backup and restore
- Cloud Discovery is stateless and will provide a point in time enumeration of all your GCP resources when it is run.

### Image updates
- Simply pull the latest container available.

### Scaling
- n/a

### Deletion
- Simply stop the container and delete it from your system.

