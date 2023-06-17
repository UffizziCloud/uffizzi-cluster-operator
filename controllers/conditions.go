package controllers

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func buildReadyCondition(ready bool) metav1.Condition {
	if ready {
		return metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "Ready",
			LastTransitionTime: metav1.Now(),
			Message:            "UffizziCluster is ready",
		}
	}
	return metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionFalse,
		Reason:             "NotReady",
		LastTransitionTime: metav1.Now(),
		Message:            "UffizziCluster is not ready",
	}
}
