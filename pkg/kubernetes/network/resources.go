package network

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbingress "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress"
	mongodbnetutilendpoint "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/endpoint"
	mongodbnetutilsport "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/port"
	mongodbnetutilsservice "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/service"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (newCtx *pulumi.Context, err error) {

	if err := mongodbnetutilsservice.Resources(ctx); err != nil {
		return ctx, errors.Wrap(err, "failed to add service resources")
	}
	if err := mongodbnetutilsport.Resources(ctx); err != nil {
		return ctx, errors.Wrap(err, "failed to add port resources")
	}
	if err := mongodbnetutilendpoint.Resources(ctx); err != nil {
		return ctx, errors.Wrap(err, "failed to add endpoint resources")
	}

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	if !ctxConfig.Spec.IsIngressEnabled || ctxConfig.Spec.EndpointDomainName == "" {
		return ctx, nil
	}
	if ctx, err = mongodbingress.Resources(ctx); err != nil {
		return ctx, errors.Wrap(err, "failed to add gateway resources")
	}
	return ctx, nil
}
