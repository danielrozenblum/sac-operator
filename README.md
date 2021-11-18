[![CircleCI](https://circleci.com/bb/accezz-io/sac-operator.svg?style=svg)](https://circleci.com/bb/accezz-io/sac-operator)

# Secure-Access-Cloud Kubernetes Operator

## Introduction
The purpose of the SAC custom Kubernetes controller is to expose Kubernetes service through SAC using K8s CRD in order to be accessible securely without 
the need to expose the service through load-balancer and implement all security controls (authentication, ddos protection, firewall, etc.)

This operator based on Kubebuilder https://book.kubebuilder.io/introduction.html.

## Building
This repository using *go modules* for dependency management and using go 1.17.

1. download dependencies:`$ go get -d ./...`
2. build: `$ make build`

After changing any CRD implementation, you will need to generate the CRD templates using `$ make manifests` and then `$ make generate`

## Local Development
1. Configure your Kubernetes cluster. The simple way to configure it is by installing Docker Desktop and enable Kubernetes in the Preferences page.
2. Install CRDs into the K8s cluster: `$ make install`
3. Validate your resources installed by performing: `$ kubectl api-resources`
4. Run your operator locally: `$ make run ENABLE_WEBHOOKS=false` or deploy it on the K8s cluster: `$ make deploy`
5. Run one of the samples: `$ kubectl create -f config/samples/access_v1_site.yaml`
   At this point, you should see activity in the log file, you can then access your site using `$ kubectl get <resource-name>`

## Debugging
Debug *main.go*

## Running Tests
### Mocks
* Install mockgen tool: `go install github.com/golang/mock/mockgen@v1.6.0`
* Add the following annotation to your interface: `//go:generate mockery -name ApplicationService -inpkg -case=underscore -output MockApplicationService`
* Run `make mocks`

### Tests
* Unit-Tests: `$ make test`
* Integration-Tests: TBD
* Running All: TBD

## Configure Log-Level (TBD)
Configure environment-variable `$ export LOG_LEVEL=debug`

## Docker
Build docker image with the manager: `$ make docker-build IMG=<some-registry>/<project-name>:tag`

Push docker image with the manager: `$ make docker-push IMG=<some-registry>/<project-name>:tag`


## Environment Variables
|Environment-Variable               | Required  | Default Value | Description                                          |
|-----------------------------------|---------- |---------------|------------------------------------------------------|
|LOG_LEVEL                          | no        | info          | The log-level                                        |
|LOG_FORMAT                         | no        | text          | The logger format                                    |
|LOG_FILE                           | no        | n/a           | The logger filename                                  |

## Kubebuilder Properties
This project has been initialized with the following commands:
```shell
>> kubebuilder init --project-name secure-access-cloud --domain secure-access-cloud.symantec.com --repo bitbucket.org/accezz-io/sac-operator --skip-go-version-check
>> go mod edit -go=1.17
>> kubebuilder edit --multigroup=true 
# create api
>> kubebuilder create api --group access --version v1 --kind Site --resource --controller
>> kubebuilder create api --group access --version v1 --kind Application --resource --controller
```

## Internal Endpoints
N/A

## Troubleshooting
N/A
