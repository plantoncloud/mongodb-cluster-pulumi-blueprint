package ingress

import (
	"github.com/pkg/errors"
	mongodbistio "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/istio"
	mongodbloadbalancer "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/loadbalancer"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (newCtx *pulumi.Context, err error) {
	i := extractInput(ctx)
	switch i.IngressType {
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_load_balancer:
		ctx, err = mongodbloadbalancer.Resources(ctx)
		if err != nil {
			return ctx, errors.Wrap(err, "failed to add load balancer resources")
		}
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_ingress_controller:
		if err = mongodbistio.Resources(ctx); err != nil {
			return ctx, errors.Wrap(err, "failed to add istio resources")
		}
	}
	return ctx, nil
}
