package internal

import "github.com/anthropics/anthropic-sdk-go"

type Role string

const (
	// RoleUser represents the user role.
	RoleUser Role = "user"
	// RoleModel represents the model role.
	RoleModel Role = "model"
)

func (r Role) ToClaude() anthropic.MessageParamRole {
	switch r {
	case RoleUser:
		return "user"
	case RoleModel:
		return "assistant"
	default:
		return ""
	}
}
