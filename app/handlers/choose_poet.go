package handlers

import (
	"pierre/app/db"
	"pierre/app/types"
	"pierre/app/views/components"
	"pierre/kit"

	v "github.com/anthdm/superkit/validate"
	"github.com/google/uuid"
)

var poetFormSchema = v.Schema{
	"poetName": v.Rules(v.Min(1), v.Max(50)),
}

func HandleChoosePoet(kit *kit.Kit) error {

	var values types.PoetFormValues
	errors, ok := v.Request(kit.Request, &values, poetFormSchema)
	if !ok {
		errors.Add("token", "Something went wrong choosing poet")
		return kit.Render(components.ChoosePoet(values, errors))
	}

	var existingPoet *types.Poet
	var err error
	existingPoet, err = getPoetFromRequest(kit)
	if err != nil {
		errors.Add("token", "Could not find this poet")
		return kit.Render(components.ChoosePoet(values, errors))
	}

	existingPoet.Name = values.PoetName
	_, err = db.Query.NewUpdate().Model(existingPoet).Where("poet_id = ?", existingPoet.PoetID).Exec(kit.Request.Context())
	if err != nil {
		errors.Add("token", "Could not find this poet")
		return kit.Render(components.ChoosePoet(values, errors))
	}

	return kit.Render(components.ChoosePoet(values, v.Errors{}))
}

func getPoetFromRequest(kit *kit.Kit) (*types.Poet, error) {
	var token = kit.Request.Header.Get("poet_token")
	var poet types.Poet
	err := db.Query.NewSelect().
		Model(&poet).
		Where("token = ?", token).
		Scan(kit.Request.Context())

	return &poet, err
}

func createPoet(kit *kit.Kit, token string) (*types.Poet, error) {
	poet := types.Poet{
		Token: token,
	}
	_, err := db.Query.NewInsert().
		Model(&poet).
		Exec(kit.Request.Context())
	if err != nil {
		return &poet, err
	}
	return &poet, nil
}

func HandleGetPoet(kit *kit.Kit) error {
	var poetValues types.PoetFormValues
	var poetErrors = v.Errors{}

	var roomValues types.RoomFormValues
	var roomErrors = v.Errors{}

	var poetToken = kit.Request.Header.Get("poet_token")

	if poetToken == "" {
		// TODO this should be a proper error page, as we expect
		// the token to exist.
		poetErrors.Add("token", "No Token Found")
		return kit.Render(components.GetPoet(poetValues, poetErrors, roomValues, roomErrors))
	}

	var existingPoet *types.Poet
	var err error
	existingPoet, err = getPoetFromRequest(kit)
	if err != nil {
		// TODO this should be a proper error page, as we expect
		// the token to exist.
		poetErrors.Add("token", "No Poet Found")
		return kit.Render(components.GetPoet(poetValues, poetErrors, roomValues, roomErrors))
	}

	poetValues.PoetName = existingPoet.Name

	return kit.Render(components.GetPoet(poetValues, v.Errors{}, roomValues, roomErrors))
}

func createPoetToken() string {
	return uuid.New().String()
}

func HandleGetPoetAuth(kit *kit.Kit) error {

	var err error
	var token string = ""

	{

		var poet, err = getPoetFromRequest(kit)
		if err == nil {
			token = poet.Token
		}

	}

	if token == "" {
		token = createPoetToken()
		_, err = createPoet(kit, token)
		if err != nil {
			return kit.JSON(500, map[string]string{
				"error": "Could not create poet",
			})
		}
	}

	return kit.JSON(200, map[string]string{
		"token": token,
	})
}
