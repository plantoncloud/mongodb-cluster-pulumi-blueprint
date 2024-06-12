package mongodbcluster

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	plantoncloudmongodbmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/model"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	ResourceId       string
	Namespace        *kubernetescorev1.Namespace
	KubeServiceName  string
	ContainerSpec    *plantoncloudmongodbmodel.MongodbClusterSpecKubernetesSpecMongodbContainerSpec
	CustomHelmValues map[string]string
	Labels           map[string]string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		ResourceId:       ctxConfig.Spec.ResourceId,
		Namespace:        ctxConfig.Status.AddedResources.Namespace,
		ContainerSpec:    ctxConfig.Spec.ContainerSpec,
		CustomHelmValues: ctxConfig.Spec.CustomHelmValues,
		Labels:           ctxConfig.Spec.Labels,
		KubeServiceName:  ctxConfig.Spec.KubeServiceName,
	}
}
