package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DangeL187/erax"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	"back/internal/features/subs/domain"
	"back/internal/features/subs/dto"
	"back/internal/features/subs/usecase"
	_ "back/internal/shared/swagger"
)

type SubsHandler struct {
	crudl *usecase.CrudlUseCase
	sum   *usecase.SumUseCase
}

func NewSubsHandler(crudl *usecase.CrudlUseCase, sum *usecase.SumUseCase) *SubsHandler {
	return &SubsHandler{crudl: crudl, sum: sum}
}

// Create создаёт новую подписку.
// @Summary Create subscription
// @Description Создает новую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.SubscriptionJSON true "Subscription payload"
// @Success 201 {object} dto.SubscriptionDTOSwagger
// @Failure 400 {object} swagger.HTTPErrorDTO
// @Failure 409 {object} swagger.HTTPErrorDTO
// @Failure 500 {object} swagger.HTTPErrorDTO
// @Router /subscriptions [post]
func (h *SubsHandler) Create(c *echo.Context) error {
	subJSON, err := bindSubscription(c)
	if err != nil {
		return err
	}

	sub, err := mapSubJSONToDTO(subJSON)
	if err != nil {
		return err
	}

	err = h.crudl.Create(c.Request().Context(), sub)
	if errors.Is(err, domain.ErrSubscriptionAlreadyExists) {
		err = erax.Wrap(err, "failed to create subscription")
		zap.L().Debug("\n" + erax.Format(err))
		return echo.NewHTTPError(http.StatusConflict, "subscription already exists")
	}
	if err != nil {
		err = erax.Wrap(err, "failed to create subscription")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusCreated, sub)
}

// Get получает подписку по ID
// @Summary Get subscription by ID
// @Description Возвращает подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionDTOSwagger
// @Failure 404 {object} swagger.HTTPErrorDTO
// @Failure 500 {object} swagger.HTTPErrorDTO
// @Router /subscriptions/{id} [get]
func (h *SubsHandler) Get(c *echo.Context) error {
	id, err := parseIDParam(c)
	if err != nil {
		return err
	}

	sub, err := h.crudl.GetByID(c.Request().Context(), id)
	if errors.Is(err, domain.ErrSubscriptionNotFound) {
		err = erax.Wrap(err, "failed to find subscription by ID")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusNotFound, "subscription not found")
	}
	if err != nil {
		err = erax.Wrap(err, "failed to get subscription by ID")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusOK, sub)
}

// Update обновляет существующую подписку.
// @Summary Update subscription
// @Description Обновляет подписку по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body dto.SubscriptionJSON true "Subscription payload"
// @Success 200 {object} dto.SubscriptionDTOSwagger
// @Failure 400 {object} swagger.HTTPErrorDTO
// @Failure 404 {object} swagger.HTTPErrorDTO
// @Failure 500 {object} swagger.HTTPErrorDTO
// @Router /subscriptions/{id} [put]
func (h *SubsHandler) Update(c *echo.Context) error {
	id, err := parseIDParam(c)
	if err != nil {
		return err
	}

	subJSON, err := bindSubscription(c)
	if err != nil {
		return err
	}

	sub, err := mapSubJSONToDTO(subJSON)
	if err != nil {
		return err
	}
	sub.ID = id

	err = h.crudl.Update(c.Request().Context(), sub)
	if err != nil {
		err = erax.Wrap(err, "failed to update subscription")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusOK, sub)
}

// Delete удаляет подписку по ID.
// @Summary Delete subscription
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} swagger.HTTPErrorDTO
// @Failure 404 {object} swagger.HTTPErrorDTO
// @Failure 500 {object} swagger.HTTPErrorDTO
// @Router /subscriptions/{id} [delete]
func (h *SubsHandler) Delete(c *echo.Context) error {
	id, err := parseIDParam(c)
	if err != nil {
		return err
	}

	err = h.crudl.DeleteByID(c.Request().Context(), id)
	if err != nil {
		err = erax.Wrap(err, "failed to delete subscription by ID")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}

// List возвращает список подписок с опциональными фильтрами.
// @Summary List subscriptions
// @Description Возвращает все подписки пользователя, можно фильтровать по user_id и service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User UUID"
// @Param service query string false "Service name"
// @Success 200 {array} dto.SubscriptionDTOSwagger
// @Failure 400 {object} swagger.HTTPErrorDTO
// @Failure 500 {object} swagger.HTTPErrorDTO
// @Router /subscriptions [get]
func (h *SubsHandler) List(c *echo.Context) error {
	userID, err := parseUserIDQuery(c)
	if err != nil {
		return err
	}
	service := c.QueryParam("service")

	subs, err := h.crudl.List(c.Request().Context(), userID, &service)
	if err != nil {
		err = erax.Wrap(err, "failed to list subscriptions")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusOK, subs)
}

// GetSubscriptionsCostForPeriod считает суммарную стоимость подписок за период
// @Summary Get total subscriptions cost for period
// @Description Возвращает суммарную стоимость подписок за период с опциональными фильтрами
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start date (DD-MM-YYYY)"
// @Param to query string true "End date (DD-MM-YYYY)"
// @Param user_id query string false "User UUID"
// @Param service query string false "Service name"
// @Success 200 {integer} int64
// @Failure 400 {object} swagger.HTTPErrorDTO
// @Failure 500 {object} swagger.HTTPErrorDTO
// @Router /subscriptions/costs [get]
func (h *SubsHandler) GetSubscriptionsCostForPeriod(c *echo.Context) error {
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	userID, err := parseUserIDQuery(c)
	if err != nil {
		return err
	}
	serviceStr := c.QueryParam("service")
	var service *string = nil
	if serviceStr != "" {
		service = &serviceStr
	}

	from, err := time.Parse("02-01-2006", fromStr)
	if err != nil {
		err = erax.Wrap(err, "failed to parse `from` date")
		zap.L().Debug("\n" + erax.Format(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid `from` date")
	}
	to, err := time.Parse("02-01-2006", toStr)
	if err != nil {
		err = erax.Wrap(err, "failed to parse `to` date")
		zap.L().Debug("\n" + erax.Format(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid `to` date")
	}

	totalCost, err := h.sum.GetSubscriptionsCostForPeriod(c.Request().Context(), from, to, userID, service)
	if err != nil {
		err = erax.Wrap(err, "failed to get subscriptions cost for period")
		zap.L().Error(erax.FormatToJSONString(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusOK, totalCost)
}

// --- Helpers ---
func mapSubJSONToDTO(subJSON dto.SubscriptionJSON) (dto.SubscriptionDTO, error) {
	userID, err := uuid.Parse(subJSON.UserID)
	if err != nil {
		err = erax.Wrap(err, "failed to parse UserID")
		zap.L().Debug("\n" + erax.Format(err))
		return dto.SubscriptionDTO{}, echo.NewHTTPError(http.StatusBadRequest, "invalid user_id")
	}

	start, err := parseMonthYear(subJSON.StartDate)
	if err != nil {
		err = erax.Wrap(err, "failed to parse StartDate")
		zap.L().Debug("\n" + erax.Format(err))
		return dto.SubscriptionDTO{}, echo.NewHTTPError(http.StatusBadRequest, "invalid start_date")
	}

	end := pgtype.Date{
		Valid:            true,
		InfinityModifier: pgtype.Infinity,
	}
	if subJSON.EndDate != "" {
		var t time.Time
		t, err = parseMonthYear(subJSON.EndDate)
		if err != nil {
			err = erax.Wrap(err, "failed to parse EndDate")
			zap.L().Debug("\n" + erax.Format(err))
			return dto.SubscriptionDTO{}, echo.NewHTTPError(http.StatusBadRequest, "invalid end_date")
		}

		end.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		end.InfinityModifier = pgtype.Finite
	}

	return dto.SubscriptionDTO{
		UserID:      userID,
		ServiceName: subJSON.ServiceName,
		StartDate:   start,
		EndDate:     end,
		PricePerDay: subJSON.PricePerDay,
	}, nil
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

func bindSubscription(c *echo.Context) (dto.SubscriptionJSON, error) {
	var s dto.SubscriptionJSON
	err := c.Bind(&s)
	if err != nil {
		err = erax.Wrap(err, "failed to bind request")
		zap.L().Debug("\n" + erax.Format(err))
		return dto.SubscriptionJSON{}, echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	return s, nil
}

func parseIDParam(c *echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		err = erax.Wrap(err, "failed to parse ID param")
		zap.L().Debug("\n" + erax.Format(err))
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid ID")
	}

	return id, nil
}

func parseUserIDQuery(c *echo.Context) (*uuid.UUID, error) {
	v := c.QueryParam("user_id")
	if v == "" {
		return nil, nil
	}

	u, err := uuid.Parse(v)
	if err != nil {
		err = erax.Wrap(err, "failed to parse UserID query")
		zap.L().Debug("\n" + erax.Format(err))
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user_id")
	}

	return &u, nil
}
