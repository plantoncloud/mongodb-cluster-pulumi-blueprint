package mongodbcluster

import (
	mongodbcontextconfig "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/contextconfig"
	"github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/istio/cert"
	mongodbnetutilsport "github.com/plantoncloud/mongodb-cluster-pulumi-blueprint/pkg/kubernetes/network/ingress/netutils/port"
	plantoncloudkubeclustermodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	plantoncloudmongodbmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/mongodbcluster/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/containerresources"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/helm/mergemaps"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
	"strconv"
)

const (
	EnvVarStunnelSidecarImage = "STUNNEL_SIDECAR_IMAGE"
	StunnelContainerPort      = 17017
	StunnelCertMountPath      = "/server/ca.pem"
)

func getHelmValues(ctx *pulumi.Context, containerSpec *plantoncloudmongodbmodel.MongodbClusterSpecKubernetesSpecMongodbContainerSpec,
	customValues map[string]string) pulumi.Map {

	var ctxConfig = ctx.Value(mongodbcontextconfig.Key).(mongodbcontextconfig.ContextConfig)

	// https://github.com/bitnami/charts/blob/main/bitnami/mongodb/values.yaml
	var baseValues = pulumi.Map{
		"nameOverride":     pulumi.String(ctxConfig.Spec.ResourceName),
		"fullnameOverride": pulumi.String(ctxConfig.Spec.ResourceName),
		"resources":        containerresources.ConvertToPulumiMap(containerSpec.Resources),
		// todo: hard-coding this to 1 since we are only using `standalone` architecture,
		// need to revisit this to handle `replicaSet` architecture
		"replicaCount": pulumi.Int(1),
		"persistence": pulumi.Map{
			"enabled": pulumi.Bool(containerSpec.IsPersistenceEnabled),
			"size":    pulumi.String(containerSpec.DiskSize),
		},
		"podLabels":         pulumi.ToStringMap(ctxConfig.Spec.Labels),
		"commonLabels":      pulumi.ToStringMap(ctxConfig.Spec.Labels),
		"useStatefulSet":    pulumi.Bool(true),
		"sidecars":          pulumi.Array{getStunnelSidecar()},
		"extraVolumeMounts": pulumi.Array{getAdditionalVolumes()},
	}
	mergemaps.MergeMapToPulumiMap(baseValues, customValues)
	return baseValues
}

func getStunnelSidecarResources() *plantoncloudkubeclustermodel.ContainerResources {
	return &plantoncloudkubeclustermodel.ContainerResources{
		Requests: &plantoncloudkubeclustermodel.CpuMemory{
			Cpu:    "100m",
			Memory: "100Mi",
		},
		Limits: &plantoncloudkubeclustermodel.CpuMemory{
			Cpu:    "500m",
			Memory: "1Gi",
		},
	}
}

func getStunnelSidecar() pulumi.Map {
	stunnelSidecarImage, _ := os.LookupEnv(EnvVarStunnelSidecarImage)
	return pulumi.Map{
		"name":      pulumi.String(englishword.EnglishWord_stunnel),
		"image":     pulumi.String(stunnelSidecarImage),
		"resources": containerresources.ConvertToPulumiMap(getStunnelSidecarResources()),
		"ports": pulumi.Array{
			pulumi.Map{
				"name":          pulumi.String("mongodb-stunnel"),
				"containerPort": pulumi.Int(StunnelContainerPort),
				"protocol":      pulumi.String("TCP"),
			},
		},
		"env": pulumi.Array{
			pulumi.Map{
				"name":  pulumi.String("STUNNEL_MODE"),
				"value": pulumi.String(englishword.EnglishWord_server.String()),
			},
			pulumi.Map{
				"name":  pulumi.String("STUNNEL_LOG_LEVEL"),
				"value": pulumi.String(englishword.EnglishWord_debug.String()),
			},
			pulumi.Map{
				"name":  pulumi.String("STUNNEL_ACCEPT_PORT"),
				"value": pulumi.String(strconv.Itoa(StunnelContainerPort)),
			},
			pulumi.Map{
				"name":  pulumi.String("STUNNEL_FORWARD_HOST"),
				"value": pulumi.String(englishword.EnglishWord_localhost.String()),
			},
			pulumi.Map{
				"name":  pulumi.String("STUNNEL_FORWARD_PORT"),
				"value": pulumi.String(strconv.Itoa(mongodbnetutilsport.MongoDbPort)),
			},
		},
	}
}

func getAdditionalVolumes() pulumi.Map {
	return pulumi.Map{
		"name":             pulumi.String("stunnel-ca"),
		"mountPath":        pulumi.String(StunnelCertMountPath),
		"subPath":          pulumi.String("tls-combined.pem"),
		"targetContainers": pulumi.Array{pulumi.String(englishword.EnglishWord_stunnel)},
		"volumeSource": pulumi.Map{
			"secret": pulumi.Map{
				"secretName": pulumi.String(cert.GetCertSecretName(cert.Name)),
			},
		},
	}
}
