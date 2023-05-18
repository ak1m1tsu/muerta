package shelflifedetector

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	sldetector "github.com/romankravchuk/muerta/internal/services/shelf-life-detector"
)

type ShelfLifeDetectorHandler struct {
	svc       sldetector.DateDetectorServicer
	log       *log.Logger
	limitSize int64
}

func New(svc sldetector.DateDetectorServicer, log *log.Logger) *ShelfLifeDetectorHandler {
	return &ShelfLifeDetectorHandler{
		svc: svc,
		log: log,
		// Limit - 512KB
		limitSize: 1024 * 512,
	}
}

func (h *ShelfLifeDetectorHandler) DetectDates(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("fileToDetect")
	if err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if file.Size > h.limitSize {
		h.log.ClientError(ctx, fmt.Errorf("file size is too large: %d", file.Size))
		return fiber.ErrRequestEntityTooLarge
	}
	fileContent, err := file.Open()
	if err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	defer fileContent.Close()
	data, err := io.ReadAll(fileContent)
	if err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	dates, err := h.svc.Detect(data)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	if dates == nil {
		h.log.ServerError(ctx, fmt.Errorf("dates is nil"))
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    dates,
	})
}