package vcluster

import v1 "k8s.io/api/core/v1"

type Toleration v1.Toleration
type NodeSelector map[string]string

func (t Toleration) Notation() string {
	if t.Value == "" {
		return t.Key + ":" + string(t.Effect)
	}
	return t.Key + "=" + t.Value + ":" + string(t.Effect)
}

func (t Toleration) ToV1() v1.Toleration {
	return v1.Toleration(t)
}

func (n NodeSelector) Notation() string {
	for k, v := range n {
		return k + "=" + v
	}
	return ""
}

var (
	GvisorToleration = Toleration{
		Key:      "sandbox.gke.io/runtime",
		Operator: v1.TolerationOpEqual,
		Value:    "gvisor",
		Effect:   v1.TaintEffectNoSchedule,
	}
	GvisorNodeSelector = NodeSelector{
		"sandbox.gke.io/runtime": "gvisor",
	}
)
