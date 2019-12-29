package main

import "encoding/json"

// RoleList - role lise
type RoleList struct {
	Roles []Role
}

// Role - role
type Role struct {
	RoleName string
	Arn      string
}

// RoleMap - get aws roles
func RoleMap() (map[string]string, error) {
	res := make(map[string]string)
	data, err := run("aws", "iam", "list-roles")
	if err != nil {
		return res, err
	}

	var rList RoleList
	if err = json.Unmarshal(data, &rList); err != nil {
		return res, err
	}

	for _, v := range rList.Roles {
		res[v.RoleName] = v.Arn
	}

	return res, nil
}
