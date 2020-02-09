package trinity

import "fmt"

// ToMigrateDB migrate db function
func (t *Trinity) ToMigrateDB(m interface{}) {
	err := t.db.AutoMigrate(m).Error
	if err != nil {
		fmt.Println(err)
	}
}

// ToCreatePermission create permission func
func (t *Trinity) ToCreatePermission(pSlice []string) {
	p := Permission{Code: pSlice[0], Name: pSlice[1]}
	p.CreateOrInitPermission(t)
}

// MigrateModel to migrate model
func (t *Trinity) MigrateModel(funcToMigrateDB func(interface{}), funcToCreatePermission func([]string), modelToMigrate ...interface{}) {
	for _, v := range modelToMigrate {
		modelName := GetTypeName(v, true)
		funcToMigrateDB(v)
		PermissionList := [][]string{
			{"system.add." + modelName, "app.right.system.add." + modelName},
			{"system.view." + modelName, "app.right.system.view." + modelName},
			{"system.edit." + modelName, "app.right.system.edit." + modelName},
			{"system.delete." + modelName, "app.right.system.delete." + modelName},
		}
		for _, v := range PermissionList {
			funcToCreatePermission(v)
		}
	}
}

func (t *Trinity) initUserDefaultValue() {
	initAdminList := []string{"superadmin"}
	for _, v := range initAdminList {
		var admingroup []Group
		t.db.Where("name = ?", "superadmin").Find(&admingroup)
		var adminuser User
		adminuser.Groups = admingroup
		adminuser.Username = v
		if err := t.db.FirstOrCreate(&adminuser, map[string]interface{}{"username": v}).Error; err != nil {
			LogPrint("Init Admin user  err :" + err.Error())
		}
	}

}

func (t *Trinity) initPermissionDefaultValue() {
	//naming rule
	// Code : module.type.code
	// Desc : app.right.app.type.xxx
	PermissionList := [][]string{
		{"system.permission.superadmin", "app.right.system.role.superadmin"},
	}
	for _, v := range PermissionList {
		t.db.Where(Permission{Code: v[0]}).FirstOrCreate(&Permission{Code: v[0], Name: v[1]})
		t.db.Where(Permission{Code: v[0]}).First(&Permission{}).Update(map[string]interface{}{"name": v[1]})
	}
}

func (t *Trinity) initRoleDefaultValue() {
	roleDefaultList := [][]interface{}{
		//naming rule
		// Code : module.type.code
		// Desc : app.right.app.type.xxx
		{"system.role.superadmin", "system admin", []string{"system.permission.superadmin"}},
	}
	for _, v := range roleDefaultList {
		var rolepermission []Permission
		name, _ := v[0].(string)
		description, _ := v[1].(string)
		plist, _ := v[2].([]string)
		t.db.Where("code in (?)", plist).Find(&rolepermission)
		var role Role

		role.Name = name
		role.Description = description
		role.Permissions = rolepermission

		if err := t.db.Where(Role{Name: name}).FirstOrCreate(&role).Error; err != nil {
			LogPrint("Func initial role data  err :" + err.Error())
		}
		if err := t.db.Model(&role).Updates(map[string]interface{}{"description": description}).Error; err != nil {
			LogPrint("Func update  role   err :" + err.Error())
		}
		if err := t.db.Model(&role).Association("Permissions").Replace(role.Permissions).Error; err != nil {
			LogPrint("Func update  role permission  err :" + err.Error())
		}

	}
}

func (t *Trinity) initGroupDefaultValue() {
	GroupDefaultList := [][]interface{}{
		{"system.group.superadmin", "system.group.superadmin", []string{"system.role.superadmin"}},
	}
	for _, v := range GroupDefaultList {
		var roleList []Role
		gName, _ := v[0].(string)
		gDesc, _ := v[1].(string)
		rlist, _ := v[2].([]string)
		t.db.Where("name in (?)", rlist).Find(&roleList)
		var group Group
		group.Name = gName
		group.Roles = roleList
		if err := t.db.Where(Group{Name: gName}).FirstOrCreate(&group).Error; err != nil {
			LogPrint("Func initial group data  err :" + err.Error())
		}
		if err := t.db.Model(&group).Updates(map[string]interface{}{"description": gDesc}).Error; err != nil {
			LogPrint("Func update  group   err :" + err.Error())
		}
		if err := t.db.Model(&group).Association("Roles").Replace(group.Roles).Error; err != nil {
			LogPrint("Func update  group roles  err :" + err.Error())
		}

	}
}
