package controllers

import (
	"testing"
)

func TestBuildInitializingCondition(t *testing.T) {
	initCondition := buildInitializingCondition()
	if initCondition.Type != "Ready" {
		t.Errorf("Expected Ready, got %s", initCondition.Type)
	}
	if initCondition.Status != "Unknown" {
		t.Errorf("Expected Unknown, got %s", initCondition.Status)
	}
	if initCondition.Reason != "Initializing" {
		t.Errorf("Expected Initializing, got %s", initCondition.Reason)
	}
	if initCondition.Message != "UffizziCluster is being initialized" {
		t.Errorf("Expected UffizziCluster is being initialized, got %s", initCondition.Message)
	}
}
