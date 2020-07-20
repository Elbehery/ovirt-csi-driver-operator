package operator

import (
	"context"
	"fmt"

	dynamicclient "k8s.io/client-go/dynamic"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/operator/csi/csicontrollerset"
	goc "github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/ovirt/csi-driver-operator/pkg/apis/operator/v1alpha1"
	"github.com/ovirt/csi-driver-operator/pkg/generated"
)

const (
	// Operand and operator run in the same namespace
	defaultNamespace = "openshift-cluster-csi-drivers"
	operatorName     = "ovirt-csi-driver-operator"
	operandName      = "ovirt-csi-driver"
)

func RunOperator(ctx context.Context, controllerConfig *controllercmd.ControllerContext) error {
	// Create clientsets and informers
	kubeClient := kubeclient.NewForConfigOrDie(rest.AddUserAgent(controllerConfig.KubeConfig, operatorName))
	dynamicClient := dynamicclient.NewForConfigOrDie(rest.AddUserAgent(controllerConfig.KubeConfig, operatorName))
	kubeInformersForNamespaces := v1helpers.NewKubeInformersForNamespaces(kubeClient, defaultNamespace, "")

	// Create GenericOperatorclient. This is used by the library-go controllers created down below
	gvr := v1alpha1.SchemeGroupVersion.WithResource("ovirtcsidrivers")
	operatorClient, dynamicInformers, err := goc.NewClusterScopedOperatorClient(controllerConfig.KubeConfig, gvr)
	if err != nil {
		return err
	}

	csiControllerSet := csicontrollerset.NewCSIControllerSet(
		operatorClient,
		controllerConfig.EventRecorder,
	).WithLogLevelController().WithManagementStateController(
		operandName,
		false,
	).WithStaticResourcesController(
		"OvirtDriverStaticResources",
		kubeClient,
		kubeInformersForNamespaces,
		generated.Asset,
		[]string{
			"storageclass.yaml",
			"csidriver.yaml",
			"controller_sa.yaml",
			"node_sa.yaml",
			"rbac/attacher_binding.yaml",
			"rbac/attacher_role.yaml",
			"rbac/controller_privileged_binding.yaml",
			"rbac/node_privileged_binding.yaml",
			"rbac/privileged_role.yaml",
			"rbac/provisioner_binding.yaml",
			"rbac/provisioner_role.yaml",
			"rbac/resizer_binding.yaml",
			"rbac/resizer_role.yaml",
			"rbac/snapshotter_binding.yaml",
			"rbac/snapshotter_role.yaml",
		},
	).
		WithCredentialsRequestController(
			"OvirtDriverCredentialsRequestController",
			defaultNamespace,
			generated.MustAsset,
			"credentials.yaml",
			dynamicClient,
		).
		WithCSIDriverController(
			"OvirtDriverController",
			operandName,
			defaultNamespace,
			generated.MustAsset,
			kubeClient,
			kubeInformersForNamespaces.InformersFor(defaultNamespace),
			csicontrollerset.WithControllerService("controller.yaml"),
			csicontrollerset.WithNodeService("node.yaml"),
		)

	if err != nil {
		return err
	}

	klog.Info("Starting the informers")
	go kubeInformersForNamespaces.Start(ctx.Done())
	go dynamicInformers.Start(ctx.Done())

	klog.Info("Starting controllerset")
	go csiControllerSet.Run(ctx, 1)

	<-ctx.Done()

	return fmt.Errorf("stopped")
}
