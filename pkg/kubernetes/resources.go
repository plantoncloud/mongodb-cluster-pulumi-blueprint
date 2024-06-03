package gcp

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbclusterresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/mongodbcluster"
	mongodbnamespaceresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/namespace"
	mongodbnetworkresources "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network"
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
	addedNameSpace, err := mongodbnamespaceresources.Resources(ctx)
	if err != nil {
		return err
	}

	AddNameSpace(ctxConfig, addedNameSpace)
	ctx = ctx.WithValue(mongodbcontextconfig.Key, *ctxConfig)

	// Deploying a Mongodb Helm chart from the Helm repository.
	err = mongodbclusterresources.Resources(ctx, &mongodbclusterresources.Input{
		ContainerSpec: mongodbCluster.Spec.Kubernetes.MongodbContainer,
		Values:        mongodbCluster.Spec.HelmValues,
	})
	if err != nil {
		return err
	}

	// Deploying a Mongodb Helm chart from the Helm repository.
	err = mongodbnetworkresources.Resources(ctx)
	if err != nil {
		return err
	}

	return nil
}
