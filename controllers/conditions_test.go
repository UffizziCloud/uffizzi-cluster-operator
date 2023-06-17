package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestBuildReadyTrueCondition(t *testing.T) {
	readyCondition := buildReadyCondition(true)
	if readyCondition.Type != "Ready" {
		t.Errorf("Expected Ready condition type, got %s", readyCondition.Type)
	}
	if readyCondition.Status != metav1.ConditionTrue {
		t.Errorf("Expected Ready condition status true, got %s", readyCondition.Status)
	}
	if readyCondition.Reason != "Ready" {
		t.Errorf("Expected Ready condition reason Ready, got %s", readyCondition.Reason)
	}
	if readyCondition.Message != "UffizziCluster is ready" {
		t.Errorf("Expected Ready condition message UffizziCluster is ready, got %s", readyCondition.Message)
	}
}

func TestBuildReadyFalseCondition(t *testing.T) {
	readyCondition := buildReadyCondition(false)
	if readyCondition.Type != "Ready" {
		t.Errorf("Expected Ready condition type, got %s", readyCondition.Type)
	}
	if readyCondition.Status != metav1.ConditionFalse {
		t.Errorf("Expected Ready condition status false, got %s", readyCondition.Status)
	}
	if readyCondition.Reason != "NotReady" {
		t.Errorf("Expected Ready condition reason NotReady, got %s", readyCondition.Reason)
	}
	if readyCondition.Message != "UffizziCluster is not ready" {
		t.Errorf("Expected Ready condition message UffizziCluster is not ready, got %s", readyCondition.Message)
	}
}
