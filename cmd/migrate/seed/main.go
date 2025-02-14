package main

import (
	"log"
	"spy-cat-agency/internal/db"
	"spy-cat-agency/internal/store"
)

func main() {
	addr := "postgres://admin:adminpassword@localhost/agency?sslmode=disable"

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)

	log.Println("DB was seeded")
}
