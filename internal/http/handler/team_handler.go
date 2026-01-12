package handler

import (
	"net/http"

	"github.com/labstack/echo"

	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

/*
POST /team/add
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "alice",
      "is_active": true
    }
  ]
}
*/
type createTeamRequest struct {
	TeamName string               `json:"team_name"`
	Members  []models.TeamMember  `json:"members"`
}

func (h *TeamHandler) CreateTeam(c echo.Context) error {
	ctx := c.Request().Context()

	var req createTeamRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	if req.TeamName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "team_name is required",
			},
		})
	}

	team := &models.Team{
		TeamName: req.TeamName,
		Members:  req.Members,
	}

	err := h.teamService.CreateTeam(ctx, team)
	if err != nil {
		switch err {
		case service.ErrTeamExists:
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": map[string]string{
					"code":    "TEAM_EXISTS",
					"message": "team already exists",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not create team",
				},
			})
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status": "ok",
	})
}

/*
GET /team/get?name=backend
*/
func (h *TeamHandler) GetTeam(c echo.Context) error {
	ctx := c.Request().Context()

	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "team name is required",
			},
		})
	}

	team, err := h.teamService.GetTeam(ctx, name)
	if err != nil {
		switch err {
		case service.ErrTeamNotFound:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "TEAM_NOT_FOUND",
					"message": "team not found",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not get team",
				},
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"team": team,
	})
}
