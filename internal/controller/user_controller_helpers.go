package controller

import (
	"context"

	headscalev1 "github.com/azaurus1/headscale-operator/api/v1"
	v1 "github.com/azaurus1/headscale-operator/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sigs.k8s.io/controller-runtime/pkg/log"
)



func (r *UserReconciler) CreateUserViaService(ctx context.Context, user *headscalev1.User) error {
	// REST is not enough, we need to use gRPC
	log := log.FromContext(ctx)
	log.Info("Attempting to create a user via gRPC")

	userReq := &v1.CreateUserRequest{
		Name: user.Name,
	}

	grpcOptions := []grpc.DialOption{
		grpc.WithBlock(),
	}

	grpcOptions = append(grpcOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// grpcTarget := fmt.Sprintf("svc-" + user.Spec.HeadscaleServerRef.Name + "." + user.Spec.HeadscaleServerRef.Namespace + ".svc.cluster.local:50443")
	grpcTarget := "localhost:8082" // port forwarding to run via `make run`

	log.Info(grpcTarget)

	conn, err := grpc.DialContext(ctx, grpcTarget, grpcOptions...)
	if err != nil {
		log.Error(err, "error dialing the grpc server")
		return err
	}

	log.Info("grpc connection established..")

	client := v1.NewHeadscaleServiceClient(conn)

	_, err = client.CreateUser(ctx, userReq)
	if err != nil {
		log.Error(err, "error creating the user via the grpc server")
		return err
	}

	// log.Info("user with id:", resp.User.Id)

	return nil
}
