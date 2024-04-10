package mongodbcluster

import (
	"fmt"
	plantoncloudmongodbmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/model"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/output/custom"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AddedNamespace   *kubernetescorev1.Namespace
	MongodbClusterId string
	ContainerSpec    *plantoncloudmongodbmodel.MongodbClusterSpecKubernetesSpecMongodbContainerSpec
	Labels           map[string]string
	Values           map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) error {
	err := addHelmChart(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func addHelmChart(ctx *pulumi.Context, input *Input) error {
	var helmValues = getHelmValues(input.ContainerSpec, input.Values, input.Labels)
	// Deploying a Locust Helm chart from the Helm repository.
	_, err := helmv3.NewChart(ctx, input.MongodbClusterId, helmv3.ChartArgs{
		Chart:     pulumi.String("mongodb"),
		Version:   pulumi.String("15.1.4"), // Use the Helm chart version you want to install
		Namespace: input.AddedNamespace.Metadata.Name().Elem(),
		Values:    helmValues,
		//if you need to add the repository, you can specify `repo url`:
		// The URL for the Helm chart repository
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
		},
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "3m", Update: "3m", Delete: "3m"}),
		pulumi.Parent(input.AddedNamespace))
	if err != nil {
		return err
	}
	ctx.Export(GetRootPasswordOutputName(), pulumi.String(GetRootPasswordSecretName(input.MongodbClusterId)))
	ctx.Export(GetRootUsernameOutputName(), pulumi.String(GetRootUsernameSecretName()))
	return nil
}

func GetRootUsernameOutputName() string {
	return custom.Name("mongodb-cluster-root-username")
}

func GetRootUsernameSecretName() string {
	return custom.Name("root")
}

func GetRootPasswordOutputName() string {
	return custom.Name("mongodb-cluster-root-password-secret-name")
}

func GetRootPasswordSecretName(mongodbClusterId string) string {
	return fmt.Sprintf("%s-mongodb", mongodbClusterId)
}
