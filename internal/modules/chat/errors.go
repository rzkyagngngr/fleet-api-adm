package chat

import "errors"

var (
	ErrScheduleNotFound    = errors.New("schedule not found")
	ErrVesselCallNotFound  = errors.New("schedule topic not found")
	ErrScheduleTopicExists = errors.New("schedule topic already exists")
	ErrParticipantNotFound = errors.New("participant not found")
	ErrParticipantExists   = errors.New("participant already exists in schedule")
	ErrInvalidRole         = errors.New("invalid participant role")
	ErrInvalidCallStatus   = errors.New("invalid schedule topic status")
	ErrInternal            = errors.New("internal chat error")
)

