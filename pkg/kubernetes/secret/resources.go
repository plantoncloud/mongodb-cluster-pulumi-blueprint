package secret

import (
	"encoding/base64"
	"fmt"
	"github.com/google/martian/log"
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"k8s.io/apimachinery/pkg/util/rand"
	"time"
)

func Resources(ctx *pulumi.Context) error {
	err := addSecret(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add secret")
	}
	return nil
}

func addSecret(ctx *pulumi.Context) error {
	var i = extractInput(ctx)

	rootPassword, err := getRootPasswordValue(ctx)
	if err != nil {
		return fmt.Errorf("failed to get root password value: %w", err)
	}
	// Create or Update the secret using the provider
	_, err = kubernetescorev1.NewSecret(ctx, i.ResourceName, &kubernetescorev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(i.ResourceName),
			Namespace: pulumi.String(i.NamespaceName),
		},
		Data: pulumi.StringMap{
			MongodbRootPasswordKey: pulumi.String(rootPassword),
		},
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(i.KubeProvider), pulumi.Parent(i.Namespace))
	if err != nil {
		return fmt.Errorf("failed to create or update kubernetes secret: %w", err)
	}

	return nil
}

func getRootPasswordValue(ctx *pulumi.Context) (string, error) {
	var i = extractInput(ctx)

	// Attempt to get the secret
	secret, err := kubernetescorev1.GetSecret(ctx, fmt.Sprintf("read-%s", i.ResourceName), pulumi.ID(fmt.Sprintf("%s/%s", i.NamespaceName, i.ResourceName)),
		nil, pulumi.Provider(i.KubeProvider), pulumi.Parent(i.Namespace))
	if err != nil {
		log.Debugf("secret not found in namespace %s", i.NamespaceName)
	}
	if secret == nil {
		return generateRandomString(16), nil
	}

	myKeyBase64 := secret.Data.ApplyT(func(data map[string]string) (string, error) {
		// base64 decode the key value
		password, exists := data[MongodbRootPasswordKey]
		if !exists {
			return generateRandomString(16), nil
		}
		return password, nil
	})
	return myKeyBase64.ElementType().String(), nil
}

// generateRandomString generates a random string of a given length
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return base64.StdEncoding.EncodeToString(b)
}
