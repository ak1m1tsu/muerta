package context

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/pkg/log"
)

func New(log *log.Logger, key idKey) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id, err := common.ParseIDFromPath(ctx, key.String())
		if err != nil {
			log.ClientError(ctx, err)
			return ctx.Status(http.StatusNotFound).
				JSON(handlers.HTTPError{Error: fiber.ErrNotFound.Error()})
		}
		ctx.Locals(key, id)
		return ctx.Next()
	}
}

type idKey string

func (id idKey) Path() string {
	return fmt.Sprintf("/:%s<int>", id)
}

func (id idKey) String() string {
	return string(id)
}

const (
	ShelfLifeID idKey = "shelf_life_id"
	StatusID    idKey = "status_id"
	StorageID   idKey = "storage_id"
	TypeID      idKey = "type_id"
	ProductID   idKey = "product_id"
	MeasureID   idKey = "measure_id"
	CategoryID  idKey = "category_id"
	RecipeID    idKey = "recipe_id"
	StepID      idKey = "step_id"
	TipID       idKey = "tip_id"
	UserID      idKey = "user_id"
	SettingID   idKey = "setting_id"
	RoleID      idKey = "role_id"
)
