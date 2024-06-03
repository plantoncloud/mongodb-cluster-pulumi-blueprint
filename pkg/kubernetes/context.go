package gcp

import (
	"fmt"
	"github.com/pkg/errors"
	environmentblueprinthostnames "github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/hostnames"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbnetutilshostname "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/hostname"
	mongodbnetutilsport "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/port"
	mongodbnetutilsservice "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/service"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	plantoncloudpulumisdkkubernetes "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/automation/provider/kubernetes"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/output/custom"
	puluminamekubeoutput "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/provider/kubernetes/output"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const rootUsername = "root"

func loadConfig(ctx *pulumi.Context, resourceStack *ResourceStack) (*mongodbcontextconfig.ContextConfig, error) {

	kubernetesProvider, err := plantoncloudpulumisdkkubernetes.GetWithStackCredentials(ctx, resourceStack.Input.CredentialsInput.Kubernetes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup kubernetes provider")
	}

	var resourceId = resourceStack.Input.ResourceInput.MongodbCluster.Metadata.Id
	var resourceName = resourceStack.Input.ResourceInput.MongodbCluster.Metadata.Name
	var environmentInfo = resourceStack.Input.ResourceInput.MongodbCluster.Spec.EnvironmentInfo
	var isIngressEnabled = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.IsEnabled

	var endpointDomainName = ""
	var ingressEndpoint = ""
	var envDomainName = ""
	var ingressType = kubernetesworkloadingresstype.KubernetesWorkloadIngressType_unspecified

	if isIngressEnabled {
		endpointDomainName = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.EndpointDomainName
		ingressEndpoint = mongodbnetutilshostname.GetExternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
		envDomainName = environmentblueprinthostnames.GetExternalEnvHostname(environmentInfo.EnvironmentName, endpointDomainName)
		ingressType = resourceStack.Input.ResourceInput.MongodbCluster.Spec.Kubernetes.Ingress.IngressType
	}

	return &mongodbcontextconfig.ContextConfig{
		Spec: &mongodbcontextconfig.ContextConfigSpec{
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
		},
		Status: &mongodbcontextconfig.ContextConfigStatus{
			OutputKeyNames: &mongodbcontextconfig.OutputKeyNames{
				IngressEndpoint:               GetIngressEndpointOutputName(),
				KubeServiceName:               GetKubeServiceNameOutputName(),
				KubeEndpoint:                  GetKubernetesEndpointOutputName(),
				KubeForwardCommand:            GetKubePortForwardCommandOutputName(),
				LoadBalancerInternalIpAddress: GetInternalLoadBalancerIp(),
				LoadBalancerExternalIpAddress: GetExternalLoadBalancerIp(),
				Namespace:                     GetNamespaceNameOutputName(),
				RootUsername:                  GetRootUsernameOutputName(),
				RootPasswordSecret:            GetRootPasswordOutputName(),
			},
			OutputValues: &mongodbcontextconfig.OutputValues{
				Namespace:          resourceId,
				IngressEndpoint:    ingressEndpoint,
				KubeServiceName:    resourceName,
				KubeEndpoint:       mongodbnetutilsservice.GetKubeServiceNameFqdn(resourceName, resourceId),
				KubeForwardCommand: getKubePortForwardCommand(resourceId, resourceName),
				RootUsername:       rootUsername,
				RootPasswordSecret: GetRootPasswordSecretName(resourceId),
			},
		},
	}, nil
}

func AddNameSpace(existingConfig *mongodbcontextconfig.ContextConfig, namespace *kubernetescorev1.Namespace) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &mongodbcontextconfig.AddedResources{
			Namespace: namespace,
		}
		return
	}
	existingConfig.Status.AddedResources.Namespace = namespace
}

func GetIngressEndpointOutputName() string {
	return custom.Name("mongodb-cluster-ingress-endpoint")
}

func GetKubeServiceNameOutputName() string {
	return custom.Name("mongodb-cluster-kubernetes-service-name")
}

func GetKubernetesEndpointOutputName() string {
	return custom.Name("mongodb-cluster-kubernetes-endpoint")
}

func GetKubePortForwardCommandOutputName() string {
	return custom.Name("mongodb-cluster-kube-port-forward-command")
}

func GetExternalLoadBalancerIp() string {
	return custom.Name("mongodb-ingress-external-lb-ip")
}

func GetInternalLoadBalancerIp() string {
	return custom.Name("mongodb-ingress-internal-lb-ip")
}

func GetNamespaceNameOutputName() string {
	return puluminamekubeoutput.Name(kubernetescorev1.Namespace{}, englishword.EnglishWord_namespace.String())
}

// getKubePortForwardCommand returns kubectl port-forward command that can be used by developers.
// ex: "kubectl port-forward -n kubernetes_namespace  service/main-mongodb-cluster 8080:8080"
func getKubePortForwardCommand(namespaceName, kubeServiceName string) string {
	return fmt.Sprintf("kubectl port-forward -n %s service/%s %d:%d",
		namespaceName, kubeServiceName, mongodbnetutilsport.MongoDbPort, mongodbnetutilsport.MongoDbPort)
}

func GetRootUsernameOutputName() string {
	return custom.Name("mongodb-cluster-root-username")
}

func GetRootPasswordOutputName() string {
	return custom.Name("mongodb-cluster-root-password-secret-name")
}

func GetRootPasswordSecretName(mongodbClusterId string) string {
	return fmt.Sprintf("%s-mongodb", mongodbClusterId)
}
