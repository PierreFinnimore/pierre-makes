package handlers

import (
	"pierre/app/views/thoughts"

	"pierre/kit"
)

func HandleThoughtsIndex(kit *kit.Kit) error {
	return kit.Render(thoughts.Thoughts())
}
