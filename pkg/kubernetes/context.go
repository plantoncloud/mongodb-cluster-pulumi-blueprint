package gcp

import (
	"github.com/pkg/errors"
	environmentblueprinthostnames "github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/hostnames"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbnetutilshostname "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/hostname"
	"github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/service"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	plantoncloudpulumisdkkubernetes "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/automation/provider/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func loadConfig(ctx *pulumi.Context, resourceStack *ResourceStack) (*mongodbcontextconfig.ContextConfig, error) {

	kubernetesProvider, err := plantoncloudpulumisdkkubernetes.GetWithStackCredentials(ctx, resourceStack.Input.CredentialsInput.Kubernetes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup kubernetes provider")
	}

	var resourceId = resourceStack.Input.ResourceInput.MongodbCluster.Metadata.Id
	var resourceName = resourceStack.Input.ResourceInput.MongodbCluster.Metadata.Name
	var environmentInfo = resourceStack.Input.ResourceInput.MongodbCluster.Spec.EnvironmentInfo
	var isIngressEnabled = false

	if resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress != nil {
		isIngressEnabled = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.IsEnabled
	}

	var endpointDomainName = ""
	var envDomainName = ""
	var ingressType = kubernetesworkloadingresstype.KubernetesWorkloadIngressType_unspecified
	var internalHostname = ""
	var externalHostname = ""

	if isIngressEnabled {
		endpointDomainName = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.EndpointDomainName
		envDomainName = environmentblueprinthostnames.GetExternalEnvHostname(environmentInfo.EnvironmentName, endpointDomainName)
		ingressType = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.IngressType

		internalHostname = mongodbnetutilshostname.GetInternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
		externalHostname = mongodbnetutilshostname.GetExternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
	}

	return &mongodbcontextconfig.ContextConfig{
		Spec: &mongodbcontextconfig.Spec{
			KubeProvider:       kubernetesProvider,
			ResourceId:         resourceId,
			ResourceName:       resourceName,
			Labels:             resourceStack.KubernetesLabels,
			WorkspaceDir:       resourceStack.WorkspaceDir,
			NamespaceName:      resourceId,
			EnvironmentInfo:    resourceStack.Input.ResourceInput.MongodbCluster.Spec.EnvironmentInfo,
			IsIngressEnabled:   isIngressEnabled,
			IngressType:        ingressType,
			EndpointDomainName: endpointDomainName,
			EnvDomainName:      envDomainName,
			ContainerSpec:      resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.MongodbContainer,
			CustomHelmValues:   resourceStack.Input.ResourceInput.MongodbCluster.Spec.HelmValues,
			InternalHostname:   internalHostname,
			ExternalHostname:   externalHostname,
			KubeServiceName:    service.GetKubeServiceName(resourceName),
			KubeLocalEndpoint:  service.GetKubeServiceNameFqdn(resourceName, resourceId),
		},
		Status: &mongodbcontextconfig.Status{},
	}, nil
}
