# Secure-Access-Cloud Kubernetes Operator

## Introduction
Create and Govern the access to your k8s resources the same way you create them.   
[Secure-access-cloud](https://www.broadcom.com/products/cyber-security/network/web-protection/secure-access-cloud) custom Kubernetes controller makes it easy to expose Kubernetes services through SAC using k8s [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

Currently supporting 2 CRDs:

1. site
2. HTTP application

## Installing

1. In secure-access-cloud admin portal - create an api client assign and make sure it has admin permission
   ![create-api!](assets/create-api.gif "create admin api")

2. In the desired k8s cluster, create a generic secret that the operator will use to access secure-access-cloud
```shell
>> kubectl create --save-config namespace secure-access-cloud-system
>> kubectl -n secure-access-cloud-system create secret generic secure-access-cloud-config \
--from-literal="tenantDomain=<api endpoint>" \
--from-literal="clientId=<client Id>" \
--from-literal="clientSecret=<client secret>"
```

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