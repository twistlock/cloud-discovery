# Twistlock Cloud Discovery

Cloud Discovery provides a point in time enumeration of all the cloud native platform services used in a given cloud provider account.
It is useful to audit and security practitioners that want a simple way to discover all the 'unknown unknowns' in an environment without having
to manually login to multiple provider consoles, click through many pages, and manually export the data.

Provided by [Twistlock](https://www.twistlock.com).

<img src="http://www.twistlock.com/wp-content/uploads/2017/11/Twistlock_Logo-Lockup_RGB.png" width="400">

# Environment variables

1. BASIC_AUTH_USERNAME - This variable determines the username to use for basic authentication.
2. BASIC_AUTH_PASSWORD - This variable determines the password to use for basic authentication.
3. TLS_CERT_PATH - This variable determines the path to the TLS certificate inside the container.
   By default the service generate self-signed certificates for localhost usage.
4. TLS_CERT_KEY - This variable determines the path to the TLS certificate key inside the container.

# Example usages

## Start the container

```sh
docker run -d --name cloud-discovery --restart=always \
 -e BASIC_AUTH_USERNAME=admin -e BASIC_AUTH_PASSWORD=pass -e PORT=9083 -p 9083:9083  twistlock/cloud-discovery
```

## Scan for AWS assets
Discover all AWS assets: lambda, ECS, EKS and ECR.
```sh
docker run -d --name cloud-discovery --restart=always \
 -e BASIC_AUTH_USERNAME=admin -e BASIC_AUTH_PASSWORD=pass -e PORT=9083 -p 9083:9083  twistlock/cloud-discovery
```

## Scan for AWS assets (partial data)
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

## Scan for AWS assets (full data)
```sh
curl -k -v -u admin:pass --raw --data \
'{"credentials": [{"id":"<AWS_ACCESS_KEY>","secret":"<AWS_ACCESS_PASSWORD>"}]}' https://localhost:9083/discover?format=json
```

## Scan for open ports and check for insecure apps.
Scan all open ports and automatically detects insecure apps (native cloud apps configured without proper authorization)
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



