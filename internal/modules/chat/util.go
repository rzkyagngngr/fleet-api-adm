package chat

import (
	"fmt"
	"strconv"
	"strings"
)

func parseUint64(value string) (uint64, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, fmt.Errorf("value is required")
	}
	parsed, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}
	return parsed, nil
}

func validVesselCallStatus(status string) bool {
	switch status {
	case VesselCallStatusActive, VesselCallStatusSuspended:
		return true
	default:
		return false
	}
}

func validParticipantRole(role string) bool {
	switch role {
	case ParticipantRoleAgent, ParticipantRolePBM, ParticipantRoleOperator:
		return true
	default:
		return false
	}
}
