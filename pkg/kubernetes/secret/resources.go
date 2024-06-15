package secret

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	err := addSecret(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add secret")
	}
	return nil
}

func addSecret(ctx *pulumi.Context) error {
	i := extractInput(ctx)

	// Encode the password in Base64
	base64Password := i.RandomPassword.Result.ApplyT(func(p string) (string, error) {
		return base64.StdEncoding.EncodeToString([]byte(p)), nil
	}).(pulumi.StringOutput)

	// Create or update the secret
	_, err := kubernetescorev1.NewSecret(ctx, i.ResourceName, &kubernetescorev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(i.ResourceName),
			Namespace: pulumi.String(i.NamespaceName),
		},
		Data: pulumi.StringMap{
			MongodbRootPasswordKey: base64Password,
		},
	}, pulumi.Provider(i.KubeProvider), pulumi.Parent(i.Namespace), pulumi.Parent(i.RandomPassword))

	if err != nil {
		return fmt.Errorf("failed to create or update kubernetes secret: %w", err)
	}

	return nil
}
