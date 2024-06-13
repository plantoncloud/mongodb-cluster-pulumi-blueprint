package mongodbcluster

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	err := addHelmChart(ctx)
	if err != nil {
		return err
	}
	return nil
}

func addHelmChart(ctx *pulumi.Context) error {
	var i = extractInput(ctx)

	var helmValues = getHelmValues(i)
	ctx.Export("helm-values", helmValues)
	// Deploying a Mongodb Helm chart from the Helm repository.
	_, err := helmv3.NewChart(ctx, i.ResourceId, helmv3.ChartArgs{
		Chart:     pulumi.String("mongodb"),
		Version:   pulumi.String("15.1.4"), // Use the Helm chart version you want to install
		Namespace: pulumi.String("mdb-planton-cloud-prod-test-ingress-controller"),
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
