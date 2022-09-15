package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Welcome struct {
	Name string
	Time string
}

// // API Example
// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// // albums slice to seed record album data.
// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

func main() {
	router := gin.Default()
	router.GET("/epics", getEpics)
	router.GET("/epics/:id", getEpicByID)
	router.POST("/epics", postEpics)

	router.GET("/teams", getTeams)
	router.GET("/teams/:id", getTeamByID)
	router.POST("/teams", postTeams)

	router.Run("localhost:8080")
}

// getAlbums responsds with the list of all albums as JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// postAlbums adds an album from JSON received in request body
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// func main() {
// 	welcome := Welcome{"Anonymous", time.Now().Format(time.Stamp)}

// 	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))

// 	http.Handle("/static/",
// 		http.StripPrefix("/static/",
// 			http.FileServer(http.Dir("static"))))

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if name := r.FormValue("name"); name != "" {
// 			welcome.Name = name
// 		}

// 		if err := templates.ExecuteTemplate(w, "welcome-template.html", welcome); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	})

// 	fmt.Println("Listening")
// 	fmt.Println(http.ListenAndServe(":8080", nil))
// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("<h1>Hello World!</h1>"))
// }

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "3000"
// 	}

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", indexHandler)
// 	http.ListenAndServe(":"+port, mux)
// }
