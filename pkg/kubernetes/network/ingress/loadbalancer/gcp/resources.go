package gcp

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbloadbalancercommon "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/loadbalancer/common"
	pulumikubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	// Create a Kubernetes Service of type LoadBalancer
	externalLoadBalancerService, err := addExternal(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add external load balancer")
	}
	internalLoadBalancerService, err := addInternal(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add internal load balancer")
	}

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	addLoadBalancerExternalServiceToContext(&ctxConfig, externalLoadBalancerService)
	addLoadBalancerInternalServiceToContext(&ctxConfig, internalLoadBalancerService)
	ctx = ctx.WithValue(mongodbcontextconfig.Key, ctxConfig)

	return ctx, nil
}

func addExternal(ctx *pulumi.Context) (*pulumikubernetescorev1.Service, error) {
	i := extractInput(ctx)
	addedKubeService, err := pulumikubernetescorev1.NewService(ctx,
		mongodbloadbalancercommon.ExternalLoadBalancerServiceName,
		getLoadBalancerServiceArgs(i, mongodbloadbalancercommon.ExternalLoadBalancerServiceName, i.ExternalEndpoint),
		pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "30s", Update: "30s", Delete: "30s"}), pulumi.Parent(i.Namespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes service of type load balancer")
	}
	return addedKubeService, nil
}

func addInternal(ctx *pulumi.Context) (*pulumikubernetescorev1.Service, error) {
	i := extractInput(ctx)
	addedKubeService, err := pulumikubernetescorev1.NewService(ctx,
		mongodbloadbalancercommon.InternalLoadBalancerServiceName,
		getInternalLoadBalancerServiceArgs(i, i.InternalEndpoint, i.Namespace),
		pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "30s", Update: "30s", Delete: "30s"}), pulumi.Parent(i.Namespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes service of type load balancer")
	}
	return addedKubeService, nil
}

func getInternalLoadBalancerServiceArgs(i *input, hostname string, namespace *pulumikubernetescorev1.Namespace) *pulumikubernetescorev1.ServiceArgs {
	resp := getLoadBalancerServiceArgs(i, mongodbloadbalancercommon.InternalLoadBalancerServiceName, hostname)
	resp.Metadata = &v1.ObjectMetaArgs{
		Name:      pulumi.String(mongodbloadbalancercommon.InternalLoadBalancerServiceName),
		Namespace: namespace.Metadata.Name(),
		Labels:    namespace.Metadata.Labels(),
		Annotations: pulumi.StringMap{
			"cloud.google.com/load-balancer-type":       pulumi.String("Internal"),
			"planton.cloud/endpoint-domain-name":        pulumi.String(i.EndpointDomainName),
			"external-dns.alpha.kubernetes.io/hostname": pulumi.String(hostname),
		},
	}
	return resp
}

func getLoadBalancerServiceArgs(i *input, serviceName string, hostname string) *pulumikubernetescorev1.ServiceArgs {
	return &pulumikubernetescorev1.ServiceArgs{
		Metadata: &v1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: i.Namespace.Metadata.Name().Elem(),
			Labels:    i.Namespace.Metadata.Labels(),
			Annotations: pulumi.StringMap{
				"planton.cloud/endpoint-domain-name":        pulumi.String(i.EndpointDomainName),
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
				"app.kubernetes.io/instance":  pulumi.String(i.ResourceId),
				"app.kubernetes.io/name":      pulumi.String(i.ResourceName),
			},
		},
	}
}

func addLoadBalancerExternalServiceToContext(existingConfig *mongodbcontextconfig.ContextConfig, loadBalancerService *pulumikubernetescorev1.Service) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &mongodbcontextconfig.AddedResources{
			LoadBalancerExternalService: loadBalancerService,
		}
		return
	}
	existingConfig.Status.AddedResources.LoadBalancerExternalService = loadBalancerService
}

func addLoadBalancerInternalServiceToContext(existingConfig *mongodbcontextconfig.ContextConfig, loadBalancerService *pulumikubernetescorev1.Service) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &mongodbcontextconfig.AddedResources{
			LoadBalancerInternalService: loadBalancerService,
		}
		return
	}
	existingConfig.Status.AddedResources.LoadBalancerInternalService = loadBalancerService
}
