package goresourcemanager

type CONST_PERMISSION struct {
	READ   string
	WRITE  string
	MODIFY string
}

var PERMISSION = CONST_PERMISSION{
	READ:   "read",
	WRITE:  "write",
	MODIFY: "modify",
}
