# Sample Deployment Controller

This is a sample controller. Which watches for deployment resource on kubernetes cluster. When a new
deployment resource is created the controller created a service for the deployment automatically. Also when
the deployment resource is deleted, then the corresponding service also be deleted.

### First Create a namespace called `nginx`. 
Because our deployment will create on nginx namespace. 
```shell
kubectl create ns nginx
```

### Run the controller Code (Terminal 1)
```shell
go run .
```


### Watch all resource in the nginx namespace (Terminal 2)
```shell
watch kubectl get all -n nginx
```

### Create a deployment (Terminal 3)
```shell
kubectl apply -f nginx-deployment.yaml
```

By this time you will see in terminal 2 that a new deployment,service is created.

### Delete the deployment (Terminal 3)
```shell
kubectl delete -f nginx-deployment.yaml
```

### Delete `nginx` Namespace
```shell
kubectl delete ns nginx
```