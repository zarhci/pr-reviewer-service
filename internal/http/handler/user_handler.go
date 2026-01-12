package handler

import (
	"net/http"

	"github.com/labstack/echo"

	"pr-reviewer-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

/*
POST /users/setIsActive

	{
	  "user_id": "u1",
	  "is_active": true
	}
*/
type setActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) SetActive(c echo.Context) error {
	ctx := c.Request().Context()

	var req setActiveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "user_id is required",
			},
		})
	}

	user, err := h.userService.SetActive(
		ctx,
		req.UserID,
		req.IsActive,
	)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "USER_NOT_FOUND",
					"message": "user not found",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not update user",
				},
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

/*
GET /users/:id
*/
func (h *UserHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "user id is required",
			},
		})
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "USER_NOT_FOUND",
					"message": "user not found",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not get user",
				},
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
