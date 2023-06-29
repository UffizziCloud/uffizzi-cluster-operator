package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildInitializingCondition() metav1.Condition {
	return metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionUnknown,
		Reason:             "Initializing",
		LastTransitionTime: metav1.Now(),
		Message:            "UffizziCluster is being initialized",
	}
}
