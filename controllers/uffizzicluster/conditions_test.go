package uffizzicluster

import (
	"testing"
)

func TestBuildInitializingCondition(t *testing.T) {
	initCondition := Initializing()
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

func TestBuildInitializingAPICondition(t *testing.T) {
	initCondition := InitializingAPI()
	if initCondition.Type != "APIReady" {
		t.Errorf("Expected APIReady, got %s", initCondition.Type)
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

func TestBuildInitializingDataStoreCondition(t *testing.T) {
	initCondition := InitializingDataStore()
	if initCondition.Type != "DataStoreReady" {
		t.Errorf("Expected DataStoreReady, got %s", initCondition.Type)
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

func TestBuildAPIReadyCondition(t *testing.T) {
	initCondition := APIReady()
	if initCondition.Type != "APIReady" {
		t.Errorf("Expected APIReady, got %s", initCondition.Type)
	}
	if initCondition.Status != "True" {
		t.Errorf("Expected True, got %s", initCondition.Status)
	}
	if initCondition.Reason != "APIReady" {
		t.Errorf("Expected APIReady, got %s", initCondition.Reason)
	}
	if initCondition.Message != "UffizziCluster API is ready" {
		t.Errorf("Expected UffizziCluster API is ready, got %s", initCondition.Message)
	}
}

func TestDataStoreReadyCondition(t *testing.T) {
	initCondition := DataStoreReady()
	if initCondition.Type != "DataStoreReady" {
		t.Errorf("Expected DataStoreReady, got %s", initCondition.Type)
	}
	if initCondition.Status != "True" {
		t.Errorf("Expected True, got %s", initCondition.Status)
	}
	if initCondition.Reason != "DataStoreReady" {
		t.Errorf("Expected DataStoreReady, got %s", initCondition.Reason)
	}
	if initCondition.Message != "UffizziCluster Datastore is ready" {
		t.Errorf("Expected UffizziCluster Datastore is ready, got %s", initCondition.Message)
	}
}

func TestAPINotReadyCondition(t *testing.T) {
	initCondition := APINotReady()
	if initCondition.Type != "APIReady" {
		t.Errorf("Expected APIReady, got %s", initCondition.Type)
	}
	if initCondition.Status != "False" {
		t.Errorf("Expected False, got %s", initCondition.Status)
	}
	if initCondition.Reason != "APINotReady" {
		t.Errorf("Expected APINotReady, got %s", initCondition.Reason)
	}
	if initCondition.Message != "UffizziCluster API is not ready" {
		t.Errorf("Expected UffizziCluster API is not ready, got %s", initCondition.Message)
	}
}
