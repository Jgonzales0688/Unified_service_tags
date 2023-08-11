# Admission Controller Project

## Introduction

For my summer project, I developed a Mutating Admission Controller implemented as a webhook. This project allowed me to gain valuable experience beyond my coursework and delve into the infrastructure side of software engineering.

## Motivation

I chose to work on this project because of the company's commitment to making a positive impact on the world. The company's values aligned with what I sought in a workplace, and I was excited to contribute to their mission.

## About the Admission Controller

An admission controller acts as a guard at the door of a Kubernetes cluster, ensuring that resources meet specific criteria before entering. Kubernetes orchestrates containerized applications within clusters, which are collections of nodes (servers) working together.

I developed this admission controller to enhance the capabilities of a monitoring tool. 
It assists in:
- Identifying deployment impact with trace and container metrics filtered by version.
- Seamlessly navigating traces, metrics, and logs using consistent tags.
- Viewing service data based on environment or version uniformly.

## How It Works

I accomplished this by automating metadata addition through Kubernetes tags:
- First tag: Specifies the environment of the workload deployment.
- Second tag: Identifies the service name, showing the name of the deployed workload.
- Third tag: Indicates the version of the deployed workload.

## Non-Tech Example

Think of an admission controller as a bouncer at a party. Just like bouncers check wristbands, the admission controller validates resources using tags. If a resource lacks the necessary tags, it goes through a mutation phase to have them added.

## Technical Implementation

I created a mutating admission controller from scratch in Go Lang. Prior to this project, I had no experience with Go Lang, Kubernetes, or Docker, making this a significant learning opportunity. I'm proud to have earned the title of an expert on admission controllers within our team.

## Process Flow

- The admission controller receives an API request when a resource is created.
- The resource undergoes validation and is checked for required tags.
- If tags are missing, the resource goes through the mutation phase to add them.
- The resource is then subjected to object schema validation, ensuring proper structure.
- Resources with valid tags are moved to Kubernetes data storage (ETCD).

## Installation
This project can fully run locally and includes automation to deploy a local Kubernetes cluster (using Kind).

### Requirements
* Docker
* kubectl
* [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* Go >=1.16 (optional)

## Usage
### Create Cluster
First, we need to create a Kubernetes cluster:
```
❯ make cluster

🔧 Creating Kubernetes cluster...
kind create cluster --config dev/manifests/kind/kind.cluster.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.21.1) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

Make sure that the Kubernetes node is ready:
```
❯ kubectl get nodes
NAME                 STATUS   ROLES                  AGE     VERSION
kind-control-plane   Ready    control-plane,master   3m25s   v1.21.1
```

And that system pods are running happily:
```
❯ kubectl -n kube-system get pods
NAME                                         READY   STATUS    RESTARTS   AGE
coredns-558bd4d5db-thwvj                     1/1     Running   0          3m39s
coredns-558bd4d5db-w85ks                     1/1     Running   0          3m39s
etcd-kind-control-plane                      1/1     Running   0          3m56s
kindnet-84slq                                1/1     Running   0          3m40s
kube-apiserver-kind-control-plane            1/1     Running   0          3m54s
kube-controller-manager-kind-control-plane   1/1     Running   0          3m56s
kube-proxy-4h6sj                             1/1     Running   0          3m40s
kube-scheduler-kind-control-plane            1/1     Running   0          3m54s
```

### Generate Self-Signed Certificates

The certificates the pre-populated Kubernetes secret defined under `/dev/manifests/webhook/webhook.tls.secret.yaml` are already expired. To fix this, you will need to run the script to re-generate them. Simply go to the `/dev/` ( `cd dev/`) directory and from there, run the shell script which generates the secrets with the correct service name already defined. This will update the `webhook.tls.secret.yaml` file with the encoded version of the newly and valid certificates.

You will notice that the script will output the caBundle. Copy everything after the `>> MutatingWebhookConfiguration caBundle:` to the `==` (include the doble equal signs). Paste this into the caBundle field for the `mutating.config.yaml` file under `dev/manifests/cluster-config/mutating.config.yaml`. Make sure to indent it correctly (All the way to the letter 'B' in "caBundle"). Do the same thing for the file `validating.config.yaml`
Once this is done, go back one directory `cd ..`.

### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```
❯ make deploy

📦 Building simple-kubernetes-webhook Docker image...
docker build -t simple-kubernetes-webhook:latest .
[+] Building 14.3s (13/13) FINISHED
...

📦 Pushing admission-webhook image into Kind's Docker daemon...
kind load docker-image simple-kubernetes-webhook:latest
Image: "simple-kubernetes-webhook:latest" with ID "sha256:46b8603bcc11a8fa1825190d3ed99c099096395b22a709e13ec6e7ae2f54014d" not yet present on node "kind-control-plane", loading...

⚙️  Applying cluster config...
kubectl apply -f dev/manifests/cluster-config/
namespace/apps created
mutatingwebhookconfiguration.admissionregistration.k8s.io/simple-kubernetes-webhook.acme.com created
validatingwebhookconfiguration.admissionregistration.k8s.io/simple-kubernetes-webhook.acme.com created

🚀 Deploying simple-kubernetes-webhook...
kubectl apply -f dev/manifests/webhook/
deployment.apps/simple-kubernetes-webhook created
service/simple-kubernetes-webhook created
secret/simple-kubernetes-webhook-tls created
```

Then, make sure the admission webhook pod is running (in the `default` namespace):
```
❯ kubectl get pods
NAME                                        READY   STATUS    RESTARTS   AGE
simple-kubernetes-webhook-77444566b7-wzwmx   1/1     Running   0          2m21s
```

You can stream logs from it:
```
❯ make logs

🔍 Streaming simple-kubernetes-webhook logs...
kubectl logs -l app=simple-kubernetes-webhook -f
time="2021-09-03T04:59:10Z" level=info msg="Listening on port 443..."
time="2021-09-03T05:02:21Z" level=debug msg=healthy uri=/health
```

And hit it's health endpoint from your local machine:
```
❯ curl -k https://localhost:8443/health
OK
```

### Deploying pods
Deploy a valid test pod that gets succesfully created:
```
❯ make pod

🚀 Deploying test pod...
kubectl apply -f dev/manifests/pods/good-pod.yaml
pod/good-pod created
```
You should see in the admission webhook logs that the pod got mutated and validated.

Deploy a non valid pod that gets rejected:
```
❯ make bad-pod

🚀 Deploying "bad" pod...
kubectl apply -f dev/manifests/pods/bad-pod.yaml
Error from server: error when creating "dev/manifests/pods/bad-pod.yaml": admission webhook "simple-kubernetes-webhook.acme.com" denied the request: pod name contains "offensive"
```
You should see in the admission webhook logs that the pod validation failed. It's possible you will also see that the pod was mutated, as webhook configurations are not ordered.

## Testing
Unit tests can be run with the following command:
```
$ make test
go test ./...
?   	github.com/slackhq/simple-kubernetes-webhook	[no test files]
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/admission	0.611s
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/mutation	1.064s
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/validation	0.749s
```

## Admission Logic
A set of validations and mutations are implemented in an extensible framework. Those happen on the fly when a pod is deployed.

#### How to add a new pod validation
To add a new pod mutation, create a file `pkg/validation/MUTATION_NAME.go`, then create a new struct implementing the `validation.podValidator` interface.

#### How to add a new pod mutation
To add a new pod mutation, create a file `pkg/mutation/MUTATION_NAME.go`, then create a new struct implementing the `mutation.podMutator` interface.
