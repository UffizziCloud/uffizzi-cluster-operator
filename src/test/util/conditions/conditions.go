package conditions

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Checks if the required conditions are present and match in the actual conditions slice.
// Both requiredConditions and actualConditions are slices of metav1.Condition.
func ContainsAllConditions(requiredConditions, actualConditions []metav1.Condition) bool {
	for _, requiredCondition := range requiredConditions {
		found := false
		for _, actualCondition := range actualConditions {
			if actualCondition.Type == requiredCondition.Type &&
				actualCondition.Status == requiredCondition.Status {
				// Add more condition checks here if necessary (e.g., Reason, Message)
				found = true
				d := cmp.Diff(requiredConditions, actualConditions)
				ginkgo.GinkgoWriter.Printf(diff.PrintWantGot(d))
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
