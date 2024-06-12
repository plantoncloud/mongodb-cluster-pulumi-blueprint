package ingress

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	IngressType kubernetesworkloadingresstype.KubernetesWorkloadIngressType
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	return &input{
		IngressType: ctxConfig.Spec.IngressType,
	}
}
