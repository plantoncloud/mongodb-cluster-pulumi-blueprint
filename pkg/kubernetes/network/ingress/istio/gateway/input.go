package gateway

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	WorkspaceDir     string
	KubeProvider     *pulumikubernetes.Provider
	ResourceId       string
	Labels           map[string]string
	EnvDomainName    string
	ExternalHostname string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		WorkspaceDir:     ctxConfig.Spec.WorkspaceDir,
		KubeProvider:     ctxConfig.Spec.KubeProvider,
		ResourceId:       ctxConfig.Spec.ResourceId,
		Labels:           ctxConfig.Spec.Labels,
		EnvDomainName:    ctxConfig.Spec.EnvDomainName,
		ExternalHostname: ctxConfig.Spec.ExternalHostname,
	}
}
