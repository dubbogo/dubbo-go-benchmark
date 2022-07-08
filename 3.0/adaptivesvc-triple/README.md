# Adaptive Service Tests

To simplify adapive service tests, we run tests on Kubernetes. You should ensure that the Kubernetes has been installed on your server properly before we go next.

The test were validated on:

- kubernetes v1.22.3
- debian 9 x86_64

The configuration of our server is:

- 8-core CPU
- 16 GB RAM

## Getting Started

### Building Image

> Before going next, please check if the Dockerfile is correct, especially for git branch.

Go to `server` directory.

```shell
cd server
```

Build the server image via docker, in this case we use "xavierniu" as account of Docker Hub.

```shell
docker build . -t xavierniu/dubbogo-server-adasvc
```

Push the image to Docker Hub, before going next, make sure that you have executed `docker login`.

```shell
docker push xavierniu/dubbogo-server-adasvc
```

Then do the samethings for the dubbo client.

```shell
cd client
docker build . -t xavierniu/dubbogo-client-adasvc
docker push xavierniu/dubbogo-client-adasvc
```

### Deployment on K8s

> Please note that we choose two specific images:
> 
> - dubbogo-client: xavierniu/dubbogo-client-adasvc
> - dubbogo-server: xavierniu/dubbogo-server-adasvc

You may change the images and envs as required.

```shell
# Go to `yamls` directory.
cd yamls

# Install namespace.
kubectl apply -f namespace.yml

# Install ZooKeeper.
kubectl apply -f zookeeper.yml

# Create Server Config
kubectl apply -f server_conf.yaml

# Install dubbogo-server.
kubectl apply -f server.yml
# If you want to create servers of different configs, run
kubectl apply -f servers

# Create Client Config
kubectl apply -f client_conf.yaml

# Install dubbogo-client.
kubectl apply -f client.yml
```

Check the deployment.

```shell
# Service
kubectl get svc -n dubbogo-adaptivesvc
# NAME                TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
# zookeeper-service   ClusterIP   x.x.x.x        <none>        2181/TCP   9m43s

# Pod
kubectl get po -n dubbogo-adaptivesvc
# NAME                                            READY   STATUS    RESTARTS   AGE
# dubbogo-server-65cbf5954d-7sjtx                 1/1     Running   0          92m
# dubbogo-server-65cbf5954d-hbclz                 1/1     Running   0          92m
# dubbogo-server-65cbf5954d-vdd5f                 1/1     Running   0          92m
# dubbogo-zookeeper-deployment-587cc5675b-7782j   1/1     Running   0          3d1h

# Deployment
kubectl get deployment -n dubbogo-adaptivesvc
# NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
# dubbogo-server                 3/3     3            3           92m
# dubbogo-zookeeper-deployment   1/1     1            1           3d1h

# Job
kubectl get jobs -n dubbogo-adaptivesvc
```

Undeployment.

```shell
# Delete the dubbogo-client
kubectl delete -f client.yml
```

Update images.

```shell
# Update dubbogo-server
docker pull xavierniu/dubbogo-server-adasvc
kubectl rollout restart deployment dubbogo-server -n dubbogo-adaptivesvc

# Update dubbogo-client
docker pull xavierniu/dubbogo-client-adasvc
```

Update Configs.
> If you want to update dubbogo'config, just update in client_conf.yaml or server_conf.yaml and apply them.
> Besides, you should delete recreate the servers and clients

### Visualize

> If you want to visualize your test results using tools like grafana, you need to deploy prometheus inside the cluster. You can refer to the directory `3.0/deploy/kubernetes` for deployment and configuration.
