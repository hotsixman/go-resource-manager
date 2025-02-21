package app

import (
	"app/class"
	"app/util"
	"fmt"
	"net/http"
	"strconv"
)

type ResourceManagerServer struct {
	resourceManager *class.ResourceManager
	mux             *http.ServeMux
}

/*
해당 포트에서 서버 시작
*/
func (s *ResourceManagerServer) Listen(port int) {
	s.mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		username := req.Header.Get("User-Name")
		groupname := req.Header.Get("Group-Name")

		switch req.Method {
		case ("PUT"):
			{
				// 이미 해당 경로에 리소스가 있는지 확인
				if s.resourceManager.GetResourceObject(req.URL.Path) != nil {
					res.WriteHeader(409)
					return
				}

				// 부모 리소스 객체 존재 확인
				var parentResourceObject *class.ResourceObject
				for {
					parentPath, err := util.GetParentDirectory(req.URL.Path)
					if err != nil {
						res.WriteHeader(400)
						return
					}
					parentResourceObject = s.resourceManager.GetResourceObject(parentPath)
					if parentResourceObject != nil || parentPath == "/" {
						break
					}
				}
				if parentResourceObject == nil {
					res.WriteHeader(400)
					return
				}

				// 권한 확인
				hasWritePermission :=
					parentResourceObject.CheckGroupPermission(groupname, "all") ||
						parentResourceObject.CheckGroupPermission(groupname, "write") ||
						parentResourceObject.CheckUserPermission(username, "all") ||
						parentResourceObject.CheckUserPermission(username, "write")
				if !hasWritePermission {
					res.WriteHeader(403)
					return
				}

				// 헤더에서 리소스가 디렉토리인지 아닌지 여부
				isDirectory := false
				if req.Header.Get("Is-Directory") == "true" {
					isDirectory = true
				}

				// 리소스 생성
				success, _ := s.resourceManager.CreateResource(req.URL.Path, isDirectory)
				if success {
					res.WriteHeader(201)
				} else {
					res.WriteHeader(500)
				}
				return
			}
		}
	})

	err := http.ListenAndServe(":"+strconv.Itoa(port), s.mux)
	if err != nil {
		fmt.Println("서버 시작에 오류가 발생했습니다: ", err)
		return
	}
}

/*
ResourceManagerServer 시작
*/
func NewServer(resourceManager *class.ResourceManager) *ResourceManagerServer {
	mux := http.NewServeMux()

	return &ResourceManagerServer{
		resourceManager: resourceManager,
		mux:             mux,
	}
}
