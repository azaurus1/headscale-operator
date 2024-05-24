package controller

import (
	"context"

	headscalev1 "github.com/azaurus1/headscale-operator/api/v1"
)

func (r *ApiKeyReconciler) createAPIKey(ctx context.Context, apiKey *headscalev1.ApiKey, timeToExpire int32) (string, error) {
	// create the api key in headscale
	return "", nil
}

func (r *ApiKeyReconciler) createAPISecret(ctx context.Context, apiKey string) (headscalev1.NamespacedName, error) {
	// create a secret in our namespace
	// store the apiKey in there

	return headscalev1.NamespacedName{}, nil
}

func (r *ApiKeyReconciler) DeleteExternalResources(ctx context.Context, apiKey *headscalev1.ApiKey) error {
	// this is where use the /api/v1/apikey/expire endpoint
	// but we dont have the prefix...
	//
	//
	return nil
}
