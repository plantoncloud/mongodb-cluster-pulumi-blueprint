package loadbalancer

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/loadbalancer/gcp"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/enums/kubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (newCtx *pulumi.Context, err error) {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	if ctxConfig.Spec.EnvironmentInfo.KubernetesProvider == kubernetesprovider.KubernetesProvider_gcp_gke {
		ctx, err = gcp.Resources(ctx)
		if err != nil {
			return ctx, errors.Wrap(err, "failed to create load balancer resources for gke cluster")
		}
	}
	return ctx, nil
}
