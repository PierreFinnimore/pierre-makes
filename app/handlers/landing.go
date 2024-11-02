package handlers

import (
	"pierre/app/views/landing"

	"pierre/kit"
)

func HandleLandingIndex(kit *kit.Kit) error {
	return kit.Render(landing.Index())
}
