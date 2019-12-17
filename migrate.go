package trinity

// ToMigrateDB migrate db function
func (t *Trinity) ToMigrateDB(m interface{}) {
	t.db.AutoMigrate(m)
}

// ToCreatePermission create permission func
func (t *Trinity) ToCreatePermission(pSlice []string) {
	p := Permission{Code: pSlice[0], Name: pSlice[1]}
	p.CreateOrInitPermission(t)
}

// MigrateModel to migrate model
func (t *Trinity) MigrateModel(funcToMigrateDB func(interface{}), funcToCreatePermission func([]string), modelToMigrate ...interface{}) {
	for _, v := range modelToMigrate {
		modelName := GetTypeName(v)
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
func (t *Trinity) initEnumtype() {
	createlanguage := "create type language as enum ('zh-CN','en-US');"
	if err := t.db.Exec(createlanguage).Error; err != nil {
		t.logger.Print("Func initEnumtype createlanguage err :" + err.Error())
	}

}
func (t *Trinity) initUserGroup() {
	sql := "CREATE  TABLE \"" + t.setting.Database.TablePrefix + "user_group\" " +
		"( \"id\"  serial ," +
		"\"group_key\" varchar(50) NOT NULL, " +
		"\"user_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.setting.Database.TablePrefix + "user_group_unique_group_key_user_key unique(\"group_key\",\"user_key\")" +
		");"
	if err := t.db.Exec(sql).Error; err != nil {
		t.logger.Print("Func initUserRole err :" + err.Error())
	}
}
func (t *Trinity) initRolePermission() {
	sql := "CREATE  TABLE \"" + t.setting.Database.TablePrefix + "role_permission\" " +
		"( \"id\"  serial ," +
		"\"role_key\" varchar(50) NOT NULL, " +
		"\"permission_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.setting.Database.TablePrefix + "role_permission_unique_role_key_permission_key unique(\"role_key\",\"permission_key\")" +
		");"
	if err := t.db.Exec(sql).Error; err != nil {
		t.logger.Print("Func initRolePermission err :" + err.Error())
	}
}

func (t *Trinity) initUserPermission() {
	sql := "CREATE  TABLE \"" + t.setting.Database.TablePrefix + "user_permission\" " +
		"( \"id\"  serial ," +
		"\"permission_key\" varchar(50) NOT NULL, " +
		"\"user_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.setting.Database.TablePrefix + "user_permission_unique_permission_key_user_key unique(\"permission_key\",\"user_key\")" +
		");"
	if err := t.db.Exec(sql).Error; err != nil {
		t.logger.Print("Func initUserPermission err :" + err.Error())
	}
}
func (t *Trinity) initGroupPermission() {
	sql := "CREATE  TABLE \"" + t.setting.Database.TablePrefix + "group_permission\" " +
		"( \"id\"  serial ," +
		"\"permission_key\" varchar(50) NOT NULL, " +
		"\"group_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.setting.Database.TablePrefix + "user_permission_unique_permission_key_group_key unique(\"permission_key\",\"group_key\")" +
		");"
	if err := t.db.Exec(sql).Error; err != nil {
		t.logger.Print("Func initGroupPermission err :" + err.Error())
	}
}
func (t *Trinity) initGroupRole() {
	sql := "CREATE  TABLE \"" + t.setting.Database.TablePrefix + "group_role\" " +
		"( \"id\"  serial ," +
		"\"group_key\" varchar(50) NOT NULL, " +
		"\"role_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.setting.Database.TablePrefix + "group_role_unique_group_key_role_key unique(\"group_key\",\"role_key\")" +
		");"
	if err := t.db.Exec(sql).Error; err != nil {
		t.logger.Print("Func initGroupRole err :" + err.Error())
	}
}

func (t *Trinity) initUserDefaultValue() {
	initAdminList := []string{"Superadmin"}
	for _, v := range initAdminList {
		var admingroup []Group
		t.db.Where("name = ?", "Superadmin").Find(&admingroup)
		var adminuser User
		adminuser.UserGroup = admingroup
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
		{"system.role.superadmin", "app.right.system.role.superadmin"},
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
		{"Superadmin", "system admin", []string{"system.role.superadmin"}},
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
		// Name        string
		// Description string       `json:"description" gorm:"type:varchar(100);"`
		// Permissions []Permission `json:"permissions" gorm:"many2many:group_permission;AssociationForeignkey:Key;ForeignKey:Key;"`
		// PGroup      *Group       `json:"p_group"  gorm:"AssociationForeignKey:PKey;Foreignkey:Key;"`
		// PKey        string       `json:"p_key"`
		// Roles       []Role
		{"Superadmin", "system.group.Superadmin", []string{"Superadmin"}},
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

func (t *Trinity) migrate() {
	t.GetDB().LogMode(false)
	t.initEnumtype()
	// initial releationship table
	t.initUserGroup()       //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.initUserPermission()  //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.initGroupPermission() //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.initRolePermission()  //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.initGroupRole()       //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.MigrateModel(
		t.ToMigrateDB,
		t.ToCreatePermission,
		&Migration{},
		&Permission{},
		&AppError{},
		&Role{},
		&Group{},
		&User{},
	)
	t.initPermissionDefaultValue()
	t.initRoleDefaultValue()
	t.initGroupDefaultValue()
	t.initUserDefaultValue()
}
