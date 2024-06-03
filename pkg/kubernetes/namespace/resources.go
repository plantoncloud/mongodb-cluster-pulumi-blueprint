package namespace

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*kubernetescorev1.Namespace, error) {
	namespace, err := addNamespace(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}
	return namespace, nil
}

func addNamespace(ctx *pulumi.Context) (*kubernetescorev1.Namespace, error) {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	ns, err := kubernetescorev1.NewNamespace(ctx, ctxConfig.Spec.NamespaceName, &kubernetescorev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("AddedNamespace"),
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:   pulumi.String(ctxConfig.Spec.NamespaceName),
			Labels: pulumi.ToStringMap(ctxConfig.Spec.Labels),
		}),
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "3m", Update: "3m", Delete: "3m"}),
		pulumi.Provider(ctxConfig.Spec.KubeProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s namespace", ctxConfig.Spec.NamespaceName)
	}
	ctx.Export(ctxConfig.Status.OutputKeyNames.Namespace, pulumi.String(ctxConfig.Status.OutputValues.Namespace))
	return ns, nil
}
