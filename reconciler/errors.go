package reconciler

import (
	"time"

	"github.com/6RiverSystems/operator-toolkit/apis"
	corev1 "k8s.io/api/core/v1"
)

// ConditionError marks error as failed condition one
type ConditionError struct {
	Type   apis.ConditionType
	Reason apis.ConditionReason
	Err    error
}

func (e *ConditionError) Error() string {
	return e.Err.Error()
}

// Condition converts consdition error into condition structure
func (e *ConditionError) Condition() apis.Condition {
	return apis.Condition{
		Type:    e.Type,
		Status:  corev1.ConditionFalse,
		Reason:  e.Reason,
		Message: e.Error(),
	}
}

// NewConditionError makes new condition error
func NewConditionError(ctype apis.ConditionType, issue error) error {
	return NewConditionErrorWithReason(ctype, "Failed", issue)
}

// NewConditionErrorWithReason makes new condition error with provided reason
func NewConditionErrorWithReason(
	ctype apis.ConditionType,
	reason apis.ConditionReason,
	issue error,
) error {
	return &ConditionError{
		Type:   ctype,
		Reason: reason,
		Err:    issue,
	}
}

type notRetriableError struct {
	issue error
}

func (e *notRetriableError) Error() string {
	return e.issue.Error()
}

// NewNotRetriableError returns notretriable error
func NewNotRetriableError(issue error) error {
	return &notRetriableError{issue}
}

type retriableError struct {
	issue      error
	retryAfter time.Duration
}

func (e *retriableError) Error() string {
	return e.issue.Error()
}

// NewRetriableError returns retriable error
func NewRetriableError(retryAfter time.Duration, issue error) error {
	return &retriableError{
		retryAfter: retryAfter,
		issue:      issue,
	}
}
