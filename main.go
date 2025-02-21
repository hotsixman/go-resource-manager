package main

import (
	"fmt"
	"os"
	resourcemanager "resource-manager/module"
)

func main() {
	json, err := os.ReadFile("./test.json")
	if err != nil {
		return
	}
	resourceManager := resourcemanager.FromJsonResourceManager(string(json))
	resourceManager.AddUserPermission("/test", "snom", "read")
	resourceManager.CreateResource("/test/ass.txt", false)
	fmt.Println(resourceManager.ToJson())
}
