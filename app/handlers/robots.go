package handlers

import (
	"net/http"
	"pierre/kit"
)

func HandleRobotsTxt(kit *kit.Kit) error {
	return kit.Text(http.StatusOK, "")
}
