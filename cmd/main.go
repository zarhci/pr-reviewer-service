package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"pr-reviewer-service/internal/http/handler"
	r "pr-reviewer-service/internal/repository"
	"pr-reviewer-service/internal/service"
)

func main() {
	// --- Database ---
	database, err := r.New(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	// --- Repositories ---
	userRepo := r.NewUserRepository(database.DB)
	teamRepo := r.NewTeamRepository(database.DB)
	prRepo := r.NewPullRequestRepository(database.DB)

	// --- Services ---
	userService := service.NewUserService(userRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)
	prService := service.NewPRService(prRepo, userRepo, teamRepo)

	// --- Handlers ---
	teamHandler := handler.NewTeamHandler(teamService)
	userHandler := handler.NewUserHandler(userService)
	prHandler := handler.NewPRHandler(prService)

	// --- Echo ---
	e := echo.New()

	// --- Middleware ---
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- Routes ---

	// Team
	e.POST("/team/add", teamHandler.CreateTeam)
	e.GET("/team/get", teamHandler.GetTeam)

	// User
	e.POST("/users/setIsActive", userHandler.SetActive)
	e.GET("/users/:id", userHandler.GetByID)

	// Pull Requests
	e.POST("/pullRequest/create", prHandler.CreatePR)
	e.POST("/pullRequest/merge", prHandler.MergePR)
	e.POST("/pullRequest/reassign", prHandler.ReassignReviewer)

	// Healthcheck
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// --- Start server ---
	log.Println("Starting server on port 8080...")
	log.Fatal(e.Start(":8080"))
}
