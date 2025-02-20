package resourcemanager

import (
	"encoding/json"
)

type ResourceObject struct {
	isCollection       bool
	path               string
	name               string
	userPermissionMap  map[string]([]string)
	groupPermissionMap map[string]([]string)
	children           [](*ResourceObject)
	childrenMap        map[string](*ResourceObject)
}

type ResourceConstructorParam struct {
	IsCollection       bool
	Path               string
	Name               string
	UserPermissionMap  map[string]([]string)
	GroupPermissionMap map[string]([]string)
}

/*
해당 유저가 가진 권한 배열을 반환
*/
func (p ResourceObject) GetUserPermissions(username string) []string {
	permissions := p.userPermissionMap[username]
	if permissions == nil {
		permissions = []string{}
	}

	return permissions
}

/*
해당 그룹이 가진 권한 배열을 반환
*/
func (p ResourceObject) GetGroupPermissions(groupname string) []string {
	permissions := p.groupPermissionMap[groupname]
	if permissions == nil {
		permissions = []string{}
	}

	return permissions
}

/*
해당 유저가 특정 권한을 가지고 있는지 여부 반환
*/
func (p ResourceObject) CheckUserPermission(username string, permission string) bool {
	permissions := p.GetUserPermissions(username)
	hasPermission := false
	for _, v := range permissions {
		if v == permission {
			hasPermission = true
			break
		}
	}
	return hasPermission
}

/*
해당 그룹이 특정 권한을 가지고 있는지 여부 반환
*/
func (p ResourceObject) CheckGroupPermission(groupname string, permission string) bool {
	permissions := p.GetGroupPermissions(groupname)
	hasPermission := false
	for _, v := range permissions {
		if v == permission {
			hasPermission = true
			break
		}
	}
	return hasPermission
}

/*
유저 권한 추가
*/
func (p *ResourceObject) AddUserPermission(username string, permission string) bool {
	isValid := false
	for _, v := range PERMISSIONS {
		if v == permission {
			isValid = true
			break
		}
	}
	if !isValid {
		return false
	}

	userPermissions := p.userPermissionMap[username]
	if userPermissions == nil {
		userPermissions = []string{}
	}

	alreadyHasPermission := false
	for _, v := range userPermissions {
		if v == permission {
			alreadyHasPermission = true
		}
	}
	if alreadyHasPermission {
		return true
	}

	userPermissions = append(userPermissions, permission)
	p.userPermissionMap[username] = userPermissions
	return true
}

/*
그룹 권한 추가
*/
func (p *ResourceObject) AddGroupPermission(groupname string, permission string) bool {
	isValid := false
	for _, v := range PERMISSIONS {
		if v == permission {
			isValid = true
			break
		}
	}
	if !isValid {
		return false
	}

	groupPermissions := p.groupPermissionMap[groupname]
	if groupPermissions == nil {
		groupPermissions = []string{}
	}

	alreadyHasPermission := false
	for _, v := range groupPermissions {
		if v == permission {
			alreadyHasPermission = true
		}
	}
	if alreadyHasPermission {
		return true
	}

	groupPermissions = append(groupPermissions, permission)
	p.groupPermissionMap[groupname] = groupPermissions
	return true
}

/*
유저 권한 제거
*/
func (p *ResourceObject) DeleteUserPermission(username string, permission string) {
	userPermissions := p.userPermissionMap[username]
	if userPermissions == nil {
		return
	}

	permissionIndex := -1
	for i, v := range userPermissions {
		if v == permission {
			permissionIndex = i
			break
		}
	}
	if permissionIndex >= 0 {
		userPermissions = append(userPermissions[:permissionIndex], userPermissions[permissionIndex+1:]...)
		p.userPermissionMap[username] = userPermissions
	}
}

/*
그룹 권한 제거
*/
func (p *ResourceObject) DeleteGroupPermission(groupname string, permission string) {
	groupPermissions := p.groupPermissionMap[groupname]
	if groupPermissions == nil {
		return
	}

	permissionIndex := -1
	for i, v := range groupPermissions {
		if v == permission {
			permissionIndex = i
			break
		}
	}
	if permissionIndex >= 0 {
		groupPermissions = append(groupPermissions[:permissionIndex], groupPermissions[permissionIndex+1:]...)
		p.groupPermissionMap[groupname] = groupPermissions
	}
}

/*
리소스가 폴더(collection)인지 여부 반환
*/
func (p ResourceObject) IsCollection() bool {
	return p.isCollection
}

/*
리소스가 파일인지 여부 반환
*/
func (p ResourceObject) IsFile() bool {
	return !p.isCollection
}

/*
리소스의 경로 반환
*/
func (p ResourceObject) GetPath() string {
	return p.path
}

/*
하위 리소스 반환
값이 없을 수 있으니 `nil`인지 확인할 것
*/
func (p ResourceObject) GetChild(name string) *ResourceObject {
	return p.childrenMap[name]
}

/*
하위 리소스 생성

- @return {bool} 성공 여부

- @return {*ResourceObject} 생성한 하위 리소스, 성공 여부가 false이면 nil
*/
func (p *ResourceObject) CreateChild(name string, isCollection bool) (bool, *ResourceObject) {
	child := p.GetChild(name)
	if child != nil {
		return false, nil
	}

	childPath := ""
	if p.GetPath() == "/" {
		childPath = "/" + name
	} else {
		childPath = p.GetPath() + "/" + name
	}

	constructorParam := ResourceConstructorParam{
		Name:               name,
		IsCollection:       isCollection,
		Path:               childPath,
		UserPermissionMap:  p.userPermissionMap,
		GroupPermissionMap: p.groupPermissionMap,
	}
	child = NewResourceObject(constructorParam)

	p.children = append(p.children, child)
	p.childrenMap[name] = child

	return true, child
}

/*
하위 리소스 삭제
*/
func (p *ResourceObject) RemoveChild(name string) bool {
	child := p.GetChild(name)
	if child == nil {
		return false
	}

	delete(p.childrenMap, name)
	childIndex := -1
	for i, v := range p.children {
		if v.name == name {
			childIndex = i
			break
		}
	}
	if childIndex >= 0 {
		p.children = append(p.children[:childIndex], p.children[childIndex+1:]...)
	}
	return true
}

/*
Map화
*/
func (p ResourceObject) ToMap() map[string]any {
	childrenMap := map[string](map[string]any){}
	children := []map[string]any{}
	for key, value := range p.childrenMap {
		childMap := value.ToMap()
		childrenMap[key] = childMap
	}
	for _, value := range p.children {
		if childrenMap[value.name] != nil {
			children = append(children, childrenMap[value.name])
		} else {
			children = append(children, value.ToMap())
		}
	}
	resourceObjectMap := map[string]any{
		"isCollection":       p.isCollection,
		"path":               p.path,
		"name":               p.name,
		"userPermissionMap":  p.userPermissionMap,
		"groupPermissionMap": p.groupPermissionMap,
		"children":           children,
		"childrenMap":        childrenMap,
	}
	return resourceObjectMap
}

/*
JSON화
*/
func (p ResourceObject) ToJSON() (string, any) {
	jsonData, err := json.Marshal(p.ToMap())

	if err != nil {
		return "", err
	}

	return string(jsonData), err
}

/*
ResourceObject 생성자 함수
*/
func NewResourceObject(param ResourceConstructorParam) *ResourceObject {
	p := &ResourceObject{
		isCollection:       param.IsCollection,
		path:               param.Path,
		name:               param.Name,
		userPermissionMap:  param.UserPermissionMap,
		groupPermissionMap: param.GroupPermissionMap,
		children:           []*ResourceObject{},
		childrenMap:        map[string]*ResourceObject{},
	}
	return p
}
