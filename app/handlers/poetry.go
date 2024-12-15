package handlers

import (
	"pierre/app/views/poetry"

	"pierre/kit"

	"github.com/go-chi/chi/v5"
)

func HandleConsequencesIndex(kit *kit.Kit) error {
	return kit.Render(poetry.PoetryConsequences())
}

func HandleGetRoomIndex(kit *kit.Kit) error {
	var roomCode string = chi.URLParam(kit.Request, "code")
	return kit.Render(poetry.PoetryRoom(roomCode))
}
