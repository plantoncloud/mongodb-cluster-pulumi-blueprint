package secret

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	MongodbRootPasswordKey = "mongodb-root-password"
)

type input struct {
	NamespaceName string
	ResourceName  string
	Labels        map[string]string
	KubeProvider  *pulumikubernetes.Provider
	Namespace     *kubernetescorev1.Namespace
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		NamespaceName: ctxConfig.Spec.NamespaceName,
		Labels:        ctxConfig.Spec.Labels,
		KubeProvider:  ctxConfig.Spec.KubeProvider,
		ResourceName:  ctxConfig.Spec.ResourceName,
		Namespace:     ctxConfig.Status.AddedResources.Namespace,
	}
}
