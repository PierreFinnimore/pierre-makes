package handlers

import (
	"pierre/app/types"

	"pierre/kit"
)

func HandleAuthentication(kit *kit.Kit) (kit.Auth, error) {
	return types.AuthUser{}, nil
}
