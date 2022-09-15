package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Team struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
	Epics   []string `json:"epics"`
}

var teams = []Team{
	{ID: "T1", Name: "Catalog", Members: []string{"Marcus Allen", "Eamon Scullion", "Kristine Boyd"}, Epics: []string{"1", "2"}},
	{ID: "T2", Name: "Syndication", Members: []string{"Tiago Ramalho", "David Loughridge", "James Nelson"}, Epics: []string{}},
}

func teamIDExists(id string) bool {
	for _, team := range teams {
		if team.ID == id {
			return true
		}
	}
	return false
}

func getTeam(id string) (int, Team) {
	for i, team := range teams {
		if team.ID == id {
			return i, team
		}
	}
	return -1, Team{}
}

func createStubTeam(id string) {
	newTeam := Team{ID: id, Name: "StubTeam-" + id, Members: []string{}, Epics: []string{}}
	teams = append(teams, newTeam)
}

func addEpicToTeam(teamID string, epicID string) {
	if !teamIDExists(teamID) {
		createStubTeam(teamID)
	}

	i, team := getTeam(teamID)
	team.Epics = append(team.Epics, epicID)
	teams[i] = team
}

func getTeams(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teams)
}

func getTeamByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over all teams to find the value matching id
	if i, team := getTeam(id); i == -1 {
		// negative index returned therefore doesn't exist
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Team not found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, team)
	}
}

func postTeams(c *gin.Context) {
	var newTeam Team

	// Call bindJson to bind the received JSON to newTeam
	if err := c.BindJSON(&newTeam); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	if teamIDExists(newTeam.ID) {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "ID Already Exists"})
		return
	}

	// Add the new team to the slice
	teams = append(teams, newTeam)
	c.IndentedJSON(http.StatusCreated, newTeam)
}
