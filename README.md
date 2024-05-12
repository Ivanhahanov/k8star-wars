# Welcome to k8StarWars
## Preparation for к k8StarsWars
### Install kind

kind create cluster --config universe-config.yml

### create ingress
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

wait for ingress ready
```bash
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

## First clone

Example of clone config 

```yaml
apiVersion: v1 
kind: Pod
metadata:
  name: clone
spec:
  containers:
  - name: clone
    image: explabs/k8star-wars-clone
    ports:
    - containerPort: 80
```

Apply new clone
```bash
kubectl apply -f supply-clones/clone.yml 
```
Command returns
```
pod/clone created
```

Show clone
```bash
kubectl get pods
```
Result is
```
NAME    READY   STATUS    RESTARTS   AGE
clone   1/1     Running   0          37s
```

## Squad of clones
Squad example
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: clone
spec:
  replicas: 3
  selector:
    matchLabels:
      person: clone
      class: empire
  template:
    metadata:
      labels:
        person: clone
        class: empire
    spec:
      containers:
      - name: clone
        image: explabs/k8star-wars-clone
        ports:
        - containerPort: 80
```

Apply new clones squad
```bash
kubectl apply -f supply-clones/squad.yml 
```
Command returns
```
replicaset.apps/clone created
```
Show squad
```
kubectl get rs
```
Output is
```
NAME    DESIRED   CURRENT   READY   AGE
clone   3         3         3       58s
```

Show clones
```
kubectl get pods
```

```
NAME          READY   STATUS    RESTARTS   AGE
clone         1/1     Running   0          2m19s
clone-tvxfb   1/1     Running   0          1m9s
clone-tzng5   1/1     Running   0          1m9s
clone-z4vpv   1/1     Running   0          1m9s
```

## Send a LAAT ship
LAAT(Low Altitude Assault Transport)
Land clones with new LAAT ship
```bash
kubectl apply -f supply-clones/squad.yml 
```
Command returns
```
deployment.apps/clone configured
```

Show Laat
```
kubectl get deploy 
```
Output is
```
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
clone   3/3     3            3           10s
```


Show clones
```
kubectl get pods
```


```
NAME                    READY   STATUS    RESTARTS   AGE
clone                   1/1     Running   0          14m
clone-bbd859cbf-9dlf6   1/1     Running   0          35s
clone-bbd859cbf-nfzzp   1/1     Running   0          35s
clone-bbd859cbf-rkk69   1/1     Running   0          35s
clone-tvxfb             1/1     Running   0          14m
clone-tzng5             1/1     Running   0          14m
clone-z4vpv             1/1     Running   0          14m
```