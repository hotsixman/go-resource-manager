package main

import (
	"os"
	resourcemanager "resource-manager/module"
)

func main() {
	resourceManager := resourcemanager.NewResourceManager()
	resourceManager.CreateResource("/foo/bar.txt", false)
	resourceManager.DeleteResource("/foo/bar.txt")
	resourceManager.CreateResource("/test/asjdwd/ansdwd.txt", false)
	resourceManager.DeleteResource("/test/asjdwd")
	json, _ := resourceManager.ToJson()
	os.WriteFile("./test.json", []byte(json), 0644)
}
