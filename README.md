# k8s-crd

## Custom Resource Definition (CRD):

### Define the structure of the custom resource, including the allowed fields, data types, and validations. 

- Run the command to create the CRD.

```sh
kubectl apply -f ./k8s/CustomResourceDefinition.yaml
```

- Check if the CRD has been created.
```sh
kubectl get crd buckets.levi.com
```

## Custom Resource Instances: 

### They are individual objects created based on the definition of a Custom Resource Definition.

- Run the command to create an instance of the CRD.
```sh
kubectl apply -f ./k8s/CustomResourceInstance.yaml
```

- Check if the Instance has been created.

```sh
kubectl delete crd buckets.levi.com
```

## Remove resources.


```sh
kubectl delete crd buckets.levi.com
```

```sh
kubectl delete bucket example-bucket
```