package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	apiv1 "github.com/rbales79/r-operator/api/v1"
)

type ArgoPipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ArgoPipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pipeline apiv1.ArgoPipeline
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Parse pipelineYaml from the CR
	manifest := []byte(pipeline.Spec.PipelineYaml)
	// Try to decode as a runtime.Object
	obj, gvk, err := r.Scheme.Codecs.UniversalDeserializer().Decode(manifest, nil, nil)
	if err != nil {
		// Optionally update status or log error
		return ctrl.Result{}, fmt.Errorf("failed to decode pipeline YAML: %w", err)
	}

	// Set namespace if not specified in YAML
	accessor, _ := runtime.UnstructuredAccessor(obj)
	if accessor.GetNamespace() == "" {
		accessor.SetNamespace(pipeline.Namespace)
	}

	// Create or update the Argo Workflow/CronWorkflow
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		objKey := client.ObjectKey{Name: accessor.GetName(), Namespace: accessor.GetNamespace()}
		existing := &runtime.Unknown{}
		if getErr := r.Get(ctx, objKey, existing); getErr != nil {
			// If not found, create
			return r.Create(ctx, obj)
		}
		// If exists, update
		accessor.SetResourceVersion(existing.GetResourceVersion())
		return r.Update(ctx, obj)
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create/update Argo resource: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *ArgoPipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.ArgoPipeline{}).
		Complete(r)
}