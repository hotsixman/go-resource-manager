package goresourcemanager

type ResourceObjectInterface interface {
	GetPermissions(username string) []string
	CheckPermission(username string, permission string) bool
	IsCollection() bool
	IsFile() bool
	GetPath() string
}

type ResourceObject struct {
	isCollection      bool
	path              string
	userPermissionMap map[string]([]string)
	children          []ResourceObject
}

type ResourceConstructorParam struct {
	isCollection      bool
	path              string
	userPermissionMap map[string]([]string)
}

/*
해당 유저가 가진 권한 배열을 반환
*/
func (p ResourceObject) GetPermissions(username string) []string {
	permissions := p.userPermissionMap[username]
	if permissions == nil {
		permissions = []string{}
	}

	return permissions
}

/*
해당 유저가 특정 권한을 가지고 있는지 여부 반환
*/
func (p ResourceObject) CheckPermission(username string, permission string) bool {
	permissions := p.GetPermissions(username)
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
ResourceObject 생성자 함수
*/
func NewPermissionObject(arg ResourceConstructorParam) *ResourceObjectInterface {
	p := ResourceObject{
		isCollection:      arg.isCollection,
		path:              arg.path,
		userPermissionMap: arg.userPermissionMap,
		children:          []ResourceObject{},
	}
	var o ResourceObjectInterface = p
	return &o
}
