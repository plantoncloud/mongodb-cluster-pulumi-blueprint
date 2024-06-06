// Package service adds a kubernetes service which is required to forward traffic from istio pods to stunnel sidecar containers running alongside postgres pods.
package stunnel

import (
	"github.com/pkg/errors"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/mongodbcluster"
	mongodbnetutilport "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/port"
	pulumikubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	pulumikubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	StunnelServiceName = "stunnel"
)

func Resources(ctx *pulumi.Context) error {
	if _, err := addService(ctx); err != nil {
		return errors.Wrap(err, "failed to add stunnel service")
	}
	return nil
}

/*
apiVersion: v1
kind: Service
metadata:

	name: stunnel
	namespace: planton-pcs-dev-postgres-apr

spec:

	type: ClusterIP
	ports:
	- name: postgresql
	  port: 5432
	  protocol: TCP
	  targetPort: 15432
	selector:
	  application: spilo
	  cluster-name: pcs-apr
	  team: pcs
*/
func addService(ctx *pulumi.Context) (*pulumikubernetescorev1.Service, error) {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	svc, err := pulumikubernetescorev1.NewService(ctx, StunnelServiceName, &pulumikubernetescorev1.ServiceArgs{
		Metadata: pulumikubernetesmetav1.ObjectMetaArgs{
			Name:      pulumi.String(StunnelServiceName),
			Namespace: ctxConfig.Status.AddedResources.Namespace.Metadata.Name(),
		},
		Spec: &pulumikubernetescorev1.ServiceSpecArgs{
			Type: pulumi.String("ClusterIP"),
			Selector: pulumi.StringMap{
				"app.kubernetes.io/component": pulumi.String("mongodb"),
				"app.kubernetes.io/instance":  pulumi.String(ctxConfig.Spec.ResourceId),
				"app.kubernetes.io/name":      pulumi.String(ctxConfig.Spec.ResourceName),
			},
			Ports: pulumikubernetescorev1.ServicePortArray{
				&pulumikubernetescorev1.ServicePortArgs{
					Name:       pulumi.String("mongodb"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(mongodbnetutilport.MongoDbPort),
					TargetPort: pulumi.Int(mongodbcluster.StunnelContainerPort),
				},
			},
		},
	}, pulumi.Parent(ctxConfig.Status.AddedResources.Namespace))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add service")
	}
	return svc, nil
}
