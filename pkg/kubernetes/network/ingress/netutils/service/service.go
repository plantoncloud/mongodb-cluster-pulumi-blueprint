package service

import (
	"fmt"
	"github.com/plantoncloud-inc/go-commons/kubernetes/network/dns"
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	ctx.Export(ctxConfig.Status.OutputKeyNames.KubeServiceName, pulumi.String(ctxConfig.Status.OutputValues.KubeServiceName))
	return nil
}

func GetKubeServiceNameFqdn(mongodbClusterName, namespace string) string {
	return fmt.Sprintf("%s.%s.%s", GetKubeServiceName(mongodbClusterName), namespace, dns.DefaultDomain)
}

func GetKubeServiceName(mongodbClusterName string) string {
	return fmt.Sprintf(mongodbClusterName)
}
