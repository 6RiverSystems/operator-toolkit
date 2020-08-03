package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReadyStatus defines base fields for resource readiness status
// +k8s:deepcopy-gen=true
type ReadyStatus struct {
	// LastUpdate: last time the status has been updated
	LastUpdate metav1.Time `json:"lastUpdate"`

	// Reason: contains possible failure reson in human readable format
	Reason string `json:"reason,omitempty"`

	// Ready: database is ready to be used
	Ready bool `json:"ready"`
}

// FailedReadyStatus returns failed status
func FailedReadyStatus(err error) ReadyStatus {
	return ReadyStatus{
		LastUpdate: metav1.Time{Time: clock.Now()},
		Reason:     err.Error(),
		Ready:      false,
	}
}

// ReadyStatusOK returns OK status
func ReadyStatusOK() ReadyStatus {
	return ReadyStatus{
		LastUpdate: metav1.Time{Time: clock.Now()},
		Ready:      true,
	}
}
