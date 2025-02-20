package main

import (
	"fmt"
	resourcemanager "resource-manager/module"
)

func main() {
	resource := resourcemanager.NewResourceObject(resourcemanager.ResourceConstructorParam{
		IsCollection:       true,
		Name:               "",
		Path:               "/",
		UserPermissionMap:  map[string]([]string){},
		GroupPermissionMap: map[string]([]string){},
	})

	resource.CreateChild("test", true)
	resource.GetChild("test").CreateChild("sub-test", false)

	jsonData, _ := resource.ToJSON()
	fmt.Println(jsonData)
}
