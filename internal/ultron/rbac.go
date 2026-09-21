package ultron

// rolePerms is the static RBAC matrix.
var rolePerms = map[Role][]Perm{
	RoleOwner: {
		PermIDEAccess, PermCompanyRead, PermCompanyWrite,
		PermFleetRead, PermFleetInvoke, PermOSBridge,
		PermUsersRead, PermUsersWrite, PermAuditRead,
	},
	RoleAdmin: {
		PermIDEAccess, PermCompanyRead, PermCompanyWrite,
		PermFleetRead, PermFleetInvoke, PermOSBridge,
		PermUsersRead, PermAuditRead,
	},
	RoleOperator: {
		PermIDEAccess, PermCompanyRead,
		PermFleetRead, PermFleetInvoke, PermOSBridge,
	},
	RoleClient: {
		PermIDEAccess, PermCompanyRead, PermFleetRead,
	},
	RoleViewer: {
		PermIDEAccess, PermCompanyRead, PermFleetRead,
	},
}

// Can reports whether role has perm.
func Can(role Role, perm Perm) bool {
	for _, p := range rolePerms[role] {
		if p == perm {
			return true
		}
	}
	return false
}

// CompanyScoped is true when the role only sees assigned companies.
func CompanyScoped(role Role) bool {
	return role == RoleClient || role == RoleViewer
}

// AllowsCompany checks tenancy for scoped roles.
func AllowsCompany(u User, companyID string) bool {
	if !CompanyScoped(u.Role) {
		return true
	}
	for _, id := range u.CompanyIDs {
		if id == companyID {
			return true
		}
	}
	return false
}

// (p *Plane) Authorize checks login + permission (+ optional company scope).
func (p *Plane) Authorize(token string, perm Perm, companyID string) (User, bool, string) {
	u, ok := p.UserForToken(token)
	if !ok || !u.Active {
		return User{}, false, "unauthorized"
	}
	if !Can(u.Role, perm) {
		p.Audit(u.Email, "deny."+string(perm), companyID, false, "rbac")
		return u, false, "forbidden"
	}
	if companyID != "" && !AllowsCompany(u, companyID) {
		p.Audit(u.Email, "deny.company", companyID, false, "tenancy")
		return u, false, "forbidden company"
	}
	return u, true, ""
}
