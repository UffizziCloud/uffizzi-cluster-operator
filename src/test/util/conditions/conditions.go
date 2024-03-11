package conditions

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/diff"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Checks if the required conditions are present and match in the actual conditions slice.
// Both requiredConditions and actualConditions are slices of metav1.Condition.
func ContainsConditionsDecorator(requiredConditions, actualConditions []metav1.Condition, flip bool) bool {
	var foundQ = func(og bool) bool {
		if flip {
			return !og
		}
		return og
	}
	var didFindQ = func() bool {
		return foundQ(true)
	}
	var notFoundQ = func() bool {
		return !foundQ(true)
	}
	for _, requiredCondition := range requiredConditions {
		found := notFoundQ()
		for _, actualCondition := range actualConditions {
			if actualCondition.Type == requiredCondition.Type &&
				actualCondition.Status == requiredCondition.Status {
				// Add more condition checks here if necessary (e.g., Reason, Message)
				found = didFindQ()
				break
			}
		}
		if !found {
			return notFoundQ()
		}
	}
	return didFindQ()
}

// Return true if zero condtions match between required and actual conditions
func ContainsAllConditions(requiredConditions, actualConditions []metav1.Condition) bool {
	return ContainsConditionsDecorator(requiredConditions, actualConditions, true)
}

// Return true if zero condtions match between required and actual conditions
func ContainsNoConditions(requiredConditions, actualConditions []metav1.Condition) bool {
	return ContainsConditionsDecorator(requiredConditions, actualConditions, false)
}

func CreateConditionsCmpDiff(requiredConditions, actualConditions []metav1.Condition) string {
	// create a clone of actual conditions which only has the keys required conditions has
	// then compare the two slices
	actualConditionsForCmp := []metav1.Condition{}
	requiredConditionsForCmp := []metav1.Condition{}
	for _, actualCondition := range actualConditions {
		actualConditionForCmp := metav1.Condition{
			Type:   actualCondition.Type,
			Status: actualCondition.Status,
			Reason: actualCondition.Reason,
		}
		actualConditionsForCmp = append(actualConditionsForCmp, actualConditionForCmp)
	}
	for _, requiredCondition := range requiredConditions {
		requiredConditionForCmp := metav1.Condition{
			Type:   requiredCondition.Type,
			Status: requiredCondition.Status,
			Reason: requiredCondition.Reason,
		}
		requiredConditionsForCmp = append(requiredConditionsForCmp, requiredConditionForCmp)
	}

	return diff.PrintWantGot(cmp.Diff(requiredConditionsForCmp, actualConditionsForCmp))
}
