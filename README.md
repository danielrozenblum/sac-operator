# Secure-Access-Cloud Kubernetes Operator

## Introduction
Create and Govern the access to your K8s resources the same way you create them.   
[Secure-access-cloud](https://www.broadcom.com/products/cyber-security/network/web-protection/secure-access-cloud) custom Kubernetes controller makes it easy to expose Kubernetes services through SAC using K8s [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

## Usage

Currently supporting 2 CRDs:

1. Sites
2. Web application

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
api endpoint = your tenant URL Client Id = from step 1. Client Secret = from step 1.


2. Clone the repository
```shell
>> git clone git@github.com:danielrozenblum/sac-operator
>> cd sac-operator
```

3. Run the following commands:
```shell
>> make install ## Install CRDs into the K8s cluster specified in ~/.kube/config. 
>> make deploy ## Deploy controller to the K8s cluster specified in ~/.kube/config.
```

3. Create site
In the desired K8s cluster, apply kind:site .yaml
- Check the site [sample](config/samples/site.yaml)
```shell
>> kubectl apply -f site.yaml namespace secure-access-cloud-system
```

4. Create application
In the desired K8s cluster, apply kind:HttpApplication .yaml
- Check the application [sample](config/samples/http-application.yaml)
```shell
>> kubectl apply -f http-application.yaml namespace secure-access-cloud-system
```


## Uninstall

1. Delete all resources related to the operator.
```shell
>> make undeploy ## Uninstall CRDs and delete namespace secure-access-cloud-system in K8s cluster specified in ~/.kube/config. 
```