package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
)

type User struct {
	ID int `json:"id"`
	FirstName string  `json:"firstName"`
	Lastname string  `json:"lastName"`
	UserName string  `json:"userName"`
	Description string `json:"description"`
	IsActive bool `json:"isActive"`
}

func main(){
	fmt.Println("Hello World")
	app := fiber.New()

	users := []User{}

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello world"})
	})

	app.Post("/api/users", func(c fiber.Ctx) error {
		user := User{}


		if err := c.Bind().Body(user); err != nil {		return err
	 }

	 if user.UserName == "" {
		return c.JSON(fiber.Map{"error":"No username found."})
	 }

	 user.ID = len(users) + 1
	 users = append(users, user)


	 return c.JSON(fiber.Map{"message": "User created successfully!", "user": user})
	})


	log.Fatal(app.Listen(":4000"))
}

