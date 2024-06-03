package gcp

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbloadbalancercommon "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/loadbalancer/common"
	mongodbnetutilshostname "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/hostname"
	pulumicommonsloadbalancerservice "github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/loadbalancer/service"
	pulumikubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	// Create a Kubernetes Service of type LoadBalancer
	if err := addExternal(ctx); err != nil {
		return errors.Wrap(err, "failed to add external load balancer")
	}
	if err := addInternal(ctx); err != nil {
		return errors.Wrap(err, "failed to add internal load balancer")
	}
	return nil
}

func addExternal(ctx *pulumi.Context) error {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	var resourceId = ctxConfig.Spec.ResourceId
	var nameSpace = ctxConfig.Status.AddedResources.Namespace

	hostname := mongodbnetutilshostname.GetExternalHostname(resourceId, ctxConfig.Spec.EnvironmentInfo.EnvironmentName, ctxConfig.Spec.EndpointDomainName)
	addedKubeService, err := pulumikubernetescorev1.NewService(ctx,
		mongodbloadbalancercommon.ExternalLoadBalancerServiceName,
		getLoadBalancerServiceArgs(ctxConfig, mongodbloadbalancercommon.ExternalLoadBalancerServiceName, hostname, nameSpace), pulumi.Parent(nameSpace))
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes service of type load balancer")
	}

	var ipAddress = pulumicommonsloadbalancerservice.GetIpAddress(addedKubeService)
	ctx.Export(ctxConfig.Status.OutputKeyNames.LoadBalancerExternalIpAddress, pulumi.String(ipAddress))
	return nil
}

func addInternal(ctx *pulumi.Context) error {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	var resourceId = ctxConfig.Spec.ResourceId
	var nameSpace = ctxConfig.Status.AddedResources.Namespace

	hostname := mongodbnetutilshostname.GetInternalHostname(resourceId, ctxConfig.Spec.EnvironmentInfo.EnvironmentName, ctxConfig.Spec.EndpointDomainName)
	addedKubeService, err := pulumikubernetescorev1.NewService(ctx,
		mongodbloadbalancercommon.InternalLoadBalancerServiceName,
		getInternalLoadBalancerServiceArgs(ctxConfig, hostname, nameSpace), pulumi.Parent(nameSpace))
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes service of type load balancer")
	}

	var ipAddress = pulumicommonsloadbalancerservice.GetIpAddress(addedKubeService)
	ctx.Export(ctxConfig.Status.OutputKeyNames.LoadBalancerInternalIpAddress, pulumi.String(ipAddress))
	return nil
}

func getInternalLoadBalancerServiceArgs(ctxConfig mongodbcontextconfig.ContextConfig, hostname string, namespace *pulumikubernetescorev1.Namespace) *pulumikubernetescorev1.ServiceArgs {
	resp := getLoadBalancerServiceArgs(ctxConfig, mongodbloadbalancercommon.InternalLoadBalancerServiceName, hostname, namespace)
	resp.Metadata = &v1.ObjectMetaArgs{
		Name:      pulumi.String(mongodbloadbalancercommon.InternalLoadBalancerServiceName),
		Namespace: namespace.Metadata.Name(),
		Labels:    namespace.Metadata.Labels(),
		Annotations: pulumi.StringMap{
			"cloud.google.com/load-balancer-type":       pulumi.String("Internal"),
			"planton.cloud/endpoint-domain-name":        pulumi.String(ctxConfig.Spec.EndpointDomainName),
			"external-dns.alpha.kubernetes.io/hostname": pulumi.String(hostname),
		},
	}
	return resp
}

func getLoadBalancerServiceArgs(ctxConfig mongodbcontextconfig.ContextConfig, serviceName, hostname string, namespace *pulumikubernetescorev1.Namespace) *pulumikubernetescorev1.ServiceArgs {
	return &pulumikubernetescorev1.ServiceArgs{
		Metadata: &v1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: namespace.Metadata.Name(),
			Labels:    namespace.Metadata.Labels(),
			Annotations: pulumi.StringMap{
				"planton.cloud/endpoint-domain-name":        pulumi.String(ctxConfig.Spec.EndpointDomainName),
				"external-dns.alpha.kubernetes.io/hostname": pulumi.String(hostname)}},
		Spec: &pulumikubernetescorev1.ServiceSpecArgs{
			Type: pulumi.String("LoadBalancer"), // Service type is LoadBalancer
			Ports: pulumikubernetescorev1.ServicePortArray{
				&pulumikubernetescorev1.ServicePortArgs{
					Name:       pulumi.String("mongodb"),
					Port:       pulumi.Int(27017),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.String("mongodb"), // This assumes your Mongodb pod has a port named 'mongodb'
				},
			},
			Selector: pulumi.StringMap{
				"app.kubernetes.io/component": pulumi.String("mongodb"),
				"app.kubernetes.io/instance":  namespace.Metadata.Name().Elem(),
				"app.kubernetes.io/name":      pulumi.String("mongodb"),
			},
		},
	}
}
