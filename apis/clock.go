package apis

import (
	kubeclock "k8s.io/apimachinery/pkg/util/clock"
)

// clock is used to set status condition timestamps.
// This variable makes it easier to test conditions.
var clock kubeclock.Clock = &kubeclock.RealClock{}
