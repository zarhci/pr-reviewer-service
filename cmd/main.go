package main

import (
	"log"
	"net/http"

	"pr-reviewer-service/internal/http/handler"
	"pr-reviewer-service/internal/repository/memory"
	"pr-reviewer-service/internal/service"
)

func main() {
	storage := memory.NewStorage()
	userRepo := memory.NewUserRepository(storage)
	teamRepo := memory.NewTeamRepository(storage)
	prRepo := memory.NewPRRepository(storage)

	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(userRepo)
	prService := service.NewPRService(prRepo, userRepo, teamRepo)

	teamHandler := handler.NewTeamHandler(teamService)
	userHandler := handler.NewUserHandler(userService, prService)
	prHandler := handler.NewPRHandler(prService)

	mux := http.NewServeMux()

	mux.HandleFunc("/team/add", teamHandler.CreateOrUpdate)
	mux.HandleFunc("/team/get", teamHandler.GetTeam)
	mux.HandleFunc("/users/setIsActive", userHandler.SetActive)
	mux.HandleFunc("/users/getReview", userHandler.GetReviews)
	mux.HandleFunc("/pullRequest/create", prHandler.CreatePR)
	mux.HandleFunc("/pullRequest/merge", prHandler.MergePR)
	mux.HandleFunc("/pullRequest/reassign", prHandler.ReassignReviewer)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
