package contextconfig

import (
	code2cloudenvironmentmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/environment/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
)

const (
	Key = "ctx-config"
)

type ContextConfig struct {
	Spec   *ContextConfigSpec
	Status *ContextConfigStatus
}

type ContextConfigSpec struct {
	KubeProvider       *kubernetes.Provider
	ResourceId         string
	ResourceName       string
	Labels             map[string]string
	WorkspaceDir       string
	NamespaceName      string
	EnvironmentInfo    *code2cloudenvironmentmodel.ApiResourceEnvironmentInfo
	IsIngressEnabled   bool
	IngressType        kubernetesworkloadingresstype.KubernetesWorkloadIngressType
	EndpointDomainName string
	EnvDomainName      string
}

type ContextConfigStatus struct {
	AddedResources *AddedResources
	OutputKeyNames *OutputKeyNames
	OutputValues   *OutputValues
}

type AddedResources struct {
	Namespace *kubernetescorev1.Namespace
}

type OutputKeyNames struct {
	IngressEndpoint               string
	RootPasswordSecret            string
	RootUsername                  string
	KubeServiceName               string
	KubeEndpoint                  string
	KubeForwardCommand            string
	LoadBalancerInternalIpAddress string
	LoadBalancerExternalIpAddress string
	Namespace                     string
}

type OutputValues struct {
	IngressEndpoint    string
	RootPasswordSecret string
	RootUsername       string
	KubeServiceName    string
	KubeEndpoint       string
	KubeForwardCommand string
	Namespace          string
}
