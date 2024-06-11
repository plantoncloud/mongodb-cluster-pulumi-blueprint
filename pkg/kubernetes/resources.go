package gcp

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbclusterresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/mongodbcluster"
	mongodbnamespaceresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/namespace"
	mongodbnetworkresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network"
	mongodboutputs "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/outputs"
	model "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/stack/kubernetes/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	WorkspaceDir     string
	Input            *model.MongodbClusterKubernetesStackInput
	KubernetesLabels map[string]string
}

func (resourceStack *ResourceStack) Resources(ctx *pulumi.Context) error {
	// https://artifacthub.io/packages/helm/bitnami/mongodb
	var mongodbCluster = resourceStack.Input.ResourceInput.MongodbCluster

	var ctxConfig, err = loadConfig(ctx, resourceStack)
	if err != nil {
		return errors.Wrap(err, "failed to initiate context config")
	}
	ctx = ctx.WithValue(mongodbcontextconfig.Key, *ctxConfig)

	// Create the namespace resource
	ctx, err = mongodbnamespaceresources.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace resource")
	}

	// Deploying a Mongodb Helm chart from the Helm repository.
	err = mongodbclusterresources.Resources(ctx, &mongodbclusterresources.Input{
		ContainerSpec: mongodbCluster.Spec.Kubernetes.MongodbContainer,
		Values:        mongodbCluster.Spec.HelmValues,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create mongodb cluster")
	}

	// Deploying a Mongodb Helm chart from the Helm repository.
	ctx, err = mongodbnetworkresources.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create mongodb network resources")
	}

	err = mongodboutputs.Export(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to extract mongodb cluster outputs")
	}

	return nil
}
