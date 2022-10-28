package reconciliation

import (
	"context"

	v1alpha1 "github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	statusUpdaterLogName = "StatusUpdater"

	// Condition Types
	Ready     = "Ready"
	Succeeded = "Succeeded"

	// Condition Reasons
	CustomResourceChanged = "CustomResourceChanged"
	Success               = "Success"
)

// =====================================================================
// Utility struct to update CR status
// =====================================================================
type StatusUpdater struct {
	client client.Client
	log    logr.Logger
}

func NewStatusUpdater(client client.Client, log logr.Logger) *StatusUpdater {
	logger := log.WithName(statusUpdaterLogName)
	return &StatusUpdater{
		client,
		logger,
	}
}

func (su *StatusUpdater) SetReconcileStarted(ctx context.Context, key types.NamespacedName, total int) (*v1alpha1.EntandoAppV2, error) {
	cr, err := su.updateStatus(ctx, key, func(cr *v1alpha1.EntandoAppV2) {
		cr.Status.Progress = 0
		cr.Status.Total = total
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Ready,
			Status:  metav1.ConditionFalse,
			Reason:  CustomResourceChanged,
			Message: "Reconciliation process started",
		})
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Succeeded,
			Status:  metav1.ConditionUnknown,
			Reason:  CustomResourceChanged,
			Message: "Reconciliation process started",
		})
	})

	if err != nil {
		su.log.Error(err, "Unable to update EntandoAppV2's status for reconciliation started")
	}

	return cr, err
}

func (su *StatusUpdater) SetReconcileProcessingComponent(ctx context.Context, key types.NamespacedName, componentName string) (*v1alpha1.EntandoAppV2, error) {
	cr, err := su.updateStatus(ctx, key, func(cr *v1alpha1.EntandoAppV2) {
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Ready,
			Status:  metav1.ConditionFalse,
			Reason:  CustomResourceChanged,
			Message: "Processing " + componentName,
		})
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Succeeded,
			Status:  metav1.ConditionUnknown,
			Reason:  CustomResourceChanged,
			Message: "Processing " + componentName,
		})
	})

	if err != nil {
		su.log.Error(err, "Unable to update EntandoAppV2's status for reconciliation started")
	}

	return cr, err
}

func (su *StatusUpdater) SetReconcileSuccessfullyCompleted(ctx context.Context, key types.NamespacedName) (*v1alpha1.EntandoAppV2, error) {
	cr, err := su.updateStatus(ctx, key, func(cr *v1alpha1.EntandoAppV2) {
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Ready,
			Status:  metav1.ConditionTrue,
			Reason:  Success,
			Message: "Reconciliation process completed",
		})
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Succeeded,
			Status:  metav1.ConditionTrue,
			Reason:  Success,
			Message: "Reconciliation process completed",
		})
	})

	if err != nil {
		su.log.Error(err, "Unable to update EntandoAppV2's status for reconciliation completed")
	}

	return cr, err
}

func (su *StatusUpdater) SetReconcileFailed(ctx context.Context, key types.NamespacedName, reason string) (*v1alpha1.EntandoAppV2, error) {
	cr, err := su.updateStatus(ctx, key, func(cr *v1alpha1.EntandoAppV2) {
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Ready,
			Status:  metav1.ConditionTrue,
			Reason:  reason,
			Message: "Reconciliation process failed",
		})
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    Succeeded,
			Status:  metav1.ConditionFalse,
			Reason:  reason,
			Message: "Reconciliation process failed",
		})
	})

	if err != nil {
		su.log.Error(err, "Unable to update EntandoAppV2's status for reconciliation failed")
	}

	return cr, err
}

func (su *StatusUpdater) IncrementProgress(ctx context.Context, key types.NamespacedName) (*v1alpha1.EntandoAppV2, error) {
	cr, err := su.updateStatus(ctx, key, func(cr *v1alpha1.EntandoAppV2) {
		su.log.Info("Updating progress", "progress", cr.Status.Progress+1)
		cr.Status.Progress = cr.Status.Progress + 1
	})

	if err != nil {
		su.log.Error(err, "Unable to update EntandoAppV2's progress status", "progress", cr.Status.Progress)
	}

	return cr, err
}

func (su *StatusUpdater) updateStatus(ctx context.Context, key types.NamespacedName, updateFields func(cr *v1alpha1.EntandoAppV2)) (*v1alpha1.EntandoAppV2, error) {
	cr := &v1alpha1.EntandoAppV2{}

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := su.client.Get(ctx, key, cr); err == nil {
			updateFields(cr)
			cr.Status.ObservedGeneration = cr.ObjectMeta.Generation
			return su.client.Status().Update(ctx, cr)
		} else {
			return err
		}
	})

	return cr, err
}
