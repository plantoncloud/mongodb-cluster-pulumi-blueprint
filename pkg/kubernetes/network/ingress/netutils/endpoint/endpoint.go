package endpoint

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	ctx.Export(ctxConfig.Status.OutputKeyNames.KubeEndpoint, pulumi.String(ctxConfig.Status.OutputValues.KubeEndpoint))
	ctx.Export(ctxConfig.Status.OutputKeyNames.IngressEndpoint, pulumi.String(ctxConfig.Status.OutputValues.IngressEndpoint))
	return nil
}
