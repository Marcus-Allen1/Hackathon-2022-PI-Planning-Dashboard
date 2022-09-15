package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/epics", getEpics)
	router.GET("/epics/:id", getEpicByID)
	router.POST("/epics", postEpics)
	router.PATCH("/epics/:id", updateEpic)
	router.DELETE("/epics/:id", deleteEpic)

	router.GET("/teams", getTeams)
	router.GET("/teams/:id", getTeamByID)
	router.POST("/teams", postTeams)
	router.PATCH("/teams/:id", updateTeam)
	router.DELETE("/teams/:id", deleteTeam)

	router.Run("localhost:8080")
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
