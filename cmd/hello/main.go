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

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetHello(writer echo.Context) error {
	msg, err := h.dbProvider.SelectHello()
	if err != nil {
		return writer.String(500, err.Error())
	}

	writer.String(200, msg)
	return nil
}
func (h *Handlers) PostHello(writer echo.Context) error {
	input := struct {
		Msg string `json:"msg"`
	}{}

	err := writer.Bind(&input)
	if err != nil {
		return writer.String(400, err.Error())
	}

	err = h.dbProvider.InsertHello(input.Msg)
	if err != nil {
		return writer.String(500, err.Error())
	}

	return writer.String(201, "Имя "+input.Msg+" было добавлено в БД")
}

// Методы для работы с базой данных
func (dp *DatabaseProvider) SelectHello() (string, error) {
	var msg string

	// Получаем одно сообщение из таблицы hello, отсортированной в случайном порядке
	row := dp.db.QueryRow("SELECT message FROM hello ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}

	return msg, nil
}
func (dp *DatabaseProvider) InsertHello(msg string) error {
	_, err := dp.db.Exec("INSERT INTO hello (message) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Создание соединения с сервером postgres
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем провайдер для БД с набором методов
	dp := DatabaseProvider{db: db}
	// Создаем экземпляр структуры с набором обработчиков
	h := Handlers{dbProvider: dp}

	e := echo.New()
	e.GET("/", h.GetHello)
	e.POST("/", h.PostHello)
	e.Logger.Fatal(e.Start(":8081"))
}
