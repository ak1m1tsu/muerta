package shelflifedetector

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	sldetector "github.com/romankravchuk/muerta/internal/services/shelf-life-detector"
)

type ShelfLifeDetectorController struct {
	svc       sldetector.DateDetectorServicer
	log       logger.Logger
	limitSize int64
}

func New(svc sldetector.DateDetectorServicer, log logger.Logger) *ShelfLifeDetectorController {
	return &ShelfLifeDetectorController{
		svc: svc,
		log: log,
		// Limit - 512KB
		limitSize: 1024 * 512,
	}
}

// DetectDates - detects shelf life dates from file
//
//	@Summary		Detect shelf life dates from file
//	@Description	detect shelf life dates from file
//	@Tags			Shelf Life Detector
//	@Accept			json
//	@Produce		json
//	@Param			fileToDetect	formData	file	true	"file to detect"
//	@Success		200				{object}	handlers.HTTPError
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-life-detector [post]
//	@Security		Bearer
func (h *ShelfLifeDetectorController) DetectDates(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("fileToDetect")
	if err != nil {
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if file.Size > h.limitSize {
		h.log.Error(ctx, logger.Client, fmt.Errorf("file size is too large: %d", file.Size))
		return fiber.ErrRequestEntityTooLarge
	}
	fileContent, err := file.Open()
	if err != nil {
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	defer fileContent.Close()
	data, err := io.ReadAll(fileContent)
	if err != nil {
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	dates, err := h.svc.Detect(data)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	if dates == nil {
		h.log.Error(ctx, logger.Server, fmt.Errorf("dates is nil"))
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    dates,
	})
}
