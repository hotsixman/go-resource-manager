package resourcemanager

import (
	"encoding/json"
	util "resource-manager/util"
)

type ResourceObject struct {
	isDirectory        bool
	path               string
	name               string
	userPermissionMap  map[string]([]string)
	groupPermissionMap map[string]([]string)
	childrenMap        map[string](*ResourceObject)
	isLocked           bool
	lockToken          string
}

type ResourceConstructorParam struct {
	IsDirectory        bool
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

	return util.CloneSlice(permissions)
}

/*
해당 그룹이 가진 권한 배열을 반환
*/
func (p ResourceObject) GetGroupPermissions(groupname string) []string {
	permissions := p.groupPermissionMap[groupname]
	if permissions == nil {
		permissions = []string{}
	}

	return util.CloneSlice(permissions)
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
func (p *ResourceObject) AddUserPermission(username string, permission string) []string {
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
		return util.CloneSlice(p.GetUserPermissions(username))
	}

	userPermissions = append(userPermissions, permission)
	p.userPermissionMap[username] = userPermissions
	return util.CloneSlice(p.GetUserPermissions(username))
}

/*
그룹 권한 추가
*/
func (p *ResourceObject) AddGroupPermission(groupname string, permission string) []string {
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
		return util.CloneSlice(p.GetGroupPermissions(groupname))
	}

	groupPermissions = append(groupPermissions, permission)
	p.groupPermissionMap[groupname] = groupPermissions
	return util.CloneSlice(p.GetGroupPermissions(groupname))
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
잠금 여부 반환
*/
func (p ResourceObject) IsLocked() bool {
	return p.isLocked
}

/*
잠금 토큰 반환
*/
func (p ResourceObject) GetLockToken() string {
	return p.lockToken
}

/*
리소스 잠금
  - @return {bool} 잠금 성공 여부
  - @return {string} 잠금 토큰
*/
func (p *ResourceObject) Lock(isDepthInfinity bool, lockToken string) (bool, string) {
	if p.IsLocked() {
		return false, ""
	}

	if lockToken == "" {
		lockToken = util.GenerateRandomString(25)
	}
	if isDepthInfinity {
		for _, child := range p.childrenMap {
			child.Lock(isDepthInfinity, lockToken)
		}

	}
	p.isLocked = true
	p.lockToken = lockToken

	return true, p.lockToken
}

/*
리소스 잠금 해제
*/
func (p *ResourceObject) Unlock(lockToken string) bool {
	if !p.IsLocked() {
		return false
	}

	if lockToken != p.lockToken {
		return false
	}

	p.isLocked = false
	p.lockToken = ""
	return true
}

/*
리소스 강제 잠금 해제
*/
func (p *ResourceObject) UnlockForce() bool {
	if !p.IsLocked() {
		return false
	}

	p.isLocked = false
	p.lockToken = ""
	return true
}

/*
리소스가 폴더(Directory)인지 여부 반환
*/
func (p ResourceObject) IsDirectory() bool {
	return p.isDirectory
}

/*
리소스가 파일인지 여부 반환
*/
func (p ResourceObject) IsFile() bool {
	return !p.isDirectory
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
func (p *ResourceObject) CreateChild(name string, isDirectory bool) (bool, *ResourceObject) {
	if name == "" {
		return false, nil
	}

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
		IsDirectory:        isDirectory,
		Path:               childPath,
		UserPermissionMap:  p.userPermissionMap,
		GroupPermissionMap: p.groupPermissionMap,
	}
	child = NewResourceObject(constructorParam)

	p.childrenMap[name] = child

	return true, child
}

/*
하위 리소스 삭제
*/
func (p *ResourceObject) DeleteChild(name string) bool {
	child := p.GetChild(name)
	if child == nil {
		return false
	}

	for key := range child.childrenMap { // 혹시몰라서 ㅎㅎ
		child.DeleteChild(key)
	}
	delete(p.childrenMap, name)
	return true
}

/*
Map화
*/
func (p ResourceObject) ToMap() map[string]any {
	childrenMap := map[string](map[string]any){}
	for key, value := range p.childrenMap {
		childMap := value.ToMap()
		childrenMap[key] = childMap
	}
	resourceObjectMap := map[string]any{
		"isDirectory":        p.isDirectory,
		"path":               p.path,
		"name":               p.name,
		"userPermissionMap":  p.userPermissionMap,
		"groupPermissionMap": p.groupPermissionMap,
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
		isDirectory:        param.IsDirectory,
		path:               param.Path,
		name:               param.Name,
		userPermissionMap:  param.UserPermissionMap,
		groupPermissionMap: param.GroupPermissionMap,
		childrenMap:        map[string](*ResourceObject){},
		isLocked:           false,
		lockToken:          "",
	}
	return p
}

/*
json으로부터 리소스 객체를 생성
*/
func FromJsonResourceObject(jsonData string) *ResourceObject {
	var resourceObjectMap map[string]any
	json.Unmarshal([]byte(jsonData), &resourceObjectMap)
	return FromMapResourceObject(resourceObjectMap)
}

/*
map으로부터 리소스 객체를 생성
*/
func FromMapResourceObject(resourceObjectMap map[string]any) *ResourceObject {
	isDirectory := resourceObjectMap["isDirectory"].(bool)
	path := resourceObjectMap["path"].(string)
	name := resourceObjectMap["name"].(string)

	userPermissionMap := map[string][]string{}
	for key, value := range resourceObjectMap["userPermissionMap"].(map[string]any) {
		permissions := []string{}
		for _, permission := range value.([]any) {
			permissions = append(permissions, permission.(string))
		}
		userPermissionMap[key] = permissions
	}
	groupPermissionMap := map[string][]string{}
	for key, value := range resourceObjectMap["groupPermissionMap"].(map[string]any) {
		permissions := []string{}
		for _, permission := range value.([]any) {
			permissions = append(permissions, permission.(string))
		}
		groupPermissionMap[key] = permissions
	}

	childrenMap := map[string](*ResourceObject){}
	for key, value := range resourceObjectMap["childrenMap"].(map[string]any) {
		childMap := value.(map[string]any)
		childrenMap[key] = FromMapResourceObject(childMap)
	}

	resourceObject := &ResourceObject{
		isDirectory:        isDirectory,
		path:               path,
		name:               name,
		userPermissionMap:  userPermissionMap,
		groupPermissionMap: groupPermissionMap,
		childrenMap:        childrenMap,
		isLocked:           false,
		lockToken:          "",
	}

	return resourceObject
}
