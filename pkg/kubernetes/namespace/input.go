package namespace

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	NamespaceName string
	Labels        map[string]string
	KubeProvider  *pulumikubernetes.Provider
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		NamespaceName: ctxConfig.Spec.NamespaceName,
		Labels:        ctxConfig.Spec.Labels,
		KubeProvider:  ctxConfig.Spec.KubeProvider,
	}
}
