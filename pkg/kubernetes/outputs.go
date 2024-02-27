package gcp

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/stack/kubernetes/model"
)

func Outputs(ctx context.Context, input *model.MongoDBClusterKubernetesStackInput) (*model.MongoDBClusterKubernetesStackOutputs, error) {
	return &model.MongoDBClusterKubernetesStackOutputs{}, nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *model.MongoDBClusterKubernetesStackInput) *model.MongoDBClusterKubernetesStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &model.MongoDBClusterKubernetesStackOutputs{}
	}
	return &model.MongoDBClusterKubernetesStackOutputs{}
}
