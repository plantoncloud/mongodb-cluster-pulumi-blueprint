package loadbalancer

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/loadbalancer/gcp"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/enums/kubernetesprovider"
	code2cloudv1envmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/environment/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	EnvironmentInfo    *code2cloudv1envmodel.ApiResourceEnvironmentInfo
	EndpointDomainName string
}

func Resources(ctx *pulumi.Context) error {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	if ctxConfig.Spec.EnvironmentInfo.KubernetesProvider == kubernetesprovider.KubernetesProvider_gcp_gke {
		if err := gcp.Resources(ctx); err != nil {
			return errors.Wrap(err, "failed to create load balancer resources for gke cluster")
		}
	}
	return nil
}
