package mongodbcluster

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
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

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	var helmValues = getHelmValues(ctx, input.ContainerSpec, input.Values)
	// Deploying a Mongodb Helm chart from the Helm repository.
	_, err := helmv3.NewChart(ctx, ctxConfig.Spec.ResourceId, helmv3.ChartArgs{
		Chart:     pulumi.String("mongodb"),
		Version:   pulumi.String("15.1.4"), // Use the Helm chart version you want to install
		Namespace: ctxConfig.Status.AddedResources.Namespace.Metadata.Name().Elem(),
		Values:    helmValues,
		//if you need to add the repository, you can specify `repo url`:
		// The URL for the Helm chart repository
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
		},
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "3m", Update: "3m", Delete: "3m"}),
		pulumi.Parent(ctxConfig.Status.AddedResources.Namespace))
	if err != nil {
		return err
	}
	ctx.Export(ctxConfig.Status.OutputKeyNames.RootPasswordSecret, pulumi.String(ctxConfig.Status.OutputValues.RootPasswordSecret))
	ctx.Export(ctxConfig.Status.OutputKeyNames.RootUsername, pulumi.String(ctxConfig.Status.OutputValues.RootUsername))
	return nil
}
