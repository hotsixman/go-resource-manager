package resourcemanager

import (
	"encoding/json"
	"strings"
)

type ResourceManager struct {
	rootResource *ResourceObject
}

/*
경로에 대해 특정 유저가 특정 권한을 가지고 있는지 확인
*/
func (m ResourceManager) CheckUserPermission(path string, username string, permission string) bool {
	if path[0] != '/' {
		return false
	}

	if path == "/" {
		return m.rootResource.CheckUserPermission(username, permission)
	}

	names := strings.Split(path, "/")
	currentResource := m.rootResource
	for i, name := range names {
		if i == 0 {
			continue
		}

		if name == "" {
			return false
		}

		currentResource = currentResource.GetChild(name)
		if currentResource == nil {
			return false
		}
	}

	return currentResource.CheckUserPermission(username, permission)
}

/*
경로에 대해 특정 그룹이 특정 권한을 가지고 있는지 확인
*/
func (m ResourceManager) CheckGroupPermission(path string, groupname string, permission string) bool {
	if path[0] != '/' {
		return false
	}

	if path == "/" {
		return m.rootResource.CheckGroupPermission(groupname, permission)
	}

	names := strings.Split(path, "/")
	currentResource := m.rootResource
	for i, name := range names {
		if i == 0 {
			continue
		}

		if name == "" {
			return false
		}

		currentResource = currentResource.GetChild(name)
		if currentResource == nil {
			return false
		}
	}

	return currentResource.CheckGroupPermission(groupname, permission)
}

/*
경로에 리소스 생성
*/
func (m *ResourceManager) CreateResource(path string, isDirectory bool) (bool, *ResourceObject) {
	if path[0] != '/' || path == "/" {
		return false, nil
	}

	names := strings.Split(path, "/")
	currentResource := m.rootResource
	for i, name := range names { // 빈 문자열이 있는 지 검사
		if i == 0 {
			continue
		}

		if name == "" {
			return false, nil
		}
	}
	namesLen := len(names)
	for i, name := range names {
		if i == 0 {
			continue
		}

		childResource := currentResource.GetChild(name)

		if childResource == nil {
			var _isDirectory bool
			if i == namesLen-1 {
				_isDirectory = isDirectory
			} else {
				_isDirectory = true
			}

			var result bool
			result, childResource = currentResource.CreateChild(name, _isDirectory)
			if !result {
				return false, nil
			}
		}

		currentResource = childResource
	}

	return true, currentResource
}

/*
Map화
*/
func (m ResourceManager) ToMap() map[string]any {
	resourceManagerMap := map[string]any{
		"rootResource": m.rootResource.ToMap(),
	}
	return resourceManagerMap
}

/*
JSON화
*/
func (m ResourceManager) ToJson() (string, error) {
	jsonData, err := json.Marshal(m.ToMap())

	if err != nil {
		return "", err
	}

	return string(jsonData), err
}

/*
ResourceManager 생성자 함수
*/
func NewResourceManager() *ResourceManager {
	m := &ResourceManager{
		rootResource: NewResourceObject(ResourceConstructorParam{
			Name:               "",
			Path:               "/",
			IsDirectory:        true,
			UserPermissionMap:  map[string][]string{},
			GroupPermissionMap: map[string][]string{},
		}),
	}
	return m
}

/*
JSON으로 ResourceManager 생성
*/
func FromJsonResourceManager(jsonData string) *ResourceManager {
	var resourceManagerMap map[string]any
	json.Unmarshal([]byte(jsonData), &resourceManagerMap)

	rootResource := FromJsonResourceObject(resourceManagerMap["rootResource"].(string))
	m := &ResourceManager{
		rootResource: rootResource,
	}
	return m
}
