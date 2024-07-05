# k8s-crd

## Custom Resource Definition (CRD)

### Define the structure of the custom resource, including the allowed fields, data types, and validations

1- Run the command to create the CRD.

```sh
kubectl apply -f ./k8s/CustomResourceDefinition.yaml
```

2- Check if the CRD has been created.

```sh
kubectl get crd bucket.levi.com
```

## Run localstack

3- Run command

```sh
docker-compose up -d
```

Access the [LocalStack URL](https://app.localstack.cloud/inst/default/resources/s3) to view the list of buckets.

## Run the application

4- Start application

```sh
PATH_KUBE_CONFIG="$HOME/.kube/config" \
AWS_PROFILE_NAME="local" \
AWS_CONFIG_ENDPOINT="http://localhost:4566" \
AWS_CONFIG_ACCESS_KEY_ID="localstack-key-id" \
AWS_CONFIG_SECRET_ACCESS_KEY="localstack-access-key" \
go run .
```

## Custom Resource Instances

### They are individual objects created based on the definition of a Custom Resource Definition

5- Run the command to create an instance of the CRD.

```sh
kubectl apply -f ./k8s/CustomResourceInstance.yaml
```

6- Check if the Instance has been created.

```sh
kubectl get bucket
```

Access [LocalStack URL](https://app.localstack.cloud/inst/default/resources/s3) and refresh the page to view the created bucket.

## Remove resources

7- Remove CRD Instance

```sh
kubectl delete -f ./k8s/CustomResourceInstance.yaml 
```

> Check if the bucket has been removed in  [LocalStack URL](https://app.localstack.cloud/inst/default/resources/s3).

8- Remove CRD

```sh
kubectl delete -f ./k8s/CustomResourceDefinition.yaml
```
