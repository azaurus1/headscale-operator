/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	headscalev1 "github.com/azaurus1/headscale-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApiKeyReconciler reconciles a ApiKey object
type ApiKeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=apikeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=apikeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=apikeys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApiKey object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ApiKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	apiKey := &headscalev1.ApiKey{}

	log.Info("Reconcilling API Key")

	if err := r.Get(ctx, req.NamespacedName, apiKey); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// we cant ensure, because we dont have a primary key in the spec
	// therefore, we just create a key

	apiSecretKey, err := r.createAPIKey(ctx, apiKey, apiKey.Spec.TimeToExpire)
	if err != nil {
		log.Error(err, "unable to create the api key", "api key", apiKey.Spec.TimeToExpire)
	}

	// we have now created the key
	// we now need to create a secret, and put the key in there

	secretName, err := r.createAPISecret(ctx, apiSecretKey)
	if err != nil {
		log.Error(err, "unable to create the api secret")
	}

	// set the secret name and the createdAt in the status
	apiKey.Status.CreatedAt = metav1.Now()
	apiKey.Status.KeySecret = secretName.Name

	// otherwise, if the deletion timestamp is set, delete this
	//

	if apiKey.DeletionTimestamp != nil {
		// Cleanup Resources
		log.Info("Finalizer found, cleaning up resources")
		if err := r.DeleteExternalResources(ctx, apiKey); err != nil {
			// retry if failed
			log.Error(err, "failed to cleanup resources")
			return ctrl.Result{}, err
		}
		// Remove the finalizer
		controllerutil.RemoveFinalizer(apiKey, finalizerName)
		if err := r.Update(ctx, apiKey); err != nil {
			log.Error(err, "unable to update apiKey")
			return ctrl.Result{}, err
		}
	} else {
		// key is not being deleted, add the finaliser if not present
		if !controllerutil.ContainsFinalizer(apiKey, finalizerName) {
			apiKey.ObjectMeta.Finalizers = append(apiKey.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, apiKey); err != nil {
				log.Error(err, "unabled to update apiKey")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&headscalev1.ApiKey{}).
		Complete(r)
}
