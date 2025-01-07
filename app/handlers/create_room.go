package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pierre/app/db"
	"pierre/app/types"
	"pierre/app/views/components"
	"pierre/kit"
	"sort"
	"time"

	"math/rand"

	v "github.com/anthdm/superkit/validate"
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

var roomFormSchema = v.Schema{
	"roomCode": v.Rules(v.Min(4), v.Max(4)),
}

func JoinRoomWithCode(kit *kit.Kit, code string) error {
	var err error
	var room types.Room = types.Room{}
	room, err = getRoomByCode(kit, code)
	if err != nil {
		errors := v.Errors{}
		values := types.RoomFormValues{}
		errors.Add("roomCode", "Could not find a room with that code")
		return kit.Render(components.CreateOrJoinGame(values, errors))
	}

	return kit.Render(components.WaitForPoem(room, true))
}

func HandleJoinRoom(kit *kit.Kit) error {
	var values types.RoomFormValues
	errors, ok := v.Request(kit.Request, &values, roomFormSchema)
	if !ok {
		return kit.Render(components.CreateOrJoinGame(values, errors))
	}

	kit.Response.Header().Add("HX-Push-Url", "/poetry/room/"+values.RoomCode)

	return JoinRoomWithCode(kit, values.RoomCode)
}

func HandleGetAvailablePoem(kit *kit.Kit) error {

	var roomCode string = chi.URLParam(kit.Request, "code")

	var existingPoet *types.Poet
	var err error
	existingPoet, err = getPoetFromRequest(kit)
	if err != nil {
		return err
	}

	var room types.Room = types.Room{}
	room, err = getRoomByCode(kit, roomCode)
	if err != nil {
		return err
	}

	var allPoems []types.Poem
	var now = time.Now().Unix()
	err = db.Query.NewSelect().Model(&allPoems).
		Relation("Submissions").
		Relation("Submissions.Lines").
		Relation("Submissions.Poet").
		Where("room_id = ?", room.RoomID).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.WhereOr("reserved_poet_id = ?", existingPoet.PoetID).
				WhereOr("reserved_until_timestamp < ?", now).
				WhereOr("reserved_poet_id IS NULL")
			return q
		}).
		Scan(kit.Request.Context())

	if err != nil {
		return err
	}

	var incompletePoems []types.Poem
	var completedCount = 0
	for idx := range allPoems {
		if allPoems[idx].IsComplete {
			completedCount += 1
		} else {
			incompletePoems = append(incompletePoems, allPoems[idx])
		}
	}

	if completedCount == len(allPoems) {
		return kit.Render(components.ViewCompletedPoems(allPoems))
	}

	for idx := range incompletePoems {
		sort.Slice(incompletePoems[idx].Submissions, func(a, b int) bool {
			return incompletePoems[idx].Submissions[a].Position < incompletePoems[idx].Submissions[b].Position
		})
	}

	var bestIndex = -1
	var highestScore = 0.0
	var userPreviouslyReserved = false
	for idx, poem := range incompletePoems {

		if poem.ReservedPoetID != nil && *poem.ReservedPoetID == existingPoet.PoetID {
			bestIndex = idx
			userPreviouslyReserved = true
			break
		}

		var scoreMultiplier = 1.0
		if poem.ReservedPoetID == nil {
			scoreMultiplier = 1.5
		}

		var distance = 0.0
		var userHasSubmittedToThisPoem = false
		for i := len(poem.Submissions) - 1; i >= 0; i-- {
			if poem.Submissions[i].PoetID == existingPoet.PoetID {
				userHasSubmittedToThisPoem = true
				break
			}
			distance += 1.0
		}

		var isGreaterThanMinimumDistance = false
		var freshPoemBonus = 0.0
		if !userHasSubmittedToThisPoem {
			isGreaterThanMinimumDistance = true
			freshPoemBonus = 2.0
		}

		if distance >= float64(room.MinimumLineDistance) {
			isGreaterThanMinimumDistance = true
		}

		var score = scoreMultiplier*distance + freshPoemBonus

		if score > highestScore && isGreaterThanMinimumDistance {
			highestScore = score
			bestIndex = idx
		}

	}

	if bestIndex == -1 {
		return kit.Render(components.WaitForPoem(room, false))
	}

	var targetPoem = incompletePoems[bestIndex]

	var reservedUntilTimestamp = time.Now().Add(time.Duration(room.SecondsPerSubmission) * time.Second).Unix()

	if userPreviouslyReserved {
		reservedUntilTimestamp = targetPoem.ReservedUntilTimestamp
	} else {
		targetPoem.ReservedUntilTimestamp = reservedUntilTimestamp
	}

	_, err = db.Query.NewUpdate().Model(&targetPoem).
		Set("reserved_poet_id = ?", existingPoet.PoetID).
		Set("reserved_until_timestamp = ?", reservedUntilTimestamp).
		Where("poem_id = ?", targetPoem.PoemID).
		Exec(kit.Request.Context())

	if err != nil {
		return err
	}

	var lines []types.Line
	lines, err = getLines(kit, targetPoem, room)

	if err != nil {
		return kit.Render(components.WaitForPoem(room, false))
	}

	return kit.Render(components.LineSubmission(types.SubmissionFormTwoLineValues{}, v.Errors{}, lines, room, targetPoem))
}

func HandlePoemSubmission(kit *kit.Kit) error {
	var roomCode string = chi.URLParam(kit.Request, "code")
	var poemID string = chi.URLParam(kit.Request, "poemid")

	var err error

	var existingPoet *types.Poet

	existingPoet, err = getPoetFromRequest(kit)
	if err != nil {
		return err
	}

	var room types.Room = types.Room{}
	room, err = getRoomByCode(kit, roomCode)
	if err != nil {
		return err
	}

	var poem types.Poem
	err = db.Query.NewSelect().Model(&poem).
		Relation("Room").
		Relation("Submissions").
		Where("reserved_poet_id = ?", existingPoet.PoetID).
		Where("poem_id = ?", poemID).
		Where("code = ?", roomCode).
		Scan(kit.Request.Context())

	if err != nil {
		return err
	}

	var values types.SubmissionFormTwoLineValues
	errors, ok := v.Request(kit.Request, &values, v.Schema{
		"line0":    v.Rules(v.Min(1), v.Max(200)),
		"line1":    v.Rules(v.Min(1), v.Max(200)),
		"lastLine": v.Rules(),
	})

	if !ok {
		var lines []types.Line
		lines, err = getLines(kit, poem, room)

		if err != nil {
			return err
		}

		return kit.Render(components.LineSubmission(values, errors, lines, room, poem))
	}

	var highestPosition = -1
	for _, submission := range poem.Submissions {
		if highestPosition < submission.Position {
			highestPosition = submission.Position
		}
	}

	var ctx = kit.Request.Context()

	err = db.Query.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		submission := types.Submission{
			PoemID:   poem.PoemID,
			PoetID:   existingPoet.PoetID,
			Position: highestPosition + 1,
		}
		if _, err := tx.NewInsert().Model(&submission).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert submission: %w", err)
		}

		var lines = []*types.Line{
			{SubmissionID: submission.SubmissionID, Position: 0, Text: values.Line0},
			{SubmissionID: submission.SubmissionID, Position: 1, Text: values.Line1},
		}

		if _, err := tx.NewInsert().Model(&lines).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert lines: %w", err)
		}

		if _, err := tx.NewUpdate().Model(&poem).
			Set("reserved_poet_id = ?", nil).
			Set("is_complete = ?", values.IsLastLine).
			Where("poem_id = ?", poemID).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to update poem: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return kit.Render(components.WaitForPoem(room, true))
}

func getLines(kit *kit.Kit, poem types.Poem, room types.Room) ([]types.Line, error) {

	var lines []types.Line
	var err = db.Query.NewSelect().Model(&lines).
		Relation("Submission").
		Where("poem_id = ?", poem.PoemID).
		Order("submission.position ASC", "line.position ASC").
		Scan(kit.Request.Context())

	if err != nil {
		return lines, err
	}
	if len(lines) < room.LinesVisible {
		return lines, err
	}

	return lines[len(lines)-room.LinesVisible:], err
}

func HandleCreateRoom(kit *kit.Kit) error {

	var errors = v.Errors{}
	var existingPoet, err = getPoetFromRequest(kit)
	var values types.RoomFormValues

	if err != nil {
		errors.Add("Token", "No poet exists")
		return kit.Render(components.CreateOrJoinGame(values, errors))
	}

	errors, ok := v.Request(kit.Request, &values, v.Schema{
		"poemCount":           v.Rules(v.GTE(1), v.LTE(10)),
		"minimumLineDistance": v.Rules(v.GTE(1), v.LTE(10)),
		"roomCode":            v.Rules(),
	})
	if !ok {
		return kit.Render(components.CreateOrJoinGame(values, errors))
	}

	var room, roomErr = createRoom(kit, values, *existingPoet)
	if roomErr != nil {
		errors.Add("Creation", "Could not create room")
		return kit.Render(components.CreateOrJoinGame(values, errors))
	}

	kit.Response.Header().Add("HX-Push-Url", "/poetry/room/"+room.Code)

	return kit.Render(components.WaitForPoem(room, true))
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func GenerateRandomCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func getRoomByCode(kit *kit.Kit, code string) (types.Room, error) {
	var room types.Room
	err := db.Query.NewSelect().
		Model(&room).
		Where("code = ?", code).
		Scan(kit.Request.Context())

	return room, err
}

func createRoom(kit *kit.Kit, values types.RoomFormValues, _poet types.Poet) (types.Room, error) {

	var uniqueCode string
	var found = true
	for i := 0; i < 5; i++ {
		uniqueCode = GenerateRandomCode(4)
		var _, err = getRoomByCode(kit, uniqueCode)
		if err != nil {
			found = false
			break
		}
	}

	room := types.Room{
		Code:                 uniqueCode,
		LinesPerSubmission:   2,
		LinesVisible:         1,
		SecondsPerSubmission: 90,
		MinimumLineDistance:  values.MinimumLineDistance,
	}

	if found {
		var err = errors.New("could not generate a unique room code")
		return room, err
	}

	_, err := db.Query.NewInsert().
		Model(&room).
		Exec(kit.Request.Context())
	if err != nil {
		return room, err
	}

	var poems = []types.Poem{}

	for range values.PoemCount {
		poems = append(poems, types.Poem{RoomID: room.RoomID, ReservedPoetID: nil, IsComplete: false})
	}

	_, err = db.Query.NewInsert().Model(&poems).Exec(kit.Request.Context())
	if err != nil {
		return room, err
	}

	return room, nil
}
