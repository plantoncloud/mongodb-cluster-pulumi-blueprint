package cert

import (
	"fmt"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	mongodbnetutilshostname "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/hostname"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/certmanager/clusterissuer"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"path/filepath"

	k8sapimachineryv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Name = "stunnel"
)

func Resources(ctx *pulumi.Context) error {
	if err := addCert(ctx); err != nil {
		return errors.Wrap(err, "failed to add cert")
	}
	return nil
}

func addCert(ctx *pulumi.Context) error {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	certObj := buildCertObject(ctxConfig)
	resourceName := fmt.Sprintf("cert-%s", certObj.Name)
	manifestPath := filepath.Join(ctxConfig.Spec.WorkspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, certObj); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{File: manifestPath}, pulumi.Parent(ctxConfig.Status.AddedResources.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to add cert manifest")
	}
	return nil
}

/*
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:

	name: stunnel
	namespace: planton-pcs-dev-postgres-apr

spec:

	additionalOutputFormats:
	- type: CombinedPEM
	secretName: cert-stunnel
	usages:
	  - server auth
	# At least one of a DNS Name, URI, or IP address is required.
	dnsNames:
	  - <postgres-cluster-id>.dev.planton.cloud
	privateKey:
	  algorithm: ECDSA
	  size: 256
	issuerRef:
	  name: self-signed
	  kind: ClusterIssuer
	  group: cert-manager.io
*/
func buildCertObject(ctxConfig mongodbcontextconfig.ContextConfig) *v1.Certificate {

	internalHostname := mongodbnetutilshostname.GetInternalHostname(ctxConfig.Spec.ResourceId, ctxConfig.Spec.EnvironmentInfo.EnvironmentName, ctxConfig.Spec.EndpointDomainName)
	externalHostname := mongodbnetutilshostname.GetExternalHostname(ctxConfig.Spec.ResourceId, ctxConfig.Spec.EnvironmentInfo.EnvironmentName, ctxConfig.Spec.EndpointDomainName)

	return &v1.Certificate{
		TypeMeta: k8sapimachineryv1.TypeMeta{
			APIVersion: "cert-manager.io/v1",
			Kind:       "Certificate",
		},
		ObjectMeta: k8sapimachineryv1.ObjectMeta{
			Name:      Name,
			Namespace: ctxConfig.Spec.NamespaceName,
			Labels:    ctxConfig.Spec.Labels,
		},
		Spec: v1.CertificateSpec{
			AdditionalOutputFormats: []v1.CertificateAdditionalOutputFormat{
				{
					Type: "CombinedPEM",
				},
			},
			SecretName: GetCertSecretName(Name),
			DNSNames:   []string{internalHostname, externalHostname},
			PrivateKey: &v1.CertificatePrivateKey{
				Algorithm: "ECDSA",
				Size:      256,
			},
			IssuerRef: cmmeta.ObjectReference{
				Kind:  "ClusterIssuer",
				Name:  clusterissuer.SelfSignedIssuerName,
				Group: "cert-manager.io",
			},
		},
	}
}

func GetCertSecretName(certName string) string {
	return fmt.Sprintf("cert-%s", certName)
}
