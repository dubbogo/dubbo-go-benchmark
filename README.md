# Dubbo-go 3.0 benchmark

## 1. Repo Structure

```text
├── 3.0
│   ├── dubbo // dubbo stress test with hessian serialization
│   │   ├── client
│   │   │   ├── client.go
│   │   │   ├── dubbogo.yml
│   │   │   └── user.go
│   │   ├── server
│   │   │   ├── dubbogo.yml
│   │   │   ├── server.go
│   │   │   ├── user.go
│   │   │   └── user_provider.go
│   │   └── stress
│   │       ├── client.go
│   │       ├── dubbogo.yml
│   │       └── user.go
│   └── triple // triple stress test with pb serialization
│       ├── api
│       │   ├── samples_api.pb.go
│       │   ├── samples_api.proto
│       │   └── samples_api_triple.pb.go
│       ├── client
│       │   ├── dubbogo.yml
│       │   └── main.go
│       ├── server
│       │   ├── dubbogo.yml
│       │   └── main.go
│       └── stress
│           ├── dubbogo.yml
│           └── main.go
├── bugreport // bug report
│   └── para200-tps5000-client-error.txt 
```

## 2. Stress Test Typology

stress -> client -> server

## 3. How to run

First, you should have 3 cloud server, under same local network. Build `stress` `client` and `server` binary file, with GOOS=linux GOARCH=amd64 if necessary, and push these three file to each server, with dubbogo.yml.

1. exec server
2. exec client
3. exec stress with Environment:
   - export tps=5000 # your expect ups
   - export payload=0 # your expect payload length for request, you can set to 10k or 50k or any number you like.
   - export parallel=20 # your expect parallel gr number.
   - the stress binary file would print real tps and rt(ns)
4. Collect CPU usage, rt and tps data.  

Pay attention that sometime the parallel is too small to meet your expect tps, for example, with 4c8g machine, you set parallel to 1 and set tps to 5000, the result tps is not as you expected, and you should increase parallel to 10 or other larger number.

## 4. Environments in AdaptiveService Benchmark
1. Random Timeout
   - export RAND_SEED=3             # random seed
   - export TIMEOUT_RATIO=0.05      # your expected timeout ratio
   - export TIMEOUT_DURATION=5s     # your expected timeout duration
2. Random Offline
   - export RAND_SEED=3             # random seed
   - export OFFLINE_RATIO=0.05      # your expected server offline ratio
   - export MIN_ONLINE_DURATION=5s  # your expected minimum online duration
   - export MAX_ONLINE_DURATION=10s # your expected maximum online duration (optional)
   - export MIN_OFFLINE_DURATION=3s # your expected minimum offline duration
   - export MAX_OFFLINE_DURATION=8s # your expected maximum offline duration (optional)
