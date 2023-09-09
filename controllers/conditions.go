package controllers

import (
	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Condition types.
const (
	// TypeReady resources are believed to be ready to handle work.
	TypeReady = "Ready"
	TypeSleep = "Sleep"
)

// Reasons a resource is or is not ready.
const (
	ReasonInitializing = "Initializing"
	ReasonSleeping     = "Sleeping"
	ReasonAwoken       = "Awoken"
)

func Initializing() metav1.Condition {
	return metav1.Condition{
		Type:               TypeReady,
		Status:             metav1.ConditionUnknown,
		Reason:             ReasonInitializing,
		LastTransitionTime: metav1.Now(),
		Message:            "UffizziCluster is being initialized",
	}
}

func Sleeping(time metav1.Time) metav1.Condition {
	return metav1.Condition{
		Type:               TypeSleep,
		Status:             metav1.ConditionTrue,
		Reason:             ReasonSleeping,
		LastTransitionTime: time,
		Message:            "UffizziCluster put to sleep manually",
	}
}

func Awoken(time metav1.Time) metav1.Condition {
	return metav1.Condition{
		Type:               TypeSleep,
		Status:             metav1.ConditionFalse,
		Reason:             ReasonAwoken,
		LastTransitionTime: time,
		Message:            "UffizziCluster awoken manually",
	}
}

func mirrorHelmReleaseConditions(helmRelease *fluxhelmv2beta1.HelmRelease, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) {
	uClusterConditions := []metav1.Condition{}
	for _, c := range helmRelease.Status.Conditions {
		helmMessage := "[HelmRelease] " + c.Message
		uClusterCondition := c
		uClusterCondition.Message = helmMessage
		uClusterConditions = append(uClusterConditions, uClusterCondition)
	}
	setConditions(uCluster, uClusterConditions...)
}

// setConditions sets the supplied conditions, replacing any existing conditions
// of the same type. This is a no-op if all supplied conditions are identical,
// ignoring the last transition time, to those already set.
func setConditions(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster, c ...metav1.Condition) {
	for _, new := range c {
		exists := false
		for i, existing := range uCluster.Status.Conditions {
			if existing.Type != new.Type {
				continue
			}
			if conditionsEqual(existing, new) {
				exists = true
				continue
			}
			uCluster.Status.Conditions[i] = new
			exists = true
		}
		if !exists {
			uCluster.Status.Conditions = append(uCluster.Status.Conditions, new)
		}
	}
}

func setCondition(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster, c metav1.Condition) {
	setConditions(uCluster, c)
}

// Equal returns true if the condition is identical to the supplied condition,
// ignoring the LastTransitionTime.
//
//nolint:gocritic // just a few bytes too heavy
func conditionsEqual(c, other metav1.Condition) bool {
	return c.Type == other.Type &&
		c.Status == other.Status &&
		c.Reason == other.Reason &&
		c.Message == other.Message
}
