package gateway

import (
	"fmt"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"

	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/controller"
	ingressnamespace "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//MongodbGatewayIdentifier is used as prefix for naming the gateway resource
	MongodbGatewayIdentifier = "mongodb"
)

func Resources(ctx *pulumi.Context) error {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	gatewayObject := buildGatewayObject(&ctxConfig)
	resourceName := fmt.Sprintf("gateway-%s", gatewayObject.Name)
	manifestPath := filepath.Join(ctxConfig.Spec.WorkspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, gatewayObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}

	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName,
		&pulumik8syaml.ConfigFileArgs{File: manifestPath}, pulumi.Provider(ctxConfig.Spec.KubeProvider))
	if err != nil {
		return errors.Wrap(err, "failed to add gateway manifest")
	}
	return nil
}

/*
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:

	creationTimestamp: "2023-08-23T22:12:03Z"
	generation: 1
	name: prod-planton-live
	namespace: istio-ingress
	resourceVersion: "69782222"
	uid: 69d1b8e1-ad08-4915-8412-be7d6e1a3d18

spec:

	selector:
	  app: istio-ingress
	  istio: ingress
	servers:
	- hosts:
	  - '*.prod.planton.live'
	  - prod.planton.live
	  name: http
	  port:
	    name: http
	    number: 80
	    protocol: HTTP
	  tls:
	    httpsRedirect: true
	- hosts:
	  - '*.prod.planton.live'
	  - prod.planton.live
	  name: https
	  port:
	    name: https
	    number: 443
	    protocol: HTTPS
	  tls:
	    credentialName: cert-prod-planton-live
	    mode: SIMPLE
*/
func buildGatewayObject(ctxConfig *mongodbcontextconfig.ContextConfig) *v1beta1.Gateway {
	return &v1beta1.Gateway{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "Gateway",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      ctxConfig.Spec.ResourceId,
			Namespace: ingressnamespace.Name,
			Labels:    ctxConfig.Spec.Labels,
		},
		Spec: networkingv1beta1.Gateway{
			Selector: controller.SelectorLabels,
			Servers: []*networkingv1beta1.Server{
				{
					Name: MongodbGatewayIdentifier,
					Port: &networkingv1beta1.Port{
						Number:   27017,
						Protocol: "TLS",
						Name:     MongodbGatewayIdentifier,
					},
					Hosts: []string{"*"},
					Tls:   &networkingv1beta1.ServerTLSSettings{},
				},
			},
		},
	}
}
