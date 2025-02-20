package main

import (
	"fmt"
	resourcemanager "resource-manager/module"
)

func main() {
	resourceManager := resourcemanager.FromJsonResourceManager("{    \"children\": [        {            \"children\": [                {                    \"children\": [],                    \"childrenMap\": {},                    \"groupPermissionMap\": {},                    \"isCollection\": false,                    \"name\": \"sub-test\",                    \"path\": \"/test/sub-test\",                    \"userPermissionMap\": {}                }            ],            \"childrenMap\": {                \"sub-test\": {                    \"children\": [],                    \"childrenMap\": {},                    \"groupPermissionMap\": {},                    \"isCollection\": false,                    \"name\": \"sub-test\",                    \"path\": \"/test/sub-test\",                    \"userPermissionMap\": {}                }            },            \"groupPermissionMap\": {},            \"isCollection\": true,            \"name\": \"test\",            \"path\": \"/test\",            \"userPermissionMap\": {}        }    ],    \"childrenMap\": {        \"test\": {            \"children\": [                {                    \"children\": [],                    \"childrenMap\": {},                    \"groupPermissionMap\": {},                    \"isCollection\": false,                    \"name\": \"sub-test\",                    \"path\": \"/test/sub-test\",                    \"userPermissionMap\": {}                }            ],            \"childrenMap\": {                \"sub-test\": {                    \"children\": [],                    \"childrenMap\": {},                    \"groupPermissionMap\": {},                    \"isCollection\": false,                    \"name\": \"sub-test\",                    \"path\": \"/test/sub-test\",                    \"userPermissionMap\": {}                }            },            \"groupPermissionMap\": {},            \"isCollection\": true,            \"name\": \"test\",            \"path\": \"/test\",            \"userPermissionMap\": {}        }    },    \"groupPermissionMap\": {},    \"isCollection\": true,    \"name\": \"\",    \"path\": \"/\",    \"userPermissionMap\": {        \"hi\": [\"read\"]    }}")

	fmt.Println(resourceManager.ToJson())
}
