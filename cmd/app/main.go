package main

import "github.com/Skrwbrgle/go-backend-learn/internal/app"


func main() {
	application := app.NewApp()
	application.Run(":8080")
}
