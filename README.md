
## 设置别名
```bash
alias k=kubectl
alias ks='kubectl -n kube-system'
```

## kube-system核心组件
master节点上的四个核心组件：
* api-server
* controller-manager
* scheduler
* etcd

worker节点上的两个核心组件：
* kubelet
* proxy
```
ks get po
```

```
NAME                                READY   STATUS    RESTARTS        AGE
coredns-6d8c4cb4d-gl6g6             1/1     Running   157 (16d ago)   249d
coredns-6d8c4cb4d-q5cwn             1/1     Running   157 (16d ago)   249d
etcd-k8smaster                      1/1     Running   165 (16d ago)   249d
kube-apiserver-k8smaster            1/1     Running   221 (16d ago)   249d
kube-controller-manager-k8smaster   1/1     Running   243 (16d ago)   249d
kube-proxy-g7sbh                    1/1     Running   159 (16d ago)   249d
kube-scheduler-k8smaster            1/1     Running   242 (16d ago)   249d
```


## etcd

### 进入etcd容器内部

```bash
ks exec -it etcd-k8smaster sh
```

```
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
sh-5.1# 
```

#### 在etcd容器中执行命令，查看所有的以/开头的key
```bash
etcdctl --endpoints https://localhost:2379 --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --cacert /etc/kubernetes/pki/etcd/ca.crt get --keys-only --prefix /
```

```
/registry/apiextensions.k8s.io/customresourcedefinitions/apiservers.operator.tigera.io
/registry/apiextensions.k8s.io/customresourcedefinitions/bgpconfigurations.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/bgppeers.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/blockaffinities.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/caliconodestatuses.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/clusterinformations.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/felixconfigurations.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/globalnetworkpolicies.crd.projectcalico.org
...
```

#### 在etcd容器中执行命令

```
etcdctl --endpoints https://localhost:2379 --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --cacert /etc/kubernetes/pki/etcd/ca.crt get --prefix /registry/services/specs/default/kubernetes 
```


输出如下所示，乱码是因为采用的是protobuffer格式
```
k8s

v1Service�
�

kubernetes�default"*$d44517f4-3c45-4989-8581-7a83ee3a15cd2����Z
	component	apiserverZ
provider
kubernetesz��
kube-apiserverUpdate�v����FieldsV1:�
�{"f:metadata":{"f:labels":{".":{},"f:component":{},"f:provider":{}}},"f:spec":{"f:clusterIP":{},"f:internalTrafficPolicy":{},"f:ipFamilyPolicy":{},"f:ports":{".":{},"k:{\"port\":443,\"protocol\":\"TCP\"}":{".":{},"f:name":{},"f:port":{},"f:protocol":{},"f:targetPort":{}}},"f:sessionAffinity":{},"f:type":{}}}Bm
�
httpsTCP��2�(�	10.96.0.1"	ClusterIP:NoneBRZ`h�
                                                    SingleStack�	10.96.0.1�IPv4�Cluster�
�"

```

#### 监听etcd对象变化

```
etcdctl --endpoints https://localhost:2379 --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --cacert /etc/kubernetes/pki/etcd/ca.crt watch --prefix /registry/services/specs/default/mynginx 
```

执行后进入到等待结果状态


## APIServer

### 创建Pod

```
cd ~/go/src/github.com/cncamp/101/module4
cat nginx-deploy.yaml
```

执行结果
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx
```

```
k create -f nginx-deploy.yaml
```

执行结果
```
deployment.apps/nginx-deployment created
```

```
k apply -f nginx-deploy.yaml -v 9
```

```log
I0901 18:15:19.412098  559679 loader.go:372] Config loaded from file:  /root/.kube/config
I0901 18:15:19.413725  559679 round_trippers.go:466] curl -v -XGET  -H "Accept: application/com.github.proto-openapi.spec.v2@v1.0+protobuf" -H "User-Agent: kubectl/v1.23.1 (linux/amd64) kubernetes/86ec240" 'https://192.168.34.2:6443/openapi/v2?timeout=32s'
I0901 18:15:19.460619  559679 round_trippers.go:510] HTTP Trace: Dial to tcp:192.168.34.2:6443 succeed
I0901 18:15:19.467115  559679 round_trippers.go:570] HTTP Statistics: DNSLookup 0 ms Dial 0 ms TLSHandshake 5 ms ServerProcessing 0 ms Duration 53 ms
I0901 18:15:19.467218  559679 round_trippers.go:577] Response Headers:
I0901 18:15:19.467268  559679 round_trippers.go:580]     Last-Modified: Sun, 14 Aug 2022 14:47:33 GMT
I0901 18:15:19.467307  559679 round_trippers.go:580]     X-Varied-Accept: application/com.github.proto-openapi.spec.v2@v1.0+protobuf
I0901 18:15:19.467343  559679 round_trippers.go:580]     Accept-Ranges: bytes
I0901 18:15:19.467379  559679 round_trippers.go:580]     Content-Type: application/octet-stream
I0901 18:15:19.467413  559679 round_trippers.go:580]     Date: Thu, 01 Sep 2022 10:15:19 GMT
I0901 18:15:19.467449  559679 round_trippers.go:580]     Etag: "D334FA1E92FB39BADDBA73798FB5CAF5BFAD240E7A0AFF4715A3A2E92BEB1CE9F2CDEE615128DC4763514BEB2A4F22A7CE006DECD2BC5C3C074FBFCD5FA59D6A"
I0901 18:15:19.467487  559679 round_trippers.go:580]     Audit-Id: 5ce155b4-d348-4a0f-a393-800e8c1ced8e
I0901 18:15:19.467522  559679 round_trippers.go:580]     Cache-Control: no-cache, private
I0901 18:15:19.467557  559679 round_trippers.go:580]     Vary: Accept-Encoding
I0901 18:15:19.467591  559679 round_trippers.go:580]     Vary: Accept
I0901 18:15:19.467627  559679 round_trippers.go:580]     X-From-Cache: 1
I0901 18:15:19.565809  559679 request.go:1179] Response Body:
00000000  0a 03 32 2e 30 12 15 0a  0a 4b 75 62 65 72 6e 65  |..2.0....Kuberne|
00000010  74 65 73 12 07 76 31 2e  32 33 2e 31 42 92 fa b3  |tes..v1.23.1B...|
00000020  01 12 8c 02 0a 22 2f 2e  77 65 6c 6c 2d 6b 6e 6f  |....."/.well-kno|
00000030  77 6e 2f 6f 70 65 6e 69  64 2d 63 6f 6e 66 69 67  |wn/openid-config|
00000040  75 72 61 74 69 6f 6e 2f  12 e5 01 12 e2 01 0a 09  |uration/........|
00000050  57 65 6c 6c 4b 6e 6f 77  6e 1a 57 67 65 74 20 73  |WellKnown.Wget s|
00000060  65 72 76 69 63 65 20 61  63 63 6f 75 6e 74 20 69  |ervice account i|
```


#### get deploy

```
k get deploy
```

```
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   1/1     1            1           17m
```

#### get replica set

```
k get rs
```

```
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-85b98978db   1         1         1       17m
```

```
k get rs nginx-deployment-85b98978db -oyaml
```

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  annotations:
    deployment.kubernetes.io/desired-replicas: "1"
    deployment.kubernetes.io/max-replicas: "2"
    deployment.kubernetes.io/revision: "1"
  creationTimestamp: "2022-09-01T10:09:39Z"
  generation: 1
  labels:
    app: nginx
    pod-template-hash: 85b98978db
  name: nginx-deployment-85b98978db
  namespace: default
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: Deployment
    name: nginx-deployment
    uid: 33e8b8a9-3517-481f-b82e-02837c770cc6
  resourceVersion: "987542"
  uid: 56362d1e-2c63-4892-abf2-74b27ff0b8bd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
      pod-template-hash: 85b98978db
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
        pod-template-hash: 85b98978db
    spec:
      containers:
      - image: nginx
        imagePullPolicy: Always
        name: nginx
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 1
  fullyLabeledReplicas: 1
  observedGeneration: 1
  readyReplicas: 1
  replicas: 1
```

#### get pod

```
k get po
```

```
NAME                                READY   STATUS    RESTARTS        AGE
nginx                               1/1     Running   151 (17d ago)   249d
nginx-deployment-85b98978db-d782x   1/1     Running   0               32m
```

```
k get po nginx-deployment-85b98978db-d782x -oyaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    cni.projectcalico.org/containerID: c65f01006c28cbcd63adcad642904238ef49e083ce742c7a3276733cc95829a0
    cni.projectcalico.org/podIP: 192.168.16.189/32
    cni.projectcalico.org/podIPs: 192.168.16.189/32
  creationTimestamp: "2022-09-01T10:09:39Z"
  generateName: nginx-deployment-85b98978db-
  labels:
    app: nginx
    pod-template-hash: 85b98978db
  name: nginx-deployment-85b98978db-d782x
  namespace: default
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: ReplicaSet
    name: nginx-deployment-85b98978db
    uid: 56362d1e-2c63-4892-abf2-74b27ff0b8bd
  resourceVersion: "987541"
  uid: 4fd3e00d-e63f-43b4-a771-9bcb20142dfd
spec:
  containers:
  - image: nginx
    imagePullPolicy: Always
    name: nginx
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-dfzs2
      readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  nodeName: k8smaster
  preemptionPolicy: PreemptLowerPriority
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext: {}
  serviceAccount: default
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  volumes:
  - name: kube-api-access-dfzs2
    projected:
      defaultMode: 420
      sources:
      - serviceAccountToken:
          expirationSeconds: 3607
          path: token
      - configMap:
          items:
          - key: ca.crt
            path: ca.crt
          name: kube-root-ca.crt
      - downwardAPI:
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
            path: namespace
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "2022-09-01T10:09:39Z"
    status: "True"
    type: Initialized
  - lastProbeTime: null
    lastTransitionTime: "2022-09-01T10:09:57Z"
    status: "True"
    type: Ready
  - lastProbeTime: null
    lastTransitionTime: "2022-09-01T10:09:57Z"
    status: "True"
    type: ContainersReady
  - lastProbeTime: null
    lastTransitionTime: "2022-09-01T10:09:39Z"
    status: "True"
    type: PodScheduled
  containerStatuses:
  - containerID: docker://ede3473765ae1177e02a24ea3ce28a8ac937dadcb05b07f7119c42ec0523e95a
    image: nginx:latest
    imageID: docker-pullable://nginx@sha256:0d17b565c37bcbd895e9d92315a05c1c3c9a29f762b011a10c54a66cd53c9b31
    lastState: {}
    name: nginx
    ready: true
    restartCount: 0
    started: true
    state:
      running:
        startedAt: "2022-09-01T10:09:57Z"
  hostIP: 10.0.2.15
  phase: Running
  podIP: 192.168.16.189
  podIPs:
  - ip: 192.168.16.189
  qosClass: BestEffort
  startTime: "2022-09-01T10:09:39Z"
```


```
k describe po nginx-deployment-85b98978db-d782x
```

```yaml
Name:         nginx-deployment-85b98978db-d782x
Namespace:    default
Priority:     0
Node:         k8smaster/10.0.2.15
Start Time:   Thu, 01 Sep 2022 18:09:39 +0800
Labels:       app=nginx
              pod-template-hash=85b98978db
Annotations:  cni.projectcalico.org/containerID: c65f01006c28cbcd63adcad642904238ef49e083ce742c7a3276733cc95829a0
              cni.projectcalico.org/podIP: 192.168.16.189/32
              cni.projectcalico.org/podIPs: 192.168.16.189/32
Status:       Running
IP:           192.168.16.189
IPs:
  IP:           192.168.16.189
Controlled By:  ReplicaSet/nginx-deployment-85b98978db
Containers:
  nginx:
    Container ID:   docker://ede3473765ae1177e02a24ea3ce28a8ac937dadcb05b07f7119c42ec0523e95a
    Image:          nginx
    Image ID:       docker-pullable://nginx@sha256:0d17b565c37bcbd895e9d92315a05c1c3c9a29f762b011a10c54a66cd53c9b31
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Thu, 01 Sep 2022 18:09:57 +0800
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-dfzs2 (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             True 
  ContainersReady   True 
  PodScheduled      True 
Volumes:
  kube-api-access-dfzs2:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  33m   default-scheduler  Successfully assigned default/nginx-deployment-85b98978db-d782x to k8smaster
  Normal  Pulling    33m   kubelet            Pulling image "nginx"
  Normal  Pulled     33m   kubelet            Successfully pulled image "nginx" in 16.044652467s
  Normal  Created    33m   kubelet            Created container nginx
  Normal  Started    33m   kubelet            Started container nginx
```


## 推荐的Add-ons

* kube-dns: 负责为整个集群提供DNS服务；
* Ingress Congroller: 为服务提供外网接口；
* MetricsServer: 提供资源监控；
* Dashboard: 提供GUI；
* Fluentd-Elasticsearch: 提供集群日志采集、存储与查询；


## kubectl

kubectl的配置文件默认位置在：~/.kube/config

```
k get po
k get po --kubeconfig ~/.kube/config
```

### 查看kubectl配置
```
k config view --minify
```

### wide输出格式
```
k get po -owide
```

```
NAME                                READY   STATUS    RESTARTS        AGE     IP               NODE        NOMINATED NODE   READINESS GATES
nginx                               1/1     Running   151 (17d ago)   249d    192.168.16.139   k8smaster   <none>           <none>
nginx-deployment-85b98978db-d782x   1/1     Running   0               4h23m   192.168.16.189   k8smaster   <none>           <none>
```

### json输出格式
```
k get po -ojson
```

```json
{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "kind": "Pod",
            "metadata": {
                "annotations": {
                    "cni.projectcalico.org/containerID": "5782af629cb726770e3c624c4c742f1fba17becd07df16c3459e870c6a16084e",
                    "cni.projectcalico.org/podIP": "192.168.16.139/32",
                    "cni.projectcalico.org/podIPs": "192.168.16.139/32"
                },
                "creationTimestamp": "2021-12-25T15:31:06Z",
                "labels": {
                    "run": "nginx"
                },
                "name": "nginx",
                "namespace": "default",
                "resourceVersion": "923379",
                "uid": "8d6c729b-aed8-41d7-9dc2-c32692e775f5"
            },
...
```

### yaml输出格式
```
k get po -oyaml
```

```yaml
apiVersion: v1
items:
- apiVersion: v1
  kind: Pod
  metadata:
    annotations:
      cni.projectcalico.org/containerID: 5782af629cb726770e3c624c4c742f1fba17becd07df16c3459e870c6a16084e
      cni.projectcalico.org/podIP: 192.168.16.139/32
      cni.projectcalico.org/podIPs: 192.168.16.139/32
    creationTimestamp: "2021-12-25T15:31:06Z"
    labels:
      run: nginx
    name: nginx
    namespace: default
    resourceVersion: "923379"
    uid: 8d6c729b-aed8-41d7-9dc2-c32692e775f5
  spec:
    containers:
    - image: nginx
      imagePullPolicy: Always
...
```

### 监视变化 -w 即 --watch
```
k get po -w -oyaml
```

```
k get po 
```

```
NAME                                READY   STATUS    RESTARTS        AGE
nginx                               1/1     Running   151 (17d ago)   249d
nginx-deployment-85b98978db-d782x   1/1     Running   0               4h33m
```

```
k edit po nginx-deployment-85b98978db-d782x
```

修改配置文件，增加label a: b
```yaml
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: v1
kind: Pod
metadata:
  annotations:
    cni.projectcalico.org/containerID: c65f01006c28cbcd63adcad642904238ef49e083ce742c7a3276733cc95829a0
    cni.projectcalico.org/podIP: 192.168.16.189/32
    cni.projectcalico.org/podIPs: 192.168.16.189/32
  creationTimestamp: "2022-09-01T10:09:39Z"
  generateName: nginx-deployment-85b98978db-
  labels:
    a: b
    app: nginx
    pod-template-hash: 85b98978db
  name: nginx-deployment-85b98978db-d782x
  namespace: default
```

查看上一条命令监视的内容，发现已经自动更新


### show-labels

```
k get po --show-labels -w
```

```
NAME                                READY   STATUS    RESTARTS        AGE     LABELS
nginx                               1/1     Running   151 (18d ago)   249d    run=nginx
nginx-deployment-85b98978db-d782x   1/1     Running   0               4h39m   a=b,app=nginx,pod-template-hash=85b98978db
nginx-deployment-85b98978db-d782x   1/1     Running   0               4h40m   a=b,app=nginx,pod-template-hash=85b98978db
```

### describe

```
k describe po
```

### exec
```
k exec -it nginx-deployment-85b98978db-d782x bash
```

### logs
```
k logs -f nginx-deployment-85b98978db-d782x
```

```
/docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
/docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
/docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
/docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
/docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
/docker-entrypoint.sh: Configuration complete; ready for start up
2022/09/01 10:09:57 [notice] 1#1: using the "epoll" event method
2022/09/01 10:09:57 [notice] 1#1: nginx/1.21.5
2022/09/01 10:09:57 [notice] 1#1: built by gcc 10.2.1 20210110 (Debian 10.2.1-6) 
2022/09/01 10:09:57 [notice] 1#1: OS: Linux 5.15.0-41-generic
2022/09/01 10:09:57 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
2022/09/01 10:09:57 [notice] 1#1: start worker processes
2022/09/01 10:09:57 [notice] 1#1: start worker process 32
2022/09/01 10:09:57 [notice] 1#1: start worker process 33
2022/09/01 10:09:57 [notice] 1#1: start worker process 34
2022/09/01 10:09:57 [notice] 1#1: start worker process 35
```

### labels

```
k get po --show-labels
```

```
NAME                                READY   STATUS    RESTARTS        AGE    LABELS
nginx                               1/1     Running   151 (19d ago)   251d   run=nginx
nginx-deployment-85b98978db-d782x   1/1     Running   0               46h    a=b,app=nginx,c=d,pod-template-hash=85b98978db
```

```
k get po -l a=b
```

```
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-85b98978db-d782x   1/1     Running   0          46h
```
## 云计算

### 云计算的层次
* Applications
* Data
* Runtime
* Middleware
* OS
* Virtualization
* Servers
* Storage
* Networking

## kubernetes核心对象

## 对象属性定义
#### TypeMeta

* Group
* Kind
* Version

#### MetaData

* namespace
* name
* label
* annotation
* finalizer
* resourceVersion

#### Spec

Spec 是用户的期望状态，由创建对象的用户端定义

### Status

Status 是对象的实际状态，由对应的控制器收集实际状态并更新


### Node

Node是Pod真正运行的主机，可以是物理机或虚拟机
为了管理Pod，每个Node节点上至少要运行container runtime(Docker 或者 Rkt)、kubelet、kube-proxy

```
 k get no -oyaml k8smaster
```

```yaml
apiVersion: v1
kind: Node
metadata:
  annotations:
    kubeadm.alpha.kubernetes.io/cri-socket: /var/run/dockershim.sock
    node.alpha.kubernetes.io/ttl: "0"
    projectcalico.org/IPv4Address: 192.168.20.1/24
    projectcalico.org/IPv4VXLANTunnelAddr: 192.168.16.128
    volumes.kubernetes.io/controller-managed-attach-detach: "true"
  creationTimestamp: "2021-12-25T12:59:23Z"
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/os: linux
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: k8smaster
    kubernetes.io/os: linux
    node-role.kubernetes.io/control-plane: ""
    node-role.kubernetes.io/master: ""
    node.kubernetes.io/exclude-from-external-load-balancers: ""
  name: k8smaster
  resourceVersion: "1003847"
  uid: 0ea06885-ce76-4b26-b956-85e3d3751e16
spec:
  podCIDR: 192.168.0.0/24
  podCIDRs:
  - 192.168.0.0/24
status:
  addresses:
  - address: 10.0.2.15
    type: InternalIP
  - address: k8smaster
    type: Hostname
  allocatable:
    cpu: "4"
    ephemeral-storage: "27895316844"
    hugepages-2Mi: "0"
    memory: 2951220Ki
    pods: "110"
  capacity:
    cpu: "4"
    ephemeral-storage: 30268356Ki
    hugepages-2Mi: "0"
    memory: 3053620Ki
    pods: "110"
  conditions:
  - lastHeartbeatTime: "2022-08-14T14:47:00Z"
    lastTransitionTime: "2022-08-14T14:47:00Z"
    message: Calico is running on this node
    reason: CalicoIsUp
    status: "False"
    type: NetworkUnavailable
  - lastHeartbeatTime: "2022-09-03T09:07:35Z"
    lastTransitionTime: "2021-12-25T12:59:21Z"
    message: kubelet has sufficient memory available
    reason: KubeletHasSufficientMemory
    status: "False"
    type: MemoryPressure
  - lastHeartbeatTime: "2022-09-03T09:07:35Z"
    lastTransitionTime: "2021-12-25T12:59:21Z"
    message: kubelet has no disk pressure
    reason: KubeletHasNoDiskPressure
    status: "False"
    type: DiskPressure
  - lastHeartbeatTime: "2022-09-03T09:07:35Z"
    lastTransitionTime: "2021-12-25T12:59:21Z"
    message: kubelet has sufficient PID available
    reason: KubeletHasSufficientPID
    status: "False"
    type: PIDPressure
  - lastHeartbeatTime: "2022-09-03T09:07:35Z"
    lastTransitionTime: "2022-06-13T11:02:19Z"
    message: kubelet is posting ready status. AppArmor enabled
    reason: KubeletReady
    status: "True"
    type: Ready
  daemonEndpoints:
    kubeletEndpoint:
      Port: 10250
  images:
  - names:
    - python@sha256:cfa62318c459b1fde9e0841c619906d15ada5910d625176e24bf692cf8a2601d
    - python:2.7
    sizeBytes: 901778544
  - names:
    - registry.aliyuncs.com/google_containers/etcd@sha256:64b9ea357325d5db9f8a723dcf503b5a449177b17ac87d69481e126bb724c263
    - registry.aliyuncs.com/google_containers/etcd:3.5.1-0
    sizeBytes: 292558922
  - names:
    - calico/cni@sha256:ce618d26e7976c40958ea92d40666946d5c997cd2f084b6a794916dc9e28061b
    - calico/cni:v3.21.2
    sizeBytes: 238868619
  - names:
    - calico/node@sha256:6912fe45eb85f166de65e2c56937ffb58c935187a84e794fe21e06de6322a4d0
    - calico/node:v3.21.2
    sizeBytes: 213767413
  - names:
    - calico/apiserver@sha256:0e947c69392a6c52bf2bc5f38419f2d8973b89b89cf64df0417c6ff14a8f1bdc
    - calico/apiserver:v3.21.2
    sizeBytes: 192714818
  - names:
    - quay.io/tigera/operator@sha256:b4e3eeccfd3d5a931c07f31c244b272e058ccabd2d8155ccc3ff52ed78855e69
    - quay.io/tigera/operator:v1.23.3
    sizeBytes: 182927858
  - names:
    - nginx@sha256:0d17b565c37bcbd895e9d92315a05c1c3c9a29f762b011a10c54a66cd53c9b31
    - nginx:latest
    sizeBytes: 141479488
  - names:
    - registry.aliyuncs.com/google_containers/kube-apiserver@sha256:f54681a71cce62cbc1b13ebb3dbf1d880f849112789811f98b6aebd2caa2f255
    - registry.aliyuncs.com/google_containers/kube-apiserver:v1.23.1
    sizeBytes: 135162256
  - names:
    - calico/kube-controllers@sha256:1f4fcdcd9d295342775977b574c3124530a4b8adf4782f3603a46272125f01bf
    - calico/kube-controllers:v3.21.2
    sizeBytes: 132282231
  - names:
    - calico/typha@sha256:9e32927a45bcadc8ff3881d1d5f040893d2cbc79c588f60145e3746b20d963b6
    - calico/typha:v3.21.2
    sizeBytes: 128169759
  - names:
    - registry.aliyuncs.com/google_containers/kube-controller-manager@sha256:a7ed87380108a2d811f0d392a3fe87546c85bc366e0d1e024dfa74eb14468604
    - registry.aliyuncs.com/google_containers/kube-controller-manager:v1.23.1
    sizeBytes: 124971684
  - names:
    - redis@sha256:db485f2e245b5b3329fdc7eff4eb00f913e09d8feb9ca720788059fdc2ed8339
    - redis:latest
    sizeBytes: 112691373
  - names:
    - registry.aliyuncs.com/google_containers/kube-proxy@sha256:e40f3a28721588affcf187f3f246d1e078157dabe274003eaa2957a83f7170c8
    - registry.aliyuncs.com/google_containers/kube-proxy:v1.23.1
    sizeBytes: 112327826
  - names:
    - ubuntu@sha256:626ffe58f6e7566e00254b638eb7e0f3b11d4da9675088f4781a50ae288f3322
    - ubuntu:latest
    sizeBytes: 72776513
  - names:
    - registry.aliyuncs.com/google_containers/kube-scheduler@sha256:8be4eb1593cf9ff2d91b44596633b7815a3753696031a1eb4273d1b39427fa8c
    - registry.aliyuncs.com/google_containers/kube-scheduler:v1.23.1
    sizeBytes: 53488305
  - names:
    - registry.aliyuncs.com/google_containers/coredns@sha256:5b6ec0d6de9baaf3e92d0f66cd96a25b9edbce8716f5f15dcd1a616b3abd590e
    - registry.aliyuncs.com/google_containers/coredns:v1.8.6
    sizeBytes: 46829283
  - names:
    - calico/pod2daemon-flexvol@sha256:b034c7c886e697735a5f24e52940d6d19e5f0cb5bf7caafd92ddbc7745cfd01e
    - calico/pod2daemon-flexvol:v3.21.2
    sizeBytes: 21327076
  - names:
    - alpine@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300
    - alpine:latest
    sizeBytes: 5585772
  - names:
    - busybox@sha256:5acba83a746c7608ed544dc1533b87c737a0b0fb730301639a0179f9344b1678
    - busybox:latest
    sizeBytes: 1239820
  - names:
    - registry.aliyuncs.com/google_containers/pause@sha256:3d380ca8864549e74af4b29c10f9cb0956236dfb01c40ca076fb6c37253234db
    - registry.aliyuncs.com/google_containers/pause:3.6
    sizeBytes: 682696
  nodeInfo:
    architecture: amd64
    bootID: 729d1e84-fb88-4aeb-8820-6c5cacf30675
    containerRuntimeVersion: docker://20.10.7
    kernelVersion: 5.15.0-41-generic
    kubeProxyVersion: v1.23.1
    kubeletVersion: v1.23.1
    machineID: e4f33ef3ff614170acdcbd965e104ace
    operatingSystem: linux
    osImage: Ubuntu 20.04.3 LTS
    systemUUID: 6db33528-876c-0747-9a62-ea7ddfd9a6c6
```

### namespace

Namespace是一组资源和对象的抽象集合，常见的pods,services,replication controllers,deployments等都是
属于某一个namesapce的(默认是default)，而node,persistentVolume等则不属于任何namespace

```
k get ns default -oyaml
```

```yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: "2021-12-25T12:59:24Z"
  labels:
    kubernetes.io/metadata.name: default
  name: default
  resourceVersion: "205"
  uid: 784f3a43-5008-4b71-8cff-8ae2077ef99f
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```

### pod

pod是一组紧密关联的容器集合，他们共享PID、IPC、Network和UTS namespace，是kubernetes调度的基本单位

```
k run --image=nginx nginx1
```

```
pod/nginx1 created
```

```
k get po nginx1 -owide
```

```
NAME     READY   STATUS    RESTARTS   AGE    IP               NODE        NOMINATED NODE   READINESS GATES
nginx1   1/1     Running   0          115s   192.168.16.146   k8smaster   <none>           <none>
```

```
k get po nginx1 -oyaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    cni.projectcalico.org/containerID: cb117c0068483a17869aedc78086c677190d9b9a3ba797195651addfd18b7348
    cni.projectcalico.org/podIP: 192.168.16.146/32
    cni.projectcalico.org/podIPs: 192.168.16.146/32
  creationTimestamp: "2022-09-03T09:22:22Z"
  labels:
    run: nginx1
  name: nginx1
  namespace: default
  resourceVersion: "1005837"
  uid: 4392f86c-a9a2-4a43-ba96-43faf9fbcca0
spec:
  containers:
  - image: nginx
    imagePullPolicy: Always
    name: nginx1
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-t5x65
      readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  nodeName: k8smaster
  preemptionPolicy: PreemptLowerPriority
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext: {}
  serviceAccount: default
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  volumes:
  - name: kube-api-access-t5x65
    projected:
      defaultMode: 420
      sources:
      - serviceAccountToken:
          expirationSeconds: 3607
          path: token
      - configMap:
          items:
          - key: ca.crt
            path: ca.crt
          name: kube-root-ca.crt
      - downwardAPI:
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
            path: namespace
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "2022-09-03T09:22:22Z"
    status: "True"
    type: Initialized
  - lastProbeTime: null
    lastTransitionTime: "2022-09-03T09:22:40Z"
    status: "True"
    type: Ready
  - lastProbeTime: null
    lastTransitionTime: "2022-09-03T09:22:40Z"
    status: "True"
    type: ContainersReady
  - lastProbeTime: null
    lastTransitionTime: "2022-09-03T09:22:22Z"
    status: "True"
    type: PodScheduled
  containerStatuses:
  - containerID: docker://76765e87399783dbbe4a899b34db0a098a9d8cba0750acdfeba47521a0afca81
    image: nginx:latest
    imageID: docker-pullable://nginx@sha256:0d17b565c37bcbd895e9d92315a05c1c3c9a29f762b011a10c54a66cd53c9b31
    lastState: {}
    name: nginx1
    ready: true
    restartCount: 0
    started: true
    state:
      running:
        startedAt: "2022-09-03T09:22:40Z"
  hostIP: 10.0.2.15
  phase: Running
  podIP: 192.168.16.146
  podIPs:
  - ip: 192.168.16.146
  qosClass: BestEffort
  startTime: "2022-09-03T09:22:22Z"
```

```
curl 192.168.16.146
```

```html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```


#### 环境变量

* 直接设置值
* 读取Pod Spec的某些属性
* 从ConfigMap读取某个值
* 从Secret读取某个值


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-env
spec:
  containers:
  - image: nginx:1.15
    name: alpine
    env:
    - name: HELLO
      value: world
```


```
k create namespace lianyz
k create -f pod-hello-env.yaml
```

```
pod/hello-env created
```

```
k exec hello-env -it -n lianyz -- sh
```

进入容器后，输入
```
env
```

结果为
```
KUBERNETES_SERVICE_PORT=443
KUBERNETES_PORT=tcp://10.96.0.1:443
HELLO=world
HOSTNAME=hello-env
HOME=/root
PKG_RELEASE=1~bullseye
TERM=xterm
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
NGINX_VERSION=1.21.5
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
KUBERNETES_PORT_443_TCP_PORT=443
NJS_VERSION=0.7.1
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
KUBERNETES_SERVICE_HOST=10.96.0.1
PWD=/
```

#### 存储卷

存储卷包括Volume和VolumeMounts两部分
Volume: 定义Pod可以使用的存储卷来源
VolumeMounts: 定义存储卷如何Mount到容器内部

```
k create -f pod-hello-volume.yaml
```

#### Pod网络

Pod的多个容器共享网络Namespace，同一个Pod中的不通容器可以彼此通过loopback地址访问

#### 资源限制

```
k set resources deployment nginx-app -c=nginx --limits=cpu=500m,memory=128Mi
```

```
deployment.apps/nginx-deployment resource requirements updated
```

#### 健康检查

```
k create -f deploy-centos-readiness.yaml
```

```
deployment.apps/centos created
```

```
k get po
```

```
NAME                               READY   STATUS    RESTARTS        AGE
centos-578b69b65f-jl9ww            0/1     Running   0               79s
hello-volume                       1/1     Running   0               38m
nginx                              1/1     Running   151 (20d ago)   252d
nginx-deployment-667c4d74b-dlnbb   1/1     Running   0               19m
nginx1                             1/1     Running   0               24h
```

```
k exec -it centos-578b69b65f-jl9ww -- bash
```

进入容器后执行
```
touch /tmp/healthy
cat /tmp/healthy
echo $?
```

退出容器后，执行以下命令，READY状态1/1表示共有1个容器，1个容器处于就绪状态

```
k get po
```

```
NAME                               READY   STATUS    RESTARTS        AGE
centos-578b69b65f-jl9ww            1/1     Running   0               2m32s
hello-volume                       1/1     Running   0               39m
nginx                              1/1     Running   151 (20d ago)   252d
nginx-deployment-667c4d74b-dlnbb   1/1     Running   0               21m
nginx1                             1/1     Running   0               24h
```