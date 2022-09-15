package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Epic struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	DRI          string   `json:"dri"`
	LinksToDocs  []string `json:"linksToDocs"`
	Size         int      `json:"size"`
	CycleTime    float64  `json:"cycleTime"`
	Status       string   `json:"status"`
	PI           string   `json:"pi"`
	Dependencies []string `json:"dependencies"`
	Team         string   `json:"team"`
}

var epics = []Epic{
	{ID: "1", Name: "Epic 1", Description: "Example Description", Type: "CSAT", DRI: "Marcus Allen", LinksToDocs: []string{}, Size: 1, CycleTime: 0, Status: "Pending", PI: "22.2", Dependencies: []string{}, Team: "T1"},
	{ID: "2", Name: "Epic 2", Description: "Example Description", Type: "RTB", DRI: "Marcus Allen", LinksToDocs: []string{}, Size: 3, CycleTime: 0, Status: "Pending", PI: "22.2", Dependencies: []string{}, Team: "T1"},
}

func epicIDExists(id string) bool {
	for _, epic := range epics {
		if epic.ID == id {
			return true
		}
	}
	return false
}

func getEpicsByTeam(c *gin.Context, team string) {
	var foundEpics []Epic

	for _, epic := range epics {
		if epic.Team == team {
			foundEpics = append(foundEpics, epic)
		}
	}
	if foundEpics == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No epics found for Team: %s", team)})
		return
	}
	c.IndentedJSON(http.StatusOK, foundEpics)
}

func getEpics(c *gin.Context) {
	if team := c.Query("team"); team != "" {
		getEpicsByTeam(c, team)
		return
	}

	c.IndentedJSON(http.StatusOK, epics)
}

func getEpicByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over all Epics to find the value matching id
	for _, epic := range epics {
		if epic.ID == id {
			c.IndentedJSON(http.StatusOK, epic)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Epic not found"})
}

func postEpics(c *gin.Context) {
	var newEpic Epic

	// Call bindJSON to bind the received JSON to newEpic
	if err := c.BindJSON(&newEpic); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	if epicIDExists(newEpic.ID) {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "ID Already Exists"})
		return
	}

	// Add the epic to the given team
	addEpicToTeam(newEpic.Team, newEpic.ID)
	// Add the new epic to the slice
	epics = append(epics, newEpic)

	c.IndentedJSON(http.StatusCreated, newEpic)
}
