package mongodbcluster

import (
	plantoncloudmongodbmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/model"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/containerresources"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/helm/mergemaps"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func getHelmValues(containerSpec *plantoncloudmongodbmodel.MongodbClusterSpecKubernetesSpecMongodbContainerSpec,
	customValues map[string]string, labels map[string]string) pulumi.Map {
	// https://github.com/bitnami/charts/blob/main/bitnami/mongodb/values.yaml
	var baseValues = pulumi.Map{
		"resources": containerresources.ConvertToPulumiMap(containerSpec.Resources),
		// todo: hard-coding this to 1 since we are only using `standalone` architecture,
		// need to revisit this to handle `replicaSet` architecture
		"replicaCount": pulumi.Int(1),
		"persistence": pulumi.Map{
			"enabled": pulumi.Bool(containerSpec.IsPersistenceEnabled),
			"size":    pulumi.String(containerSpec.DiskSize),
		},
		"podLabels":      pulumi.ToStringMap(labels),
		"commonLabels":   pulumi.ToStringMap(labels),
		"useStatefulSet": pulumi.Bool(true),
	}
	mergemaps.MergeMapToPulumiMap(baseValues, customValues)
	return baseValues
}
