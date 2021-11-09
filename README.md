# secure access cloud operator

init using the following: 

```shell
>> kubebuilder init --project-name secure-access-cloud --domain secure-access-cloud.symantec.com --repo bitbucket.org/accezz-io/secure-access-cloud-operator --skip-go-version-check

>> go mod edit -go=1.17

# allowing multigroup
>> kubebuilder edit --multigroup=true 

# create api
>> kubebuilder create api --group access --version v1 --kind Site --resource --controller
```
