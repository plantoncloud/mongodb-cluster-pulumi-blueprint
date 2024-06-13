package network

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	IsIngressEnabled   bool
	EndpointDomainName string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		IsIngressEnabled:   ctxConfig.Spec.IsIngressEnabled,
		EndpointDomainName: ctxConfig.Spec.EndpointDomainName,
	}
}
