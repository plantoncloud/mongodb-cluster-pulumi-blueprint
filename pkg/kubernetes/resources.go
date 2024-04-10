package gcp

import (
	"github.com/pkg/errors"
	mongodbclusterresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/mongodbcluster"
	mongodbnamespaceresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/namespace"
	model "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/stack/kubernetes/model"
	pulumikubernetesprovider "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/automation/provider/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input *model.MongodbClusterKubernetesStackInput
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	kubernetesProvider, err := pulumikubernetesprovider.GetWithStackCredentials(ctx, s.Input.CredentialsInput.Kubernetes)
	if err != nil {
		return errors.Wrap(err, "failed to setup kubernetes provider")
	}

	var mongodbCluster = s.Input.ResourceInput.MongodbCluster

	// Create the namespace resource
	addedNameSpace, err := mongodbnamespaceresources.Resources(ctx, &mongodbnamespaceresources.Input{
		KubernetesProvider: kubernetesProvider,
		MongodbClusterId:   mongodbCluster.Metadata.Id,
		Labels:             mongodbCluster.Metadata.Labels,
	})
	if err != nil {
		return err
	}

	// Deploying a Locust Helm chart from the Helm repository.
	err = mongodbclusterresources.Resources(ctx, &mongodbclusterresources.Input{
		AddedNamespace:   addedNameSpace,
		ContainerSpec:    mongodbCluster.Spec.Kubernetes.MongodbContainer,
		MongodbClusterId: mongodbCluster.Metadata.Id,
		Labels:           mongodbCluster.Metadata.Labels,
		Values:           mongodbCluster.Spec.HelmValues,
	})
	if err != nil {
		return err
	}

	return nil
}
