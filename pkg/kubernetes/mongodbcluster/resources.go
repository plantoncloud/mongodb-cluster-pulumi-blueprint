package mongodbcluster

import (
	plantoncloudmongodbmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/model"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	ContainerSpec *plantoncloudmongodbmodel.MongodbClusterSpecKubernetesSpecMongodbContainerSpec
	Values        map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) error {
	err := addHelmChart(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func addHelmChart(ctx *pulumi.Context, input *Input) error {
	var i = extractInput(ctx)

	var helmValues = getHelmValues(i)
	// Deploying a Mongodb Helm chart from the Helm repository.
	_, err := helmv3.NewChart(ctx, i.ResourceId, helmv3.ChartArgs{
		Chart:     pulumi.String("mongodb"),
		Version:   pulumi.String("15.1.4"), // Use the Helm chart version you want to install
		Namespace: i.Namespace.Metadata.Name().Elem(),
		Values:    helmValues,
		//if you need to add the repository, you can specify `repo url`:
		// The URL for the Helm chart repository
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
		},
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "3m", Update: "3m", Delete: "3m"}),
		pulumi.Parent(i.Namespace))
	if err != nil {
		return err
	}
	return nil
}
