package role

import "forum/pkg/constvar"

func Role2Domain(role string) string {
	switch role {
	case constvar.MuxiRole, constvar.MuxiAdminRole, constvar.SuperAdminRole:
		return constvar.MuxiDomain
	case constvar.NormalRole, constvar.NormalAdminRole:
		return constvar.NormalDomain
	}

	return ""
}
