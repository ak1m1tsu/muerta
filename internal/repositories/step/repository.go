package steps

import (
	"context"

	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type StepsRepositorer interface {
	Create(ctx context.Context, recipe *models.Step) error
}
