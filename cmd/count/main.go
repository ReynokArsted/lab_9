package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

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

type DatabaseProvider struct { // Структура с полем, которое хранит ссылку на СУБД
	db *sql.DB
}

type Handlers struct {
	DProvider DatabaseProvider
}

func (h *Handlers) GetCount(writer echo.Context) error {
	answer, err := h.DProvider.SelectCount()
	if err != nil {
		return writer.String(500, err.Error())
	}

	return writer.String(200, "Значение счётчика: "+answer)
}

func (h *Handlers) SetCount(writer echo.Context) error {
	input := struct {
		Msg string `json:msg`
	}{}

	err := writer.Bind(&input)
	if err != nil {
		return writer.String(400, err.Error())
	}

	value, err := strconv.Atoi(input.Msg)
	if err != nil {
		return writer.String(400, "Было введено не число или присутствуют пробелы в записи числа")
	}

	err = h.DProvider.UpdateCount(value)
	if err != nil {
		return writer.String(500, err.Error())
	}

	if value > 0 {
		return writer.String(201, "Значение счётчика было изменено на +"+strconv.Itoa(value))
	} else {
		return writer.String(201, "Значение счётчика было изменено на "+strconv.Itoa(value))
	}
}

func (Dp *DatabaseProvider) SelectCount() (string, error) {
	var dbAnswer string

	row := Dp.db.QueryRow("SELECT value FROM count")
	err := row.Scan(&dbAnswer) // Проверка на то, есть ли искомые данные в БД
	if err != nil {
		return "", err
	}

	return dbAnswer, nil
}

func (Dp *DatabaseProvider) UpdateCount(n int) error {
	_, err := Dp.db.Exec("UPDATE count SET value = value + ($1)", n)
	if err != nil {
		return err
	}

	return nil
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
	e.GET("/", h.GetCount)
	e.POST("/", h.SetCount)
	e.Logger.Fatal(e.Start(":8083"))
}
