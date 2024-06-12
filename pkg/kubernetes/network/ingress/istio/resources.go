package istio

import (
	"github.com/pkg/errors"
	mongodbistiogateway "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/istio/gateway"
	mongodbistiovirtualservice "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/istio/virtualservice"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	if err := mongodbistiogateway.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add gateway resources")
	}
	if err := mongodbistiovirtualservice.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add virtual-service resources")
	}
	return nil
}
