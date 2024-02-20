# Function brief description
1. According to the example.csv file, you can make replicas and resource-related adjustments to a container resource of the deploy/sts type.
2. If the current resource is less than the changed resource, it will be marked green; if the current resource is greater than the changed resource, it will be marked red.
3. Note that to obtain service quality by default, the app label needs to be standardized, that is, app=workload name.

![image](https://github.com/Einic/cops/blob/main/img/run.png)

# Batch resource changes
1. Only need to be sorted into example.csv, as shown below.

```
workload,containers_name,worktype,namespace,replicas,limits_cpu,limits_memory,requests_cpu,requests_memory
hotrod,hotrod,deployment,sample-application,1,100m,256Mi,100m,256Mi
locust-master,locust-master,deployment,sample-application,2,400m,512Mi,400m,512Mi
locust-worker,locust-worker,deployment,sample-application,2,300m,215Mi,300m,215Mi
```

2. Execute alter resource changes.

```
[root@tcs-192-168-200-132 cops]# ./bin/cops -a /root/.kube/config ./example.csv
┏━━━┳━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━┳━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━━━┓
┃   ┃ DATATIME            ┃ WORKLOAD      ┃ CONTAINERNAME ┃ WORKTYPE ┃ NAMESPACE          ┃ REPLICAS ┃ REQUESTS (CPU) ┃ REQUESTS (MEMORY) ┃ LIMITS (CPU) ┃ LIMITS (MEMORY) ┃ PODQOS     ┃ RUNSTATUS ┃ ALTERSTATUS ┃
┣━━━╋━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━┫
┃ 1 ┃ 2024-02-19 16:48:50 ┃ hotrod        ┃ hotrod        ┃ deploy   ┃ sample-application ┃ 1 -> 2   ┃ 100m -> 100m   ┃ 256Mi -> 216Mi    ┃ 100m -> 100m ┃ 256Mi -> 216Mi  ┃ Guaranteed ┃ Available ┃ Success     ┃
┣━━━╋━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━┫
┃ 2 ┃ 2024-02-19 16:49:10 ┃ locust-master ┃ locust-master ┃ deploy   ┃ sample-application ┃ 2 -> 1   ┃ 400m -> 300m   ┃ 512Mi -> 512Mi    ┃ 400m -> 300m ┃ 512Mi -> 512Mi  ┃ Guaranteed ┃ Available ┃ Success     ┃
┣━━━╋━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━┫
┃ 3 ┃ 2024-02-19 16:49:10 ┃ locust-worker ┃ locust-worker ┃ deploy   ┃ sample-application ┃ 2 -> 2   ┃ 300m -> 300m   ┃ 215Mi -> 512Mi    ┃ 300m -> 300m ┃ 215Mi -> 512Mi  ┃ Guaranteed ┃ Available ┃ Success     ┃
┗━━━┻━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━┻━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━━━━━━┛

```
