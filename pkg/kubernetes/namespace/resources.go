package namespace

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	puluminamekubeoutput "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/provider/kubernetes/output"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubernetesProvider *pulumikubernetes.Provider
	MongodbClusterId   string
	Labels             map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) (*kubernetescorev1.Namespace, error) {
	namespace, err := addNamespace(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}
	return namespace, nil
}

func addNamespace(ctx *pulumi.Context, input *Input) (*kubernetescorev1.Namespace, error) {
	ns, err := kubernetescorev1.NewNamespace(ctx, input.MongodbClusterId, &kubernetescorev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("AddedNamespace"),
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:   pulumi.String(input.MongodbClusterId),
			Labels: pulumi.ToStringMap(input.Labels),
		}),
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "3m", Update: "3m", Delete: "3m"}),
		pulumi.Provider(input.KubernetesProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s namespace", input.MongodbClusterId)
	}
	ctx.Export(GetNamespaceNameOutputName(), ns.Metadata.Name())
	return ns, nil
}

func GetNamespaceNameOutputName() string {
	return puluminamekubeoutput.Name(kubernetescorev1.Namespace{}, englishword.EnglishWord_namespace.String())
}
