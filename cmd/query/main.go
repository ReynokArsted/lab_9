package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "RA"
	password = "postgres"
	dbname   = "sandbox"
)

type DatabaseProvider struct {
	db *sql.DB
}

type Handlers struct {
	DProvider DatabaseProvider
}

func (h *Handlers) SetName(writer echo.Context) error {
	name := writer.QueryParam("name")

	if name == "" {
		return writer.String(400, "Попробуй ввести своё имя через query-параметр 'name'")
	}

	err := h.DProvider.AddName(name)
	if err != nil {
		return writer.String(500, err.Error())
	}

	return writer.String(201, "Имя было изменено на "+name)
}

func (h *Handlers) GetName(writer echo.Context) error {
	answer, err := h.DProvider.SelectName()

	if err != nil {
		return writer.String(500, err.Error())
	}

	return writer.String(200, "Привет, "+answer+"!")
}

func (Db *DatabaseProvider) AddName(name string) error {
	_, err := Db.db.Exec("UPDATE query SET name = ($1)", name)
	if err != nil {
		return err
	}

	return nil
}

func (Db *DatabaseProvider) SelectName() (string, error) {
	var answer string

	row := Db.db.QueryRow("SELECT name FROM query LIMIT 1")
	err := row.Scan(&answer)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	Db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	dp := DatabaseProvider{db: Db}
	h := Handlers{DProvider: dp}

	e := echo.New()
	e.GET("/", h.GetName)
	e.POST("/", h.SetName)
	e.Logger.Fatal(e.Start(":8082"))
}
