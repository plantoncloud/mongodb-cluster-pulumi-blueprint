package gcp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	mongodbclustermodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/model"
	mongodbclusterstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/stack/kubernetes/model"
)

func Outputs(ctx context.Context, input *mongodbclusterstackmodel.MongodbClusterKubernetesStackInput) (*mongodbclusterstackmodel.MongodbClusterKubernetesStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return OutputMapTransformer(stackOutput, input), nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *mongodbclusterstackmodel.MongodbClusterKubernetesStackInput) *mongodbclusterstackmodel.MongodbClusterKubernetesStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &mongodbclusterstackmodel.MongodbClusterKubernetesStackOutputs{}
	}
	return &mongodbclusterstackmodel.MongodbClusterKubernetesStackOutputs{
		MongodbClusterStatus: &mongodbclustermodel.MongodbClusterStatus{
			Kubernetes: &mongodbclustermodel.MongodbClusterStatusKubernetesStatus{
				Namespace:              backend.GetVal(stackOutput, GetNamespaceNameOutputName()),
				RootUsername:           backend.GetVal(stackOutput, GetRootUsernameOutputName()),
				RootPasswordSecretName: backend.GetVal(stackOutput, GetRootPasswordOutputName()),
				//Service:                backend.GetVal(stackOutput, GetKubeServiceNameOutputName()),
				//PortForwardCommand:     backend.GetVal(stackOutput, GetKubePortForwardCommandOutputName()),
				//KubeEndpoint:           backend.GetVal(stackOutput, GetKubernetesEndpointOutputName()),
				//IngressEndpoint:        backend.GetVal(stackOutput, GetIngressEndpointOutputName()),
			},
		},
	}
}
