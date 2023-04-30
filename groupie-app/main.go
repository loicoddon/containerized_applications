package main

import (
	"groupie-tracker/controller"
	"groupie-tracker/models"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	PORT := "8080"
	ADDR := "127.0.0.1"

	dbIP := os.Getenv("MYSQL_SERVER_IP")
    dbPort := os.Getenv("MYSQL_SERVER_PORT")
	db_user := os.Getenv("MYSQL_USER")
	db_password := os.Getenv("MYSQL_PASSWORD")
	db_Name := os.Getenv("MYSQL_DATABASE")
	fmt.Println(dbIP,":",dbPort,":",db_user,":",db_password,":",db_Name)

    // Create a MySQL database connection string using the host and port
    connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",db_user,db_password, dbIP, dbPort,db_Name)
	fmt.Printf(connectionString)

    // Open a connection to the database
    db, errr := sql.Open("mysql", connectionString)
    if errr != nil {
        panic(errr.Error())
    }
	// Close connection to the database when main() function die (naturally or not)
    defer db.Close()

	// Check connection to db_server
	ping_err := db.Ping()
    if ping_err != nil {
        fmt.Println("Error: Cannot connect to the database")
        panic(ping_err.Error())
    } else {
        fmt.Println("Successfully connected to the database")
    }

	_, create_err := db.Exec("CREATE TABLE IF NOT EXISTS une_super_table (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, some_text TEXT NOT NULL)")
    if create_err != nil {
        panic(create_err)
    }
	_, test_err := db.Exec("INSERT INTO une_super_table (some_text) VALUES (?)", "hello world")
    if test_err != nil {
        panic(test_err)
    }


	bandsData := &models.ApiData{}
	bandsData.FeedApi()
	bandsData.CreateCaches()

	coords := &models.ApiCoords{}

	staticFiles := http.FileServer(http.Dir("view/"))
	http.Handle("/view/", http.StripPrefix("/view/", staticFiles))

	http.HandleFunc("/", bandsData.RootHandler)

	http.HandleFunc("/map", func(webpage http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			request.ParseForm()
			if len(request.Form["ID"]) != 0 {
				id, _ := strconv.Atoi(request.Form["ID"][0])
				cities,name := bandsData.CitiesTab(id)
				coords.FeedApiCord(cities,name)
				controller.ServeFile(webpage, "map.html", coords)
			}
		} else {
			controller.ServeFile(webpage, "404.html", nil)
		}
	})

	log.Printf("[INFO] - Starting server on http://" + ADDR + ":" + PORT + "/")

	go bandsData.WaitThenRefreshApi()
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("[ERROR] - Server not started properly.\n" + err.Error())
	}
}
