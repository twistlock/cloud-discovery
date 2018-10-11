# Twistlock Cloud Discovery

Cloud Discovery provides point in time enumeration of all the cloud native platform services, such as container registries, managed Kubernetes platforms, and serverless services used across your cloud providers, accounts, and regions.  Its a powerful tool for audit and security practitioners that want a simple way to discover all the 'unknown unknowns' across environments without having to manually login to multiple provider consoles, click through many pages, and manually export the data.  

Cloud Discovery connects to cloud providers' native platform APIs to discover services and their metadata and requires only read permissions.  Cloud Discovery also has a network discovery option that uses port scanning to sweep IP ranges and discover cloud native infrastructure and apps, such as Docker Registries and Kubernetes API servers, with weak settings or authentication.  This is useful to discover 'self-installed' cloud native components not provided as a service by a cloud provider, such as a Docker Registry running on an EC2 instance.  Cloud Discovery is provided as a simple Docker container image that can be run anywhere and works well for both interactive use and automation.

Cloud Discovery is [another](https://github.com/docker/swarmkit/pull/2239) [open](https://github.com/moby/moby/pull/15365) [source](https://github.com/moby/moby/pull/20111) [contribution](https://github.com/moby/moby/pull/21556) [provided](https://github.com/docker/distribution/pull/2362 ) by [Twistlock](https://www.twistlock.com).

<img src="http://www.twistlock.com/wp-content/uploads/2017/11/Twistlock_Logo-Lockup_RGB.png" width="400">

# Environment variables

1. BASIC_AUTH_USERNAME - This variable determines the username to use for basic authentication.
2. BASIC_AUTH_PASSWORD - This variable determines the password to use for basic authentication.
3. TLS_CERT_PATH - This variable determines the path to the TLS certificate inside the container.
   By default the service generates self-signed certificates for localhost usage.
4. TLS_CERT_KEY - This variable determines the path to the TLS certificate key inside the container.

# Example usages

## Start the container

```sh
docker run -d --name cloud-discovery --restart=always \
 -e BASIC_AUTH_USERNAME=admin -e BASIC_AUTH_PASSWORD=pass -e PORT=9083 -p 9083:9083  twistlock/cloud-discovery
```

## Scan and list all AWS assets
```sh
curl -k -v -u admin:pass --raw --data \
'{"credentials": [{"id":"<AWS_ACCESS_KEY>","secret":"<AWS_ACCESS_PASSWORD>"}]}' \
 https://localhost:9083/discover
```
Output
```sh
Type    Region        ID
EKS     us-east-1     k8s-cluster-1
ECS     us-east-1     cluster-1
ECS     us-east-1     cluster-2
ECS     us-east-1     cluster-3
ECR     us-east-2     cluster-1
```

## Scan all AWS assets and show full metadata for each of them
```sh
curl -k -v -u admin:pass --raw --data \
'{"credentials": [{"id":"<AWS_ACCESS_KEY>","secret":"<AWS_ACCESS_PASSWORD>"}]}' https://localhost:9083/discover?format=json
```

## Port scan a subnet to discover cloud native infrastructure and apps
Scan all open ports and automatically detect insecure apps (native cloud apps configured without proper authorization)
Remark: If the container runs in AWS cluster, the subnet can be automatically extracted from [AWS metadata API server](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html)
```sh
curl -k -v -u admin:pass --raw   --data '{"subnet":"172.17.0.1", "debug": true}'   https://localhost:9083/nmap
```
Output
```
Host           Port      App                 Insecure
172.17.0.1     5000      docker registry     true
172.17.0.1     5003      docker registry     false
172.17.0.1     27017     mongod              true
```



