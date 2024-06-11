package gcp

import (
	"github.com/pkg/errors"
	environmentblueprinthostnames "github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/hostnames"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	plantoncloudpulumisdkkubernetes "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/automation/provider/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
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

	if isIngressEnabled {
		endpointDomainName = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.EndpointDomainName
		envDomainName = environmentblueprinthostnames.GetExternalEnvHostname(environmentInfo.EnvironmentName, endpointDomainName)
		ingressType = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.IngressType
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
		},
		Status: &mongodbcontextconfig.Status{},
	}, nil
}

func AddNameSpaceToContext(existingConfig *mongodbcontextconfig.ContextConfig, namespace *kubernetescorev1.Namespace) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &mongodbcontextconfig.AddedResources{
			Namespace: namespace,
		}
		return
	}
	existingConfig.Status.AddedResources.Namespace = namespace
}

func AddLoadBalancerExternalServiceToContext(existingConfig *mongodbcontextconfig.ContextConfig, loadBalancerService *kubernetescorev1.Service) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &mongodbcontextconfig.AddedResources{
			LoadBalancerExternalService: loadBalancerService,
		}
		return
	}
	existingConfig.Status.AddedResources.LoadBalancerExternalService = loadBalancerService
}

func AddLoadBalancerInternalServiceToContext(existingConfig *mongodbcontextconfig.ContextConfig, loadBalancerService *kubernetescorev1.Service) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &mongodbcontextconfig.AddedResources{
			LoadBalancerInternalService: loadBalancerService,
		}
		return
	}
	existingConfig.Status.AddedResources.LoadBalancerInternalService = loadBalancerService
}
