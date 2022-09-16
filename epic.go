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

func getEpic(id string) (int, Epic) {
	for i, epic := range epics {
		if epic.ID == id {
			return i, epic
		}
	}
	return -1, Epic{}
}

func removeEpicFromSlice(s []Epic, i int) []Epic {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func differenceInSlices(slice1, slice2 []string) []string {
	var diff []string

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// string not found, we add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

func removeTeamFromEpics(removeFromList []string) {
	for _, epicID := range removeFromList {
		i, epic := getEpic(epicID)
		epic.Team = ""
		epics[i] = epic
	}
}

func updateTeamInEpics(oldEpics []string, newEpics []string, teamID string) {
	removed := differenceInSlices(oldEpics, newEpics)
	added := differenceInSlices(newEpics, oldEpics)

	for _, epicID := range removed {
		i, epic := getEpic(epicID)
		epic.Team = ""
		epics[i] = epic
	}

	for _, epicID := range added {
		i, epic := getEpic(epicID)
		epic.Team = teamID
		epics[i] = epic
	}
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
	if i, epic := getEpic(id); i == -1 {
		// negative index returned therefore doesn't exist
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Epic not found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, epic)
	}
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

func updateEpic(c *gin.Context) {
	var updatedEpic Epic
	id := c.Param("id")

	// Get epic if exists
	if i, _ := getEpic(id); i == -1 {
		// negative index returned therefore doesn't exist
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Epic not found"})
		return
	} else {
		if err := c.BindJSON(&updatedEpic); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
			return
		}

		if epics[i].Team != updatedEpic.Team {
			err := updateEpicInTeams(epics[i].Team, updatedEpic.Team, id)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
		}

		epics[i] = updatedEpic
		c.IndentedJSON(http.StatusOK, updatedEpic)
	}
}

func deleteEpic(c *gin.Context) {
	id := c.Param("id")

	// Get epic if exists
	if i, epic := getEpic(id); i == -1 {
		// negative index returned therefore doesn't exist
		c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Epic %s does not exist so does not require deletion", id)})
		return
	} else {
		if err := removeEpicFromTeam(epic.Team, epic.ID); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		epics = removeEpicFromSlice(epics, i)
		c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Epic %s deleted successfully", id)})
	}
}
