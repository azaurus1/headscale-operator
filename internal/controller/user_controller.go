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
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	finalizerName = "headscale.azaurus.dev/finalizer"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=headscale.azaurus.dev,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	user := &headscalev1.User{}

	log.Info("Reconciling user")

	// fetch user instance
	if err := r.Get(ctx, req.NamespacedName, user); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Ensure the user with the name exists, if not, create it
	if err := r.ensureUser(ctx, user, user.Spec.Name); err != nil {
		log.Error(err, "unable to ensure user", "user", user.Spec.ID)
	}

	// update the user status with current state
	user.Status.CreatedAt = user.Spec.CreatedAt
	if err := r.Status().Update(ctx, user); err != nil {
		log.Error(err, "unable to update User status")
		return ctrl.Result{}, err
	}

	if user.DeletionTimestamp != nil {
		if controllerutil.ContainsFinalizer(user, finalizerName) {
			// Cleanup Resources
			log.Info("Finalizer found, cleaning up resources")
			if err := r.DeleteExternalResources(ctx, user); err != nil {
				// retry if failed
				log.Error(err, "failed to cleanup resources")
				return ctrl.Result{}, err
			}
			// Remove the finalizer
			controllerutil.RemoveFinalizer(user, finalizerName)
			if err := r.Update(ctx, user); err != nil {
				log.Error(err, "unable to update user")
				return ctrl.Result{}, err
			}
		}
	} else {
		// user is not being deleted, add the finaliser if not present
		if !controllerutil.ContainsFinalizer(user, finalizerName) {
			user.ObjectMeta.Finalizers = append(user.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, user); err != nil {
				log.Error(err, "unabled to update user")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *UserReconciler) ensureUser(ctx context.Context, user *headscalev1.User, userName string) error {
	log := log.FromContext(ctx)
	log.Info("This is where we create the user headscale")

	headscaleDeployment := appsv1.Deployment{}

	// check for which headscale server we are deploying this to (user.spec.headscaleServerRef(name/namespace))
	err := r.Get(ctx, client.ObjectKey{Name: user.Spec.HeadscaleServerRef.Name, Namespace: user.Spec.HeadscaleServerRef.Namespace}, &headscaleDeployment)
	if err != nil {
		// deployment doesnt exist -> headscale isnt deployed, error
		if apierrors.IsNotFound(err) {
			log.Error(err, "headscale server is not deployed")
			return err
		}
	}
	if !(headscaleDeployment.Status.ReadyReplicas != int32(0) && headscaleDeployment.Status.ReadyReplicas > 0) {
		log.Error(err, "headscale server has no working replicas")
		return err
	}
	// we should make sure the user doesnt exist already, if it does we will send a PUT instead

	// use the service name e.g. svc-[headscale.spec.name].[headscale.metadata.namespace].svc.cluster.local to communicate - api/v1/user
	err = r.CreateUserViaService(ctx, user)
	if err != nil {
		log.Error(err, "could not create a user via the service")
		return err
	}

	return nil
}

func (r *UserReconciler) DeleteExternalResources(ctx context.Context, user *headscalev1.User) error {
	log := log.FromContext(ctx)

	log.Info("This is where we delete a user from headscale")

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&headscalev1.User{}).
		Complete(r)
}
