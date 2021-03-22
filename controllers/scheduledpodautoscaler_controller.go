/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	kautoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	autoscalingv1 "spa.sarmadabualkaz.io/spa/api/v1"
)

// ScheduledPodAutoscalerReconciler reconciles a ScheduledPodAutoscaler object
type ScheduledPodAutoscalerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autoscaling.spa.sarmadabualkaz.io,resources=scheduledpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.spa.sarmadabualkaz.io,resources=scheduledpodautoscalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;update;patch

var (
	scheduledTimeAnnotation = "spa.sarmadabualkaz.io/scheduled-at"
)

func (r *ScheduledPodAutoscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("scheduledpodautoscaler", req.NamespacedName)

	// 1. Load the named SPA (ScheduledPodAutoscaler):
	var scheduledPodAutoscaler autoscalingv1.ScheduledPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &scheduledPodAutoscaler); err != nil {
		log.Error(err, "unable to fetch ScheduledPodAutoscaler")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Validate resource is one of the 3 main types 'scream back if its not :|':
	var passedResourceType string
	var passedResourceName string
	var resourceType string

	passedResourceType = scheduledPodAutoscaler.Spec.Resource.Type
	passedResourceName = scheduledPodAutoscaler.Spec.Resource.Name

	if (passedResourceType == "deployment") || (passedResourceType == "Deployment") {
		resourceType = "deployment"
	} else if (passedResourceType == "annotatedDeployment") || (passedResourceType == "AnnotatedDeployment") {
		resourceType = "hpaOperator"
	} else if (passedResourceType == "HPA") || (passedResourceType == "hpa") || (passedResourceType == "HorizontalPodAutoscaler") || (passedResourceType == "horizontalPodAutoscaler") {
		resourceType = "hpa"
	} else {
		err := fmt.Errorf("unrecognizable resource.type %s ResourceType", passedResourceType)
		return ctrl.Result{}, err
	}

	// 3. Get the respective resource (Deployment if resourceType = "deployment" or "hpaOperator"; HorizontalPodAutoscaler if resourceType = "hpaOperator"):
	getResourceSpec := func(resourceName string, resourceType string) (deploymentSpec *appsv1.Deployment, hpaSpec *kautoscalingv1.HorizontalPodAutoscaler, err error) {
		if (resourceType == "deployment") || (resourceType == "hpaOperator") {
			deploymentSpec = &appsv1.Deployment{}

			u := &unstructured.Unstructured{}
			u.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "apps",
				Kind:    "Deployment",
				Version: "v1",
			})

			err := r.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: req.NamespacedName.Namespace}, u)
			if err != nil {
				return &appsv1.Deployment{}, &kautoscalingv1.HorizontalPodAutoscaler{}, err
			}

			runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &deploymentSpec)

			return deploymentSpec, &kautoscalingv1.HorizontalPodAutoscaler{}, nil
		} else if resourceType == "hpa" {
			hpaSpec = &kautoscalingv1.HorizontalPodAutoscaler{}

			u := &unstructured.Unstructured{}
			u.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "autoscaling",
				Kind:    "HorizontalPodAutoscaler",
				Version: "v1",
			})

			err := r.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: req.NamespacedName.Namespace}, u)
			if err != nil {
				return &appsv1.Deployment{}, &kautoscalingv1.HorizontalPodAutoscaler{}, err
			}

			runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &hpaSpec)

			return &appsv1.Deployment{}, hpaSpec, err
		} else {
			err := fmt.Errorf("unrecognizable resource.type %s ResourceType", resourceType)
			return &appsv1.Deployment{}, &kautoscalingv1.HorizontalPodAutoscaler{}, err
		}
	}

	log.V(1).Info("Checking for resource:", "type", resourceType, "name", passedResourceName)

	deploymentSpec, hpaSpec, err := getResourceSpec(passedResourceName, resourceType)

	if err != nil {
		log.Error(err, "unable to find resource for", "resourceName", passedResourceName, "and resource type", resourceType)
		return ctrl.Result{}, err
	}

	// 4. (optional) - Check if we’re suspended (and don’t do anything else if we are)

	// not implemented atm

	// 5. Validate scaling spec do not conflict:
	// retrieve scaleUp spec:
	var scaleUpTimeStr string
	var scaleDownTimeStr string
	var scaleUpValue *int32
	var scaleDownValue *int32
	var requiredReplicas *int32

	scaleUpTimeStr = scheduledPodAutoscaler.Spec.ScaleUp.Time
	scaleUpValue = scheduledPodAutoscaler.Spec.ScaleUp.Value

	// retrieve scaleDown spec:
	scaleDownTimeStr = scheduledPodAutoscaler.Spec.ScaleDown.Time
	scaleDownValue = scheduledPodAutoscaler.Spec.ScaleDown.Value

	if scaleUpTimeStr == scaleDownTimeStr {
		// error-out in case scaling times are the same:
		err := fmt.Errorf("scaleUp.Time is equal scaleDown.Time with %s value", scaleUpTimeStr)
		return ctrl.Result{}, err
	}

	// 6. Confirm which happens earlier - scaleup or scaledown:
	curr_time := time.Now()
	cur_year, cur_month, cur_day := (time.Now()).Date()

	scaleUpTimeZero, _ := time.Parse(time.Kitchen, scaleUpTimeStr)
	scaleDownTimeZero, _ := time.Parse(time.Kitchen, scaleDownTimeStr)

	scaleUpHour, scaleUpMin, scaleUpSeconds := scaleUpTimeZero.Clock()
	scaleDownHour, scaleDownMin, scaleDownSeconds := scaleDownTimeZero.Clock()

	ScaleUpTime := time.Date(cur_year, cur_month, cur_day, scaleUpHour, scaleUpMin, scaleUpSeconds, 0, time.Local)
	ScaleDownTime := time.Date(cur_year, cur_month, cur_day, scaleDownHour, scaleDownMin, scaleDownSeconds, 0, time.Local)

	// 7. Trigger scale action if required:
	// check which action to take - later or earlier? (or later from a previous day?):
	whichAction := func(earlierTime time.Time, laterTime time.Time, curTime time.Time) (actionType string) {

		if curTime.After(laterTime) {
			actionType = "later"
		} else if curTime.After(earlierTime) {
			actionType = "earlier"
		} else {
			actionType = "laterYesterday"
		}

		return actionType
	}

	// scaleup funciton - scale only if current setup doesnt match required scale value:
	scaleResource := func(scaleValue *int32, resourceType string, deploymentSpec *appsv1.Deployment, hpaSpec *kautoscalingv1.HorizontalPodAutoscaler) (required bool, err error) {
		switch resourceType {
		case "deployment":
			if *scaleValue == *deploymentSpec.Spec.Replicas {
				return false, nil
			} else {
				u := &unstructured.Unstructured{}

				deploymentSpec.Spec.Replicas = scaleValue

				var convErr error

				u.Object, convErr = runtime.DefaultUnstructuredConverter.ToUnstructured(&deploymentSpec)
				if convErr != nil {
					return false, convErr
				}

				updateErr := r.Update(ctx, u)
				if updateErr != nil {
					return false, updateErr
				}
				return true, nil
			}
		case "hpa":
			if *scaleValue == *hpaSpec.Spec.MinReplicas {
				return false, nil
			} else {
				u := &unstructured.Unstructured{}

				hpaSpec.Spec.MinReplicas = scaleValue

				if *scaleValue > hpaSpec.Spec.MaxReplicas {
					log.V(1).Info("maxReplicas is lower than required scaling and new minReplicas", "maxReplicas", hpaSpec.Spec.MaxReplicas, "new minReplicas", scaleValue)
					log.V(1).Info("setting maxReplicas to new minReplicas", "maxReplicas", scaleValue)
					hpaSpec.Spec.MaxReplicas = *scaleValue
				}

				var convErr error

				u.Object, convErr = runtime.DefaultUnstructuredConverter.ToUnstructured(&hpaSpec)
				if convErr != nil {
					return false, convErr
				}

				updateErr := r.Update(ctx, u)
				if updateErr != nil {
					return false, updateErr
				}
				return true, nil
			}
		}
		scaleErr := fmt.Errorf("Failed to update resource %s - its neither a 'deployment' nor 'hpa'", resourceType)
		return false, scaleErr
	}

	switch ScaleUpTime.Before(ScaleDownTime) {
	// when scaleup is before scaledown
	case true:
		actionType := whichAction(ScaleUpTime, ScaleDownTime, curr_time)
		switch actionType {
		case "earlier":
			log.V(1).Info("Based on current time - current replicas must match ScaleUp.Value", "pods", scaleUpValue)
			log.V(1).Info("Checking if scaleUp is required and taking actions if necessairy")
			requiredReplicas = scaleUpValue
		case "later":
			log.V(1).Info("Based on current time - current replicas must match ScaleDown.Value", "pods", scaleDownValue)
			log.V(1).Info("Checking if scaleDown is required and taking actions if necessairy")
			requiredReplicas = scaleDownValue
		case "laterYesterday":
			log.V(1).Info("Based on current time - no actions are required for today. Current replicas must match ScaleDown.Value", "pods from yesterday", scaleDownValue)
			log.V(1).Info("Checking if scaleDown is required and taking actions if necessairy")
			requiredReplicas = scaleDownValue
		}
	// when scaleup is after scaledown
	case false:
		actionType := whichAction(ScaleDownTime, ScaleUpTime, curr_time)
		switch actionType {
		case "earlier":
			log.V(1).Info("Based on current time - current replicas must match ScaleDown.Value", "pods", scaleDownValue)
			log.V(1).Info("Checking if scaleDown is required and taking actions if necessairy")
			requiredReplicas = scaleDownValue
		case "later":
			log.V(1).Info("Based on current time - current replicas must match ScaleUp.Value", "pods", scaleUpValue)
			log.V(1).Info("Checking if scaleUp is required and taking actions if necessairy")
			requiredReplicas = scaleUpValue
		case "laterYesterday":
			log.V(1).Info("Based on current time - no actions are required for today. Current replicas must match scaleUpValue.Value", "pods from yesterday", scaleUpValue)
			log.V(1).Info("Checking if scaleUp is required and taking actions if necessairy")
			requiredReplicas = scaleUpValue
		}
	}

	// check if scaleup is required - trigger the scaleResource func:
	requiredScaling, err := scaleResource(requiredReplicas, resourceType, deploymentSpec, hpaSpec)

	// log outcome:
	if err != nil {
		log.Error(err, "unable to scale resource", "type", resourceType, "named", passedResourceName)
	} else if requiredScaling {
		log.V(1).Info("Scaling process was required and contoller successfully scaled to", "podsCount", requiredReplicas)
	} else {
		log.V(1).Info("Replica count already matched required setup with", "podsCount alreadt at", requiredReplicas)
	}

	// 8. Requeue reconciliation and return to manager:

	// retrieve the rate of requeuing reconciliation loop:
	requeueRate := os.Getenv("RequeueRate")
	var requeueRateNS time.Duration
	var prsDurErr error

	if requeueRate == "" {
		requeueRateNS, prsDurErr = time.ParseDuration("10s")
	} else {
		requeueRateNS, prsDurErr = time.ParseDuration(requeueRate)
	}

	if prsDurErr != nil {
		log.Error(prsDurErr, "unable parse duration for reconcilation requeue rate", "value", requeueRate)
	}

	// return to manager if no errors occured along the way:
	return ctrl.Result{RequeueAfter: requeueRateNS}, nil
}

func (r *ScheduledPodAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.ScheduledPodAutoscaler{}).
		Complete(r)
}
