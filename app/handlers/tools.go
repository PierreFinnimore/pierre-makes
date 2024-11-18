package handlers

import (
	"pierre/app/views/tools"

	"pierre/kit"
)

func HandleToolsIndex(kit *kit.Kit) error {
	return kit.Render(tools.Tools())
}
