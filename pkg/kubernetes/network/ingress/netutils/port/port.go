package port

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	MongoDbPort = 27017
)

func Resources(ctx *pulumi.Context) error {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	ctx.Export(ctxConfig.Status.OutputKeyNames.KubeForwardCommand, pulumi.String(ctxConfig.Status.OutputValues.KubeForwardCommand))
	return nil
}
