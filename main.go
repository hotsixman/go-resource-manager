package main

import (
	class "app/class"
	server "app/server"
	"os"
)

func main() {
	json, err := os.ReadFile("./test.json")
	if err != nil {
		return
	}
	resourceManager := class.FromJsonResourceManager(string(json))
	resourceManager.AddUserPermission("/test", "snom", "read")
	resourceManager.CreateResource("/test/ass.txt", false)
	resourceManagerServer := server.NewServer(resourceManager)
	resourceManagerServer.Listen(3000)
}
