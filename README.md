# dummy-operator
A simple dummy controller for managing nginx pods.

## Description
The dummy operator is a custom Kubernetes controller used for creating nginx pods and giving a status field to the Custom Resource that tracks the
status of the Nginx Pod associated with the CR (Dummy).

The dummy operator has three objectives:

1. Logs the Spec of the Custom Resource (`name`, `namespace` and
   the value of `spec.message`).
2. Copies the value of `spec.message` into `status.specEcho` for each CR (Dummy).
3. Creates and associates an `nginx` Pod to each Dummy API object and gives CR (Dummy) a status field that tracks the status of the Pod associated to the Dummy.

## Technologies
The project has been created using these technologies:

* **golang** as programming language
* **Kubernetes** as container-orchestration system
* **kubectl** as command-line tool to interact kubernetes
* **Kubebuilder** as a framework for building Kubernetes APIs using custom resource definitions (CRDs)

## Installation
**Install:**

1. `golang v1.19.0+` from <a href="https://go.dev/dl/">here</a>
2. `Docker 17.03+` from <a href="https://docs.docker.com/get-docker/">here</a> then enable kubernetes
3. `kubectl  v1.11.3+` from <a href="https://kubernetes.io/docs/tasks/tools/">here</a>
4. `Kubebuilder 3.7.0` from <a href="https://book.kubebuilder.io/quick-start.html">here</a>

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
>In this project, controller deployed in a Kubernetes single-node cluster inside Docker Desktop rather than KIND or remote cluster.

### Cloning the Repository
```
git clone https://github.com/mmertdogann/dummy-operator
```

### Running on the cluster

>For testing the controller we can either deploy it on a Kubernetes cluster with `make deploy IMG=<some-registry>/dummy-operator:tag` command  or
we can directly run the controller with `make run IMG=<some-registry>/dummy-operator:tag` command.

### Test the controller by deploying it on the cluster

1. Build and push the controller image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/dummy-operator:tag
```
	
2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/dummy-operator:tag
```

>>> The dummy controller image named `mmertdogann/dummy-operator:0.1` has already built and pushed to the docker hub.

So we can deploy our dummy controller using this command:

```sh
make deploy IMG=mmertdogann/dummy-operator:0.1
```

After this command, the controller will be deployed under `dummy-operator-system`namespace

Controller logs could be observed by this command:

```sh
k logs -f <pod-name> -n dummy-operator-system -c manager
```

3. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
Undeploy the controller to the cluster:

```sh
make undeploy
```

### Test the controller by running
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run the controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

3. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

