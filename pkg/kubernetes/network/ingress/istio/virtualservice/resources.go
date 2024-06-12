package virtualservice

import (
	"fmt"
	"github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/outputs"

	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	ingressnamespace "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Resources(ctx *pulumi.Context) error {
	if err := addVirtualService(ctx); err != nil {
		return errors.Wrap(err, "failed to add external virtual service")
	}
	return nil
}

func addVirtualService(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	var virtualServiceObject = buildVirtualServiceObject(i)
	resourceName := fmt.Sprintf("virtual-service-%s", virtualServiceObject.Name)
	manifestPath := filepath.Join(i.WorkspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, virtualServiceObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{
		File: manifestPath,
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "30s", Update: "30s", Delete: "30s"}), pulumi.DependsOn([]pulumi.Resource{i.Namespace}), pulumi.Parent(i.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to add virtual-service manifest")
	}
	return nil
}

/*
apiVersion: v1
items:
- apiVersion: networking.istio.io/v1beta1
  kind: VirtualService
  metadata:
    creationTimestamp: "2024-05-30T10:40:39Z"
    generation: 3
    name: ingress-controller-test
    namespace: jnk-planton-cloud-prod-ingress-controller-test
    resourceVersion: "276489649"
    uid: 662b159f-54f9-4227-970f-22d8f033ac5b
  spec:
    gateways:
    - istio-ingress/prod-planton-live
    hosts:
    - jnk-planton-cloud-prod-ingress-controller-test.prod.planton.live
    http:
    - name: jnk-planton-cloud-prod-ingress-controller-test
      route:
      - destination:
          host: ingress-controller-test.jnk-planton-cloud-prod-ingress-controller-test.svc.cluster.local
          port:
            number: 8080
kind: List
metadata:
  resourceVersion: ""
*/

func buildVirtualServiceObject(i *input) *v1beta1.VirtualService {
	return &v1beta1.VirtualService{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "VirtualService",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      i.KubeServiceName,
			Namespace: i.NameSpaceName,
		},
		Spec: networkingv1beta1.VirtualService{
			Gateways: []string{fmt.Sprintf("%s/%s", ingressnamespace.Name, i.ResourceId)},
			Hosts:    i.HostNames,
			Tcp: []*networkingv1beta1.TCPRoute{{
				Match: []*networkingv1beta1.L4MatchAttributes{
					{
						Port: outputs.MongoDbPort,
					},
				},
				Route: []*networkingv1beta1.RouteDestination{
					{
						Destination: &networkingv1beta1.Destination{
							Host: i.KubeLocalEndpoint,
							Port: &networkingv1beta1.PortSelector{
								Number: outputs.MongoDbPort,
							},
						},
					},
				},
			}},
		},
	}
}
