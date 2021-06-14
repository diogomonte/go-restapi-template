
package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

type App struct {
	Router	*mux.Router
	Db		*gorm.DB
	ApiKey	string
}

type Response struct {
	message string
}

func (app *App) Initialize(apikey string, mysqlUrl string) {
	app.Router = mux.NewRouter().StrictSlash(true)

	app.ApiKey = apikey
	app.Router.Use(app.authenticate)

	app.Router.HandleFunc("/helloworld", app.handleHelloWorld).Methods(http.MethodPost)

	db, err := gorm.Open("mysql", mysqlUrl)
	if err != nil {
		log.Fatal("failed to connect to database. ", err.Error())
	}

	app.Db = db
}

func (app *App) handleHelloWorld(res http.ResponseWriter, _ *http.Request) {

	response := Response{
		message: "hello world!",
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if err:= json.NewEncoder(res).Encode(response); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (app *App) Run(port string) {
	c := cors.New(cors.Options{
		AllowedHeaders: []string{ "Authorization" },
		AllowCredentials: true,
	})
	handler := c.Handler(app.Router)
	log.Fatal(http.ListenAndServe(":" + port, handler))
}

func getEnvOrDefault(key string, defaultValue string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}
	return env
}

func main() {
	log.Println("-- Running template --")
	app := App{}
	app.Initialize(
		getEnvOrDefault("API_KEY", ""),
		getEnvOrDefault("MYSQL_URL", "root:root@tcp(mysql:3306)/my_db?charset=utf8&parseTime=True&loc=Local"))
	app.Run(getEnvOrDefault("PORT", "8082"))
}
