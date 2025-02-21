package main

import (
	"fmt"
	resourcemanager "resource-manager/module"
)

func main() {
	resourceManager := resourcemanager.NewResourceManager()
	resourceManager.CreateResource("/foo/bar.txt", false)
	fmt.Println(resourceManager.ToJson())
}
