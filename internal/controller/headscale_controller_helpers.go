package controller

import (
	"context"
	"fmt"

	headscalev1 "github.com/azaurus1/headscale-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *HeadscaleReconciler) EnsureHeadscaleServer(ctx context.Context, server *headscalev1.Headscale, name string, version string, configMapName string) error {
	log := log.FromContext(ctx)

	headscaleServer := &apiv1.Pod{}

	// attempt to get a pod with the provided name and version
	logStr := fmt.Sprintf("Attempting to get a pod with name %s", name)
	log.Info(logStr)
	err := r.Get(ctx, client.ObjectKey{Name: name, Namespace: "default"}, headscaleServer)
	if err != nil {
		// if the pod doesnt exist, create it
		if apierrors.IsNotFound(err) {
			log.Info("Creating Headscale Server Pod")
			headscaleServer := &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default", // using default for now, TODO: Change to namespace defined in server
					Annotations: map[string]string{
						"managed-by": headscaleOperatorAnnotation,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "headscale-server",
							Image: fmt.Sprintf("headscale/headscale:%s", version),
							Args:  []string{"serve"},
							// Add Volume Mount for Config map
							VolumeMounts: []apiv1.VolumeMount{
								{
									// ConfigMap must have data with the name config.yaml
									Name:      "config-volume",
									MountPath: "/etc/headscale",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "config-volume",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
					},
				},
			}

			// attempt to create the namespace
			err = r.Create(ctx, headscaleServer)
			if err != nil {
				return err
			}
		} else {
			return err
		}

	} else {
		// if the pod already exist, check for annotations, etc.
		log.Info("Headscale server pod exists")

		if server.Annotations == nil {
			server.Annotations = map[string]string{}
		}

		// define required annotations and their desired values
		requiredAnnotations := map[string]string{
			"managed-by": headscaleOperatorAnnotation,
		}

		// iterated over required annotations and update if necessary
		for key, desiredValue := range requiredAnnotations {
			existingValue, ok := server.Annotations[key]
			if !ok || existingValue != desiredValue {
				log.Info("Updating headscale annotaion..")
				server.Annotations[key] = desiredValue
				err := r.Update(ctx, server)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *HeadscaleReconciler) DeleteExternalResources(ctx context.Context, server *headscalev1.Headscale) error {
	log := log.FromContext(ctx)

	headscaleServer := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      server.Spec.Name,
			Namespace: "default", // using default for now, TODO: Change to namespace defined in server
		},
	}

	err := r.Delete(ctx, headscaleServer)
	if err != nil {
		log.Error(err, "unable to delete headscale server pod")
		return err
	}

	log.Info("all resources deleted for headscale server")
	return nil

}
