package mongodbcluster

import (
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/containerresources"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/helm/mergemaps"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func getHelmValues(i *input) pulumi.Map {
	// https://github.com/bitnami/charts/blob/main/bitnami/mongodb/values.yaml
	var baseValues = pulumi.Map{
		"fullnameOverride":  pulumi.String(i.KubeServiceName),
		"namespaceOverride": i.Namespace.Metadata.Name(),
		"resources":         containerresources.ConvertToPulumiMap(i.ContainerSpec.Resources),
		// todo: hard-coding this to 1 since we are only using `standalone` architecture,
		// need to revisit this to handle `replicaSet` architecture
		"replicaCount": pulumi.Int(1),
		"persistence": pulumi.Map{
			"enabled": pulumi.Bool(i.ContainerSpec.IsPersistenceEnabled),
			"size":    pulumi.String(i.ContainerSpec.DiskSize),
		},
		"podLabels":      pulumi.ToStringMap(i.Labels),
		"commonLabels":   pulumi.ToStringMap(i.Labels),
		"useStatefulSet": pulumi.Bool(true),
		"auth": pulumi.Map{
			"existingSecret": pulumi.String(i.KubeServiceName),
		},
	}
	mergemaps.MergeMapToPulumiMap(baseValues, i.CustomHelmValues)
	return baseValues
}
