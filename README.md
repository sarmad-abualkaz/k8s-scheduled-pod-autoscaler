# k8s-scheduled-pod-autoscaler

This is a basic [Kubernetes Controller](https://kubernetes.io/docs/concepts/architecture/controller/) to scale pods based on a schedule using the Scheduled Pod Autoscaler (or SPA) [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/). 

As of now, this controller is meant to manage both [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) and [Horizontal Pod Autoscaler (HPAs)](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) Resources.

## How does it work?

Based on the custom resource manifest - which describes the name and type of resource (e.g. `Deployment` or `HPA`), the controller will infer which Kubernetes resource to to update, then perform scheduled scale action on the required time. For Deplyoments this translates to updating the replica under the Deployment.Spec to the value provided in the SPA scale action. For HPAs this translates to updating the minReplicas value in order to ensure at minimum the number of replicas matches the SPA scale action value, upon which the HPA will manage any remaining scale action required from there on. 

*Note: The controller uses server local timezone as a reference when comparing the times provided in th SPA resource.*

### Custom Resource - SPA:
Below is an example design of the Custom Resource managed by this controller to manage a Depoloymen with the name `test-deployment`:

```
apiVersion: autoscaling.spa.sarmadabualkaz.io/v1
kind: ScheduledPodAutoscaler
metadata:
  name: scheduledpodautoscaler-sample
spec:
  resource:
    type: Deployment
    name: test-deployment
  scaleUp:
    time: 8:15AM
    value: 20
  scaleDown:
    time: 10:00PM
    value: 5
```

Based on the above resource, from `spec.scaleUp`/`spec.scaleDown` values under `time` and `value` the controller will check the current time (local timezone) - if it's 8:15 AM, but before 10:00Pm the replicas for `deployment/test-deploymen` will have to be 20 pods. If it's past 10:00PM then it will scale down the replicas on `deployment/test-deploymen` to 5 pods. 

Simlarly for HPAs the SPA resource will look as below: 
```
apiVersion: autoscaling.spa.sarmadabualkaz.io/v1
kind: ScheduledPodAutoscaler
metadata:
  name: scheduledpodautoscaler-hpa-sample
spec:
  resource:
    type: HPA
    name: test-hpa
  scaleUp:
    time: 8:15AM
    value: 20
  scaleDown:
    time: 10:00PM
    value: 5
```
Based on the above resource, from `spec.scaleUp`/`spec.scaleDown` values under `time` and `value` the controller will check the current time (local timezone) - if it's 8:15 AM, but before 10:00Pm the minReplicas for `hpa/test-hpa` will have to be 20 pods. If it's past 10:00PM then it will scale down the minReplicas on `hpa/test-hpa` to 5. 

**Note1: if the maxReplicas on during scaleUp or scaleDown happens to be below the value the spa-controller is expected to update minReplicas to, both maxReplicas and minReplicas will be updated to the required 'new' value.*

***Note2: the controller accepts the following under `spec.resource.type`: `deployment`, `Deployment` for deployments; and `HPA`, `hpa`, `HorizontalPodAutoscaler` `horizontalPodAutoscaler` for HPAs*


## How to install on cluster?

### prerequisite:
Installing the SPA-controller on cluster with its validating and defautling webhooks require [cert-manager](https://github.com/jetstack/cert-manager) to be installed. Please follow [these steps to install cert-manager](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html) prior.

### Direcrtly from project:

1. clone this project: `git clone git@github.com:sarmad-abualkaz/k8s-scheduled-pod-autoscaler.git`

2. build and publish a new image from this projec: `make docker-build docker-push IMG=<some-registry>/k8s-scheduled-pod-autoscaler:<tag>`

    where `some-registry` is the docker registry you will publish the docker image to; and `tag` is the image tag you will publish the image with.

3. install the CRD on your cluster: `make install`

4. deploy the controller and its required resources: `make deploy IMG=<some-registry>/k8s-scheduled-pod-autoscaler:<tag>`

    again here, `some-registry` is the docker registry you published the docker image to; and `tag` is the image tag you published the image with.


### Install the controller using Helm:
  TBD


### Install and test the project locally:

1. clone this project: `git clone git@github.com:sarmad-abualkaz/k8s-scheduled-pod-autoscaler.git`

2. install the CRD on cluster: `make install`

3. run the controller from your machine: `make run ENABLE_WEBHOOKS=false` 

    the command above will run the controller locally (using your local timezone) and disable the webhook components
