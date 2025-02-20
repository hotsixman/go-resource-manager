package resourcemanager

var PERMISSIONS = [3]string{
	/*
		파일 내용을 읽을 수 있음.
		폴더의 하위 파일, 폴더 목록을 볼 수 있음.
	*/
	"read",
	/*
		폴더에 하위 폴더, 파일을 생성할 수 있음.
		파일의 내용을 변경할 수 있음.
	*/
	"write",
	/*
		파일, 폴더를 삭제할 수 있음.
		파일, 폴더의 이름을 변경할 수 있음.
	*/
	"modify",
}
