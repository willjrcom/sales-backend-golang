package rbac

import (
	"fmt"
)

// Role representa uma função em um sistema RBAC
type Role string

// User representa um usuário em um sistema RBAC
type User struct {
	Name  string
	Roles []Role
}

// Resource representa um recurso em um sistema RBAC
type Resource string

// AccessLevel representa os diferentes níveis de acesso a um recurso
type AccessLevel int

// Constantes para representar diferentes níveis de acesso
const (
	NoAccess AccessLevel = iota
	Read
	Write
	Admin
)

// RBAC representa o sistema de controle de acesso baseado em funções
type RBAC struct {
	permissions map[Role]map[Resource]AccessLevel
}

// NewRBAC cria uma nova instância do sistema RBAC
func NewRBAC() *RBAC {
	return &RBAC{
		permissions: make(map[Role]map[Resource]AccessLevel),
	}
}

// AddRolePermission adiciona permissões para uma função específica
func (r *RBAC) AddRolePermission(role Role, resource Resource, level AccessLevel) {
	if r.permissions[role] == nil {
		r.permissions[role] = make(map[Resource]AccessLevel)
	}
	r.permissions[role][resource] = level
}

// CanAccess verifica se um usuário tem permissão para acessar um recurso específico
func (r *RBAC) CanAccess(user User, resource Resource) bool {
	for _, role := range user.Roles {
		if level, ok := r.permissions[role][resource]; ok && level > NoAccess {
			return true
		}
	}

	return false
}

func Main() {
	// Criar uma instância RBAC
	rbac := NewRBAC()

	// Definir permissões
	rbac.AddRolePermission("admin", "document", Admin)
	rbac.AddRolePermission("editor", "document", Write)
	rbac.AddRolePermission("reader", "document", Read)

	// Criar usuários com diferentes funções
	adminUser := User{Name: "Admin User", Roles: []Role{"admin"}}
	editorUser := User{Name: "Editor User", Roles: []Role{"editor"}}
	readerUser := User{Name: "Reader User", Roles: []Role{"reader"}}

	// Verificar permissões
	fmt.Println("Admin can access document:", rbac.CanAccess(adminUser, "document"))
	fmt.Println("Editor can access document:", rbac.CanAccess(editorUser, "document"))
	fmt.Println("Reader can access document:", rbac.CanAccess(readerUser, "document"))
}
