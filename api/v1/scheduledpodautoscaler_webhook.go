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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var scheduledpodautoscalerlog = logf.Log.WithName("scheduledpodautoscaler-resource")

func (r *ScheduledPodAutoscaler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-autoscaling-spa-sarmadabualkaz-io-v1-scheduledpodautoscaler,mutating=true,failurePolicy=fail,groups=autoscaling.spa.sarmadabualkaz.io,resources=scheduledpodautoscalers,verbs=create;update,versions=v1,name=mscheduledpodautoscaler.kb.io

var _ webhook.Defaulter = &ScheduledPodAutoscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ScheduledPodAutoscaler) Default() {
	scheduledpodautoscalerlog.Info("default", "name", r.Name)

	// default 'Spec.Resource.Type' to 'deployment' if set blank
	if r.Spec.Resource.Type == "" {
		r.Spec.Resource.Type = "deployment"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-autoscaling-spa-sarmadabualkaz-io-v1-scheduledpodautoscaler,mutating=false,failurePolicy=fail,groups=autoscaling.spa.sarmadabualkaz.io,resources=scheduledpodautoscalers,versions=v1,name=vscheduledpodautoscaler.kb.io

var _ webhook.Validator = &ScheduledPodAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ScheduledPodAutoscaler) ValidateCreate() error {
	scheduledpodautoscalerlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateScheduledPodAutoscaler()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ScheduledPodAutoscaler) ValidateUpdate(old runtime.Object) error {
	scheduledpodautoscalerlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateScheduledPodAutoscaler()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ScheduledPodAutoscaler) ValidateDelete() error {
	scheduledpodautoscalerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *ScheduledPodAutoscaler) validateScheduledPodAutoscaler() error {
	var allErrs field.ErrorList
	if err := r.validateScheduledPodAutoscalerSpec(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateScheduledPodAutoscalerTimeEnteries(); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "autoscaling.spa.sarmadabualkaz.io", Kind: "ScheduledPodAutoscaler"},
		r.Name, allErrs)
}

func (r *ScheduledPodAutoscaler) validateScheduledPodAutoscalerSpec() *field.Error {
	// The field helpers from Kubernetes API machinery to return
	// structured validation errors

	if r.Spec.Resource.Name == "" {
		return field.Invalid(field.NewPath("spec").Child("resource").Key("name"), r.Spec.Resource.Name, "name cannot be blank and must be no more than 52 characters")
	} else if *r.Spec.ScaleDown.Value <= 0 {
		return field.Invalid(field.NewPath("spec").Child("scaleDown").Key("value"), r.Spec.ScaleDown.Value, "scalueDown.value is invalid - needs to be at least equal to 1")
	} else if *r.Spec.ScaleUp.Value <= *r.Spec.ScaleDown.Value {
		return field.Invalid(field.NewPath("spec").Child("scaleUp").Key("value"), r.Spec.ScaleUp.Value, "scalueUp.value is invalid - needs to be more than scaleDown.value")
	}
	return nil
}

func (r *ScheduledPodAutoscaler) validateScheduledPodAutoscalerTimeEnteries() *field.Error {
	// The field helpers from Kubernetes API machinery to return
	// structured validation errors

	//check if time is validly entred
	if _, err := time.Parse(time.Kitchen, r.Spec.ScaleUp.Time); err != nil {
		return field.Invalid(field.NewPath("spec").Child("scaleUp").Key("time"), r.Spec.ScaleUp.Time, err.Error())
	}
	if _, err := time.Parse(time.Kitchen, r.Spec.ScaleDown.Time); err != nil {
		return field.Invalid(field.NewPath("spec").Child("scaleDown").Key("time"), r.Spec.ScaleDown.Time, err.Error())
	}
	return nil
}
