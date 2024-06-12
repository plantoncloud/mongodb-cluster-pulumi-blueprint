package virtualservice

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	ResourceId        string
	Namespace         *kubernetescorev1.Namespace
	HostNames         []string
	WorkspaceDir      string
	NameSpaceName     string
	KubeLocalEndpoint string
	KubeServiceName   string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		ResourceId:        ctxConfig.Spec.ResourceId,
		Namespace:         ctxConfig.Status.AddedResources.Namespace,
		HostNames:         []string{ctxConfig.Spec.ExternalHostname},
		WorkspaceDir:      ctxConfig.Spec.WorkspaceDir,
		NameSpaceName:     ctxConfig.Spec.NamespaceName,
		KubeLocalEndpoint: ctxConfig.Spec.KubeLocalEndpoint,
		KubeServiceName:   ctxConfig.Spec.KubeServiceName,
	}
}
