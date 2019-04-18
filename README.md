# Kubernetes Crash Course

A very short course that aims to uncover basic principles and methodologies for using Kubernetes. The aim is to achieve that with a practical example for setting up services that talk to each other and using [Helm](https://helm.sh/) to setup [Traefik](https://traefik.io/) and [Brigade](https://brigade.sh/) in order to create a simple CD pipelnie.

## Requirements

* [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [helm](https://helm.sh/docs/using_helm/#installing-helm)
* [brig](https://docs.brigade.sh/intro/quickstart/)

## Let's begin

First start a local development cluster with `minikube`.

```console
$ minikube -p crashcourse start
```

In order to have our `Docker` images in our development `Kubernetes` cluster, we can use the `Kubernetes` `Docker` socket. You can get the needed environment variables from `minikube`.

```console
$ eval $(minikube -p crashcourse docker-env)
```

Now let's build our first image. In the same terminal where you've set your environment variables to point to the `crashcourse` `Docker` socket run:

```console
$ docker build sample-apps/api -t api:latest
```

After the build is done, you can verify that our image `api:latest` is present on the node by listing all images that `Docker` has locally on the node.

```console
$ docker images
```

Now we can use this image to start our first `Pod`. You can check if `kubectl` points to our development cluster by checking its `context`.

```console
$ kubectl config current-context
```

`minikube` sets the context when you run the cluster. If for some reason this is not the case, check hot to [Configure Access to Multiple Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)
To start our `Pod`, we need to deploy the manifest which describes how the `Pod` should be configured.

```console
$ kubectl create -f sample-apps/api/pod.yaml
```

Check if the pod is running by listing all `Pods` in the `default` namespace.

```console
$ kubectl get pods
```

To access the `Pod` we need to tell kubernetes how we want to do that. We use `Services` to do that. Let's deploy our first service.

```console
$ kubectl create -f sample-apps/api/service.yaml
```

This is a `Service` of a `NodePort` type. `NodePort` services open a port on every node that is specified by the `spec.ports.nodePort` field. This means that if everything is ok we can access our "API" with

```console
$ curl http://$(minikube -p crashcourse ip):30000
```

Using `NodePort`s is discouraged because we limit the possible nodes where our pod can be deployed. This means that we need to have a different way of accessing our services. We can use an ingress load balancer such as [Traefik](https://traefik.io/) to solve this problem. Traefik is a cloud native reverse proxy and load balancer that can route connection to our cluster based on the http `Host` header.

The easiest way to install `Traefik` is using the [Helm](https://github.com/helm/helm) package manager. 

First, lets initialize the server side component `Tiller` so we can start using Helm.

```console
$ helm init
```

`Tiller` should be runnign in the `kube-system` namespace.

```console
$ kubectl get pods -n kube-system
```

Helm packages can be configured with `yaml` files or passing keys and values on the command line. We have a Traefik configuration file in the `helm` directory.

```console
$ helm install stable/traefik --values helm/traefik-values.yaml --name crashtraefik
```

This deploys a lots of stuff. We can inspect how Traefik is packaged by inspecting its `Chart` by running the following command.

```console
$ helm fetch stable/traefik
```

Now becasue `Traefik` routes our requests based on the http `Host` header, we need to somehow send all our requests with the `Host` that is set in the `traefik-values.yaml` configuration. One way to do this is to add our cluster in `/etc/hosts/`. To do this we need the `IP` address of the cluster

```console
$ minikube -p crashcourse ip
```

I can access my cluster at `192.168.99.106`. Change the `IP` address and edi `/etc/hosts` to include the following line

```
192.168.99.106 traefik.crashcouse
```

Now we can access `Traefik` dashboard in [our broswer](http://traefik.crashcourse/dashboard/)!

Now that we don't depend on a `NodePort` to access our `api`, we can scale it by running multiple instances of it. Pods can't be scaled directly, but `Deployments` can. Lets delete our `Pod` and deploy our `api` with a `Deployment`

```console
$ kubectl delete -f sample-apps/api/pod.yaml 
$ kubectl create -f kubectl delete -f sample-apps/api/deployment.yaml
```

Also, let's deploy the `Ingress` object that instructs `Traefik` how to route to the `api`

```console
$ kubectl create -f sample-apps/api/ingress.yaml 
```

The `Host` of this ingress rule is set to `api.crashcourse` so we need to set this again in our `/etc/hosts` fileas we did for the `Traefik` dashboard.

```console
192.168.99.106 api.crashcouse
```

Also, in our `Ingress` manifest we have set the `spec.rules.paths.path.path` field which enables us to define a path on which our service will be available. Because the path is `hows-the-weather` we can reach our service like this

```console
$ curl http://api.crashcourse/hows-the-weather
```

The `Deployment` for our `api` runs 4 replicas. You can check that by running `kubectl get pods`. Now if access the `api` several times you will notice that you are accessing different replicas of it.

```console
$ curl http://api.crashcourse/hows-the-weather
Here at pod our-api-56d5ff664-qs5zv its fine.
$ curl http://api.crashcourse/hows-the-weather
Here at pod our-api-56d5ff664-tkwhc its fine.
```

One way `Pods` can cummunicate between them is by the internal `DNS` server, which resolves the `Services` created in the cluster. Let's deploy our `backend` service which will fetch information from our `api` service.

```console
$ docker build sample-apps/backend/ -t backend:latest
$ kubectl create -f sample-apps/backend/deployment.yaml
$ kubectl create -f sample-apps/backend/service.yaml
$ kubectl create -f sample-apps/backend/ingress.yaml
```

In the `Deployment` of the `backend` we have defined environment variables that will be available in the container environment. In the array `spec.template.spec.containers.env` we have an environment variable named `API_URL` with value `http://api:8081`. Kubernetes DNS server will resolve to a `Pod` of our `api` service. We are accessing on port `8081` because in the `Service` file of the `api` we have defined `spec.ports.port` to be `8081`.

Now let's introduce [brigadejs](https://brigade.sh/) which will enable us to process events such as GitHub pushes to automatically build and deploy new versions of our applications.

Let's first deploy `brigade` with `helm`

```console
$ helm install brigade/brigade --values helm/brigade-values.yaml --name crashbrigade
```

Now we can create our briage project.

```console
$ brig project create

? Project Name keitaroinc/k8s-crashcourse
? Full repository name github.com/keitaroinc/k8s-crashcourse
? Clone URL (https://github.com/your/repo.git) https://github.com/keitaroinc/k8s-crashcourse.git
? Add secrets? No
Auto-generated a Shared Secret: "qtZhEaHpDmF6gXl7Y9dRdoo4"
? Configure GitHub Access? No
? Configure advanced options Yes
? Custom VCS sidecar (enter 'NONE' for no sidecar) deis/git-sidecar:latest
? Build storage size 
? SecretKeyRef usage No
? Build storage class standard
? Job cache storage class standard
? Worker image registry or DockerHub org 
? Worker image name 
? Custom worker image tag 
? Worker image pull policy IfNotPresent
? Worker command yarn -s start
? Initialize Git submodules No
? Allow host mounts Yes
? Allow privileged jobs Yes
? Image pull secrets 
? Default script ConfigMap name 
? brigade.js file path relative to the repository root brigade/brigade.js
? Upload a default brigade.js script 
? Secret for the Generic Gateway (alphanumeric characters only). Press Enter if you want it to be a? Secret for the Generic Gateway (alphanumeric characters only). Press Enter if you want it to be auto-generated 
Auto-generated Generic Gateway Secret: xd9nG
Project ID: brigade-9b3d43323a1845104be8d58f6a2597933052e14aef294fb22cf7e6
```

We can use `brig` - the `brigadejs` CLI to test our pipeline. Our test pipeline is located in `brigade/brigade.js`. Using `brig` we can send events to `brigadejs` manually.

```console
$ brig run keitaroinc/crashcourse-k8s -f brigade/brigade.js -e push
```

The build phase is successful, but the deployment is not. By checking the logs of the `deploy-sample-apps` pod we can see that the `ServiceAccount` need privileges to do the scaling.

`ServiceAccounts` are binded to a `Role` with a `RoleBinding`.

Let's edit the `Role` object directly via `kubectl`. This approach should be avoided, but comes in handy on rare occasions.

```
# Set the `EDITOR` environment variable to your editor of choice
$ kubectl edit role crashbrigade-brigade-wrk
```

Append the following at the end of the `crash-brigade-wrk` role, in the `rules` array.

```
- apiGroups:
  - apps
  resources:
  - deployments/scale
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
```

Now if we run the pipeline again, new images should be built and everything should be redeployed.

