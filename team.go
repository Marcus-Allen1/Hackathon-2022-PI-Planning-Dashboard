package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

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

const (
	XS_WEIGHT  int = 1
	S_WEIGHT       = 2
	M_WEIGHT       = 4
	L_WEIGHT       = 8
	XL_WEIGHT      = 16
	XXL_WEIGHT     = 32
)

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

func removeTeamFromSlice(s []Team, i int) []Team {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeStringFromSlice(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findInSlice(s []string, id string) int {
	for i, item := range s {
		if item == id {
			return i
		}
	}
	return -1
}

func removeEpicFromTeam(teamID string, epicID string) error {
	i, team := getTeam(teamID)

	positionInSlice := findInSlice(team.Epics, epicID)

	if i == -1 {
		return errors.New("Epic Does not Exist for Team")
	}

	team.Epics = removeStringFromSlice(team.Epics, positionInSlice)
	teams[i] = team

	return nil
}

func updateEpicInTeams(oldTeam string, newTeam string, epicID string) error {
	i, oldTeamObject := getTeam(oldTeam)
	j, newTeamObject := getTeam(newTeam)
	if i == -1 || j == -1 {
		return errors.New("Team Does Not Exist")
	}

	positionInSlice := findInSlice(oldTeamObject.Epics, epicID)

	if positionInSlice == -1 {
		return errors.New("Epic not assigned to Team")
	}

	oldTeamObject.Epics = removeStringFromSlice(oldTeamObject.Epics, positionInSlice)

	newTeamObject.Epics = append(newTeamObject.Epics, epicID)

	teams[i] = oldTeamObject
	teams[j] = newTeamObject

	return nil
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

func updateTeam(c *gin.Context) {
	var updatedTeam Team
	id := c.Param("id")

	// Get team if exists
	if i, _ := getTeam(id); i == -1 {
		// negative index returned therefore doesn't exist
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Team not found"})
		return
	} else {
		if err := c.BindJSON(&updatedTeam); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
			return
		}

		if !reflect.DeepEqual(teams[i].Epics, updatedTeam.Epics) {
			updateTeamInEpics(teams[i].Epics, updatedTeam.Epics, id)
		}

		teams[i] = updatedTeam
		c.IndentedJSON(http.StatusOK, updatedTeam)
	}
}

func deleteTeam(c *gin.Context) {
	id := c.Param("id")

	// Get team if exists
	if i, team := getTeam(id); i == -1 {
		// negative index returned therefore doesn't exist
		c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Team %s does not exist so does not require deletion", id)})
		return
	} else {
		// TODO: Remove all Epics related to this team?
		removeTeamFromEpics(team.Epics)
		teams = removeTeamFromSlice(teams, i)
		c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Team %s deleted successfully", id)})
	}
}
