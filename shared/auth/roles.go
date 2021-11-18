package auth

import (
	pb "wz2100.net/microlobby/shared/proto/user"
)

const ROLE_SUPERADMIN = "superadmin"
const ROLE_ADMIN = "admin"
const ROLE_USER = "user"
const ROLE_SERVICE = "service"

func HasRole(user *pb.User, role string) bool {
	for _, ur := range user.Roles {
		if ur == role {
			return true
		}
	}

	return false
}

func IntersectsRoles(user *pb.User, roles ...string) bool {
	for _, ur := range user.Roles {
		for _, mr := range roles {
			if ur == mr {
				return true
			}
		}
	}

	return false
}
