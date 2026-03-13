package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Team struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	City          string `json:"city"`
	Championships int    `json:"championships"`
	Pet           string `json:"pet"`
	Arena         string `json:"arena"`
}

type Message struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

var teams []Team

func main() {
	loadTeams()

	http.HandleFunc("/api/teams/jayson", pingHandler)
	http.HandleFunc("/api/teams", teamsHandler)
	http.HandleFunc("/api/teams/", teamByIDHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadTeams() {
	file, err := os.ReadFile("ejercicio4web/data/teams.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(file, &teams)
	if err != nil {
		fmt.Println("Error unmarshalling file:", err)
		return
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := Message{
		Message: "Tatum",
	}

	writeJSON(w, http.StatusOK, response)
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		handleGetTeams(w, r)
	case http.MethodPost:
		handlePostTeam(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, Message{Message: "Method not allowed"})
	}
}

func handleGetTeams(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")

	if idParam == "" {
		writeJSON(w, http.StatusOK, teams)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	for _, team := range teams {
		if team.ID == id {
			writeJSON(w, http.StatusOK, team)
			return
		}
	}

	http.Error(w, "Team not found", http.StatusNotFound)
}

func generateID() int {
	maxID := 0
	for _, team := range teams {
		if team.ID > maxID {
			maxID = team.ID
		}
	}
	return maxID + 1
}

func saveTeams() error {
	data, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling error: %w", err)
	}
	if err = os.WriteFile("ejercicio4web/data/teams.json", data, 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

func teamByIDHandler(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/teams/")

	if idStr == "" {
		teamsHandler(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid path parameter",
			fmt.Sprintf("'%s' is not a valid integer id", idStr))
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetTeamByID(w, r, id)
	case http.MethodPut:
		handlePutTeam(w, r, id)
	case http.MethodPatch:
		handlePatchTeam(w, r, id)
	case http.MethodDelete:
		handleDeleteTeam(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed",
			fmt.Sprintf("'%s' is not supported on this endpoint", r.Method))
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, errMsg, details string) {
	writeJSON(w, status, ErrorResponse{
		Error:   errMsg,
		Code:    status,
		Details: details,
	})
}

func handleGetTeamByID(w http.ResponseWriter, r *http.Request, id int) {
	for _, team := range teams {
		if team.ID == id {
			writeJSON(w, http.StatusOK, team)
			return
		}
	}

	writeError(w, http.StatusNotFound, "Team not found", "No team found with the specified ID")
}

func handlePostTeam(w http.ResponseWriter, r *http.Request) {
	var newTeam Team
	if err := json.NewDecoder(r.Body).Decode(&newTeam); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body",
			"Request body must be valid JSON matching the Team schema")
		return
	}

	if validationErr := validateTeam(newTeam, false); validationErr != "" {
		writeError(w, http.StatusUnprocessableEntity, "Validation failed", validationErr)
		return
	}

	newTeam.ID = generateID()
	teams = append(teams, newTeam)

	if err := saveTeams(); err != nil {
		writeError(w, http.StatusInternalServerError, "Persistence error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, newTeam)
}

func handlePutTeam(w http.ResponseWriter, r *http.Request, id int) {
	idx := findIndex(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "Team not found",
			fmt.Sprintf("No team with id=%d exists", id))
		return
	}

	var updated Team
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body",
			"Request body must be valid JSON matching the Team schema")
		return
	}

	if validationErr := validateTeam(updated, false); validationErr != "" {
		writeError(w, http.StatusUnprocessableEntity, "Validation failed", validationErr)
		return
	}

	updated.ID = id
	teams[idx] = updated

	if err := saveTeams(); err != nil {
		writeError(w, http.StatusInternalServerError, "Persistence error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, teams[idx])
}

func handlePatchTeam(w http.ResponseWriter, r *http.Request, id int) {
	idx := findIndex(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "Team not found",
			fmt.Sprintf("No team with id=%d exists", id))
		return
	}

	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body",
			"Request body must be valid JSON")
		return
	}

	team := teams[idx]

	if v, ok := patch["name"].(string); ok {
		if strings.TrimSpace(v) == "" {
			writeError(w, http.StatusUnprocessableEntity, "Validation failed",
				"'name' cannot be empty")
			return
		}
		team.Name = v
	}
	if v, ok := patch["city"].(string); ok {
		if strings.TrimSpace(v) == "" {
			writeError(w, http.StatusUnprocessableEntity, "Validation failed",
				"'city' cannot be empty")
			return
		}
		team.City = v
	}
	if v, ok := patch["pet"].(string); ok {
		team.Pet = v
	}
	if v, ok := patch["arena"].(string); ok {
		team.Arena = v
	}
	if v, ok := patch["championships"].(float64); ok {
		if v < 0 {
			writeError(w, http.StatusUnprocessableEntity, "Validation failed",
				"'championships' must be >= 0")
			return
		}
		team.Championships = int(v)
	}

	teams[idx] = team

	if err := saveTeams(); err != nil {
		writeError(w, http.StatusInternalServerError, "Persistence error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, teams[idx])
}

func handleDeleteTeam(w http.ResponseWriter, r *http.Request, id int) {
	idx := findIndex(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "Team not found",
			fmt.Sprintf("No team with id=%d exists", id))
		return
	}

	deleted := teams[idx]
	teams = append(teams[:idx], teams[idx+1:]...)

	if err := saveTeams(); err != nil {
		writeError(w, http.StatusInternalServerError, "Persistence error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Message{
		Message: fmt.Sprintf("Team '%s' (id=%d) deleted successfully", deleted.Name, deleted.ID),
	})
}

func validateTeam(t Team, allowZeroChamps bool) string {
	missing := []string{}

	if strings.TrimSpace(t.Name) == "" {
		missing = append(missing, "'name'")
	}
	if strings.TrimSpace(t.City) == "" {
		missing = append(missing, "'city'")
	}
	if strings.TrimSpace(t.Pet) == "" {
		missing = append(missing, "'pet'")
	}
	if strings.TrimSpace(t.Arena) == "" {
		missing = append(missing, "'arena'")
	}
	if !allowZeroChamps && t.Championships < 0 {
		return "'championships' must be >= 0"
	}

	if len(missing) > 0 {
		return "Missing or empty required fields: " + strings.Join(missing, ", ")
	}
	return ""
}

func findIndex(id int) int {
	for i, t := range teams {
		if t.ID == id {
			return i
		}
	}
	return -1
}
