# Secure-Access-Cloud Kubernetes Operator

## Introduction
Create and Govern the access to your K8s resources the same way you create them.   
[Secure-access-cloud](https://www.broadcom.com/products/cyber-security/network/web-protection/secure-access-cloud) custom Kubernetes controller makes it easy to expose Kubernetes services through SAC using K8s [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

## Usage

Currently supporting 2 CRDs:

1. Site
2. HTTP application

## Installing

1. In Secure-Access-Cloud admin portal:
 - Create an API Client - settings->API Clients->New
 - Assign API client to tenant roles - settings->Tenant Roles->Assign Roles->Tenanat Admin->Click add under the API Client created above
   ![create-api!](assets/create-api.gif "create admin api")

2. In the desired K8s cluster, create a generic secret that the operator will be using to access Secure-Access-Cloud
```shell
>> kubectl create namespace secure-access-cloud-system --save-config
>> kubectl -n secure-access-cloud-system create secret generic secure-access-cloud-config \
--from-literal="tenantDomain=<api endpoint>" \
--from-literal="clientId=<client Id>" \
--from-literal="clientSecret=<client secret>"
```
api endpoint = your tenant URL
Client Id = from step 1.
Client Secret = from step 1.


2. Clone the repository
```shell
>> git clone git@github.com:odedpriva/sac-operator.git
>> cd sac-operator
```

3. Run the following commands:
```shell
>> make install ## Install CRDs into the K8s cluster specified in ~/.kube/config. 
>> make deploy ## Deploy controller to the K8s cluster specified in ~/.kube/config.
```

3. Create site
- check the site [sample](config/samples/site.yaml)

4. Create application
- check the site [sample](config/samples/http-application.yaml)

## Uninstall

1. delete all resources related to the operator.
```shell
>> make undeploy
```

## [Contributing](contributing.md)
