package shelflifedetector

import (
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
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
		return fiber.ErrBadRequest
	}
	if file.Size > h.limitSize {
		h.log.ClientError(ctx, fmt.Errorf("file size is too large: %d", file.Size))
		return fiber.ErrRequestEntityTooLarge
	}
	fileContent, err := file.Open()
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	defer fileContent.Close()
	data, err := io.ReadAll(fileContent)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	dates, err := h.svc.Detect(data)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	if dates == nil {
		h.log.ServerError(ctx, fmt.Errorf("dates is nil"))
		return fiber.ErrBadGateway
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    dates,
	})
}
