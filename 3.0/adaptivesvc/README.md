# Adaptive Service Tests

To simplify adapive service tests, we run tests on Kubernetes. You should ensure that the Kubernetes has been installed on your server properly before we go next.

We had tested tests on:

- kubernetes v1.22.3
- debian 9 x86_64

The configuration of server is:

- 8-core CPU
- 16 GB RAM

## Getting Started

Go to `yamls` directory.

```bash
cd yamls
```

Install namespace.

```bash
kubectl apply -f namespace.yml
```

Install zookeeper.

```bash
kubectl apply -f zookeeper.yml
```

Check zookeeper service.

```bash
kubectl get svc -n dubbogo-adaptivesvc

NAME                TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
zookeeper-service   ClusterIP   x.x.x.x        <none>        2181/TCP   9m43s
```

Check zookeeper pods.

```bash
kubectl get po -n dubbogo-adaptivesvc

NAME                                            READY   STATUS    RESTARTS   AGE
dubbogo-zookeeper-deployment-587cc5675b-7782j   1/1     Running   0          10m
```