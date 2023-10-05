### Steps to start building operator

### Perquisite
- GO installed
- Kubebuilder installed
`make sure to install samem version of kubebuilder and GoLang`

#### Create a directory, for our example we are creating k8sOperator.
 - mkdir k8sOperator

#### Initialise Go
 - go mod init k8sOperator

#### Initialise kubebuilder
 - kubebuilder init --domain sagar.com  `if you are creating the controller in the directory other then ~/go/src then please define repo path to you current working directory --repo github.com/sagar0419/k8sOperator`
-kubebuilder create api --group k8soperator  --version v1 --kind Demo
This command generates Kubernetes manifests (YAML files) for custom resources (CRDs), RBAC roles, and other resources defined in your operator. These generated manifests are usually stored in the config/crds/ and config/rbac/ directories.
 - make manifests

This command generates Go code based on the custom resource definitions (CRDs) in your project. It generates client code, informers, listers, and other code needed for interacting with your custom resources.
 - make generate

Now make changes in the code as all the files  are initialise. Once that is done follow these commands :-
#### To run operator
- make install run

#### Now deploy the sample app 
- cd config/samples
- k apply -f k8soperator_v1_demo.yaml

`Now you can describe the deployed app`

#### To make container and store it on your dockerHub / private repository

- make docker-build docker-push

#### To deploy Operator as container
- make deploy
`Now you can see that all resources, namespaces, rbac, role-binding, service,deployment has been created in a namespace `k8soperator-system``

#### To Undeploy the operator
- make undeploy

`It will remove everything from the cluster`


### Below is the default readme created by kubebuilder. 
# k8soperator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/k8soperator:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/k8soperator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

