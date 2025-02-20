package resourcemanager

import "encoding/json"

type ResourceManager struct {
	rootObject *ResourceObject
}

func (r ResourceManager) ToJson() (string, error) {
	jsonData, err := json.Marshal(r.ToMap())

	if err != nil {
		return "", err
	}

	return string(jsonData), err
}

func (r ResourceManager) ToMap() map[string]any {
	managerMap := map[string]any{
		"rootObject": r.rootObject.ToMap(),
	}
	return managerMap
}

func NewResourceManager() *ResourceManager {
	m := &ResourceManager{
		rootObject: NewResourceObject(ResourceConstructorParam{
			Name:               "",
			Path:               "/",
			IsCollection:       true,
			UserPermissionMap:  map[string][]string{},
			GroupPermissionMap: map[string][]string{},
		}),
	}
	return m
}

func FromJsonResourceManager(rootJsondata string) *ResourceManager {
	m := &ResourceManager{
		rootObject: FromJsonResourceObject(rootJsondata),
	}
	return m
}
