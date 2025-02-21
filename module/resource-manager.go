package resourcemanager

import (
	"encoding/json"
	"strings"
)

type ResourceManager struct {
	rootResource *ResourceObject
}

// 특정 경로의 리소스 객체의 포인터를 반환
func (m ResourceManager) GetResourceObject(path string) *ResourceObject {
	if path[0] != '/' {
		return nil
	}

	if path == "/" {
		return m.rootResource
	}

	names := strings.Split(path, "/")
	currentResource := m.rootResource
	for i, name := range names {
		if i == 0 {
			continue
		}

		if name == "" {
			return nil
		}

		currentResource = currentResource.GetChild(name)
		if currentResource == nil {
			return nil
		}
	}

	return currentResource
}

/*
경로에 대해 특정 유저가 특정 권한을 가지고 있는지 확인
*/
func (m ResourceManager) CheckUserPermission(path string, username string, permission string) bool {
	resource := m.GetResourceObject(path)

	if resource == nil {
		return false
	}

	return resource.CheckUserPermission(username, permission)
}

/*
경로에 대해 특정 그룹이 특정 권한을 가지고 있는지 확인
*/
func (m ResourceManager) CheckGroupPermission(path string, groupname string, permission string) bool {
	resource := m.GetResourceObject(path)

	return resource.CheckGroupPermission(groupname, permission)
}

func (m *ResourceManager) GetUserPermissions(path string, username string) []string {
	resource := m.GetResourceObject(path)

	return resource.GetUserPermissions(username)
}

/*
경로에 특정 유저의 권한 추가
*/
func (m *ResourceManager) AddUserPermission(path string, username string, permission string) []string {
	resource := m.GetResourceObject(path)

	if resource == nil {
		return []string{}
	}

	return resource.AddUserPermission(username, permission)
}

/*
경로에 특정 그룹의 권한 추가
*/
func (m *ResourceManager) AddGroupPermission(path string, groupname string, permission string) []string {
	resource := m.GetResourceObject(path)

	if resource == nil {
		return []string{}
	}

	return resource.AddGroupPermission(groupname, permission)
}

/*
경로에 특정 유저의 권한 삭제
*/
func (m *ResourceManager) DeleteUserPermission(path string, username string, permission string) {
	resource := m.GetResourceObject(path)

	if resource == nil {
		return
	}

	resource.DeleteUserPermission(username, permission)
}

/*
경로에 특정 그룹의 권한 삭제
*/
func (m *ResourceManager) DeleteGroupPermission(path string, groupname string, permission string) {
	resource := m.GetResourceObject(path)

	if resource == nil {
		return
	}

	resource.DeleteGroupPermission(groupname, permission)
}

/*
경로에 해당하는 리소스 잠금
  - @return {bool} 잠금 성공 여부
  - @return {string} 잠금 토큰
*/
func (m *ResourceManager) Lock(path string, isDepthInfinity bool) (bool, string) {
	resource := m.GetResourceObject(path)
	if resource == nil {
		return false, ""
	}

	return resource.Lock(isDepthInfinity, "")
}

/*
경로에 해당하는 리소스 잠금 해제
*/
func (m *ResourceManager) Unlock(path string, lockToken string) bool {
	resource := m.GetResourceObject(path)
	if resource == nil {
		return false
	}

	return resource.Unlock(lockToken)
}

/*
경로에 해당하는 리소스 강제 잠금 해제
*/
func (m *ResourceManager) UnlockForce(path string) bool {
	resource := m.GetResourceObject(path)
	if resource == nil {
		return false
	}

	return resource.UnlockForce()
}

/*
경로에 리소스 생성
  - @return {bool} 성공 여부
  - @return {*ResourceObject} 생성한 리소스 객체의 포인터, 실패시 nil
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
경로에 리소스 객체 삭제
  - @param {bool} 삭제 성공 여부
*/
func (m *ResourceManager) DeleteResource(path string) bool {
	if path[0] != '/' || path == "/" {
		return false
	}

	names := strings.Split(path, "/")
	currentResource := m.rootResource
	for i, name := range names { // 빈 문자열이 있는 지 검사
		if i == 0 {
			continue
		}

		if name == "" {
			return false
		}
	}
	namesLen := len(names)
	for i, name := range names {
		if i == 0 || i == namesLen-1 {
			continue
		}

		childResource := currentResource.GetChild(name)

		if childResource == nil {
			return false
		}

		currentResource = childResource
	}

	return currentResource.DeleteChild(names[namesLen-1])
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

	rootResource := FromMapResourceObject(resourceManagerMap["rootResource"].(map[string]any))
	m := &ResourceManager{
		rootResource: rootResource,
	}
	return m
}
