package vcluster

import v1 "k8s.io/api/core/v1"

var (
	GvisorToleration = Toleration{
		Key:      "sandbox.gke.io/runtime",
		Operator: v1.TolerationOpEqual,
		Value:    "gvisor",
		Effect:   v1.TaintEffectNoSchedule,
	}
	GvisorNodeSelector = map[string]string{
		"key":   "sandbox.gke.io/runtime",
		"value": "gvisor",
	}
)

type Toleration v1.Toleration

func (t Toleration) Notation() string {
	return t.Key + "=" + t.Value + ":" + string(t.Effect)
}
