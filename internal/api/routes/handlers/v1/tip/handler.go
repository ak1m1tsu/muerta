package tip

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/tip"
)

type TipHandler struct {
	svc service.TipServicer
	log *log.Logger
}

func New(svc service.TipServicer, log *log.Logger) *TipHandler {
	return &TipHandler{
		svc: svc,
		log: log,
	}
}

// Create godoc
//
//	@Summary		Create tip
//	@Description	Create tip
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			tip	body		dto.CreateTip	true	"tip"
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Failure		400	{object}	handlers.HTTPError
//	@Failure		500	{object}	handlers.HTTPError
//	@Router			/tips [post]
//	@Security		Bearer
func (h *TipHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateTip
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.CreateTip(ctx.Context(), payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

// FindOne godoc
//
//	@Summary		Find tip by id
//	@Description	Find tip by id
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip	path		int	true	"tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips/{id_tip} [get]
func (h *TipHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

// FindMany godoc
//
//	@Summary		Find many tips
//	@Description	Find many tips
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.TipFilter	true	"filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips [get]
func (h *TipHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.TipFilter)
	if err := common.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindTips(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tips": result, "count": count}},
	)
}

// Update godoc
//
//	@Summary		Update tip
//	@Description	Update tip
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip	path		int				true	"tip id"
//	@Param			tip		body		dto.UpdateTip	true	"tip"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips/{id_tip} [put]
//	@Security		Bearer
func (h *TipHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	payload := new(dto.UpdateTip)
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateTip(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete tip
//	@Description	Delete tip
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip	path		int	true	"tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips/{id_tip} [delete]
//	@Security		Bearer
func (h *TipHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Restore godoc
//
//	@Summary		Restore tip
//	@Description	Restore tip
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip	path		int	true	"tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips/{id_tip} [patch]
func (h *TipHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	if err := h.svc.RestoreTip(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindStorages godoc
//
//	@Summary		Find storages
//	@Description	Find storages
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip	path		int	true	"tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips/{id_tip}/storages [get]
func (h *TipHandler) FindStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipStorages(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
}
// FindProducts godoc
//
//	@Summary		Find products
//	@Description	Find products
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip	path		int	true	"tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/tips/{id_tip}/products [get]
func (h *TipHandler) FindProducts(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipProducts(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"products": result}})
}
// AddProduct godoc
//
//	@Summary		Add product
//	@Description	Add product
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip		path		int	true	"tip id"
//	@Param			id_product	path		int	true	"product id"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/tips/{id_tip}/products/{id_product} [post]
//	@Security		Bearer
func (h *TipHandler) AddProduct(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	productID := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.AddProductToTip(ctx.Context(), tipID, productID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"product": result}})
}
// RemoveProduct godoc
//
//	@Summary		Remove product
//	@Description	Remove product
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip		path		int	true	"tip id"
//	@Param			id_product	path		int	true	"product id"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/tips/{id_tip}/products/{id_product} [delete]
//	@Security		Bearer 
func (h *TipHandler) RemoveProduct(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	productID := ctx.Locals(context.ProductID).(int)
	if err := h.svc.RemoveProductFromTip(ctx.Context(), tipID, productID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// AddStorage godoc
//
//	@Summary		Add storage
//	@Description	Add storage
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip		path		int	true	"tip id"
//	@Param			id_storage	path		int	true	"storage id"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/tips/{id_tip}/storages/{id_storage} [post]
//	@Security		Bearer
func (h *TipHandler) AddStorage(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	storateID := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.AddStorageToTip(ctx.Context(), tipID, storateID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storage": result}})
}

// RemoveStorage godoc
//
//	@Summary		Remove storage
//	@Description	Remove storage
//	@Tags			Tips
//	@Accept			json
//	@Produce		json
//	@Param			id_tip		path		int	true	"tip id"
//	@Param			id_storage	path		int	true	"storage id"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/tips/{id_tip}/storages/{id_storage} [delete]
//	@Security		Bearer
func (h *TipHandler) RemoveStorage(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	storateID := ctx.Locals(context.StorageID).(int)
	if err := h.svc.RemoveStorageFromTip(ctx.Context(), tipID, storateID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}
