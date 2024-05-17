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
)

const (
	headscaleOperatorAnnotation = "headscale-operator"
)

// HeadscaleReconciler reconciles a Headscale object
type HeadscaleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=headscales,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=headscales/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=headscales/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Headscale object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *HeadscaleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	server := &headscalev1.Headscale{}

	log.Info("Reconciling Headscale Server..")

	if err := r.Get(ctx, req.NamespacedName, server); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.EnsureHeadscaleServer(ctx, server, server.Spec.Name, server.Spec.Version, server.Spec.Config); err != nil {
		log.Error(err, "unable to ensure headscale server")
		return ctrl.Result{}, err
	}

	if server.DeletionTimestamp != nil {
		if controllerutil.ContainsFinalizer(server, finalizerName) {
			log.Info("finalizer found, cleaning up..")
		}
		err := r.DeleteExternalResources(ctx, server)
		if err != nil {
			log.Error(err, "failed to clean up resources")
			return ctrl.Result{}, err
		}

		// remove the finalizer
		controllerutil.RemoveFinalizer(server, finalizerName)
		err = r.Update(ctx, server)
		if err != nil {
			log.Error(err, "unable to update headscale server")
			return ctrl.Result{}, err
		}
	} else {
		// server is not being deleted
		if !controllerutil.ContainsFinalizer(server, finalizerName) {
			server.ObjectMeta.Finalizers = append(server.ObjectMeta.Finalizers, finalizerName)
			err := r.Update(ctx, server)
			if err != nil {
				log.Error(err, "unabled to update headscale server")
				return ctrl.Result{}, err
			}
		}
	}

	// update the headscale server status with the current state
	server.Status.Version = server.Spec.Version
	if err := r.Status().Update(ctx, server); err != nil {
		log.Error(err, "unable to update headscale server status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HeadscaleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&headscalev1.Headscale{}).
		Complete(r)
}
