package trinity

// ToMigrateDB migrate db function
func ToMigrateDB(m interface{}) {
	Db.AutoMigrate(m)
}

// ToCreatePermission create permission func
func ToCreatePermission(pSlice []string) {
	p := Permission{Code: pSlice[0], Name: pSlice[1]}
	p.CreateOrInitPermission()
}

// MigrateModel to migrate model
func MigrateModel(funcToMigrateDB func(interface{}), funcToCreatePermission func([]string), modelToMigrate ...interface{}) {
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
	if err := t.Db.Exec(createlanguage).Error; err != nil {
		t.Logger.Print("Func initEnumtype createlanguage err :" + err.Error())
	}

}
func (t *Trinity) initUserGroup() {
	sql := "CREATE  TABLE \"" + t.Setting.Database.TablePrefix + "user_group\" " +
		"( \"id\"  serial ," +
		"\"group_key\" varchar(50) NOT NULL, " +
		"\"user_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.Setting.Database.TablePrefix + "user_group_unique_group_key_user_key unique(\"group_key\",\"user_key\")" +
		");"
	if err := t.Db.Exec(sql).Error; err != nil {
		t.Logger.Print("Func initUserRole err :" + err.Error())
	}
}

func (t *Trinity) initUserPermission() {
	sql := "CREATE  TABLE \"" + t.Setting.Database.TablePrefix + "user_permission\" " +
		"( \"id\"  serial ," +
		"\"permission_key\" varchar(50) NOT NULL, " +
		"\"user_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.Setting.Database.TablePrefix + "user_permission_unique_permission_key_user_key unique(\"permission_key\",\"user_key\")" +
		");"
	if err := t.Db.Exec(sql).Error; err != nil {
		t.Logger.Print("Func initUserPermission err :" + err.Error())
	}
}
func (t *Trinity) initGroupPermission() {
	sql := "CREATE  TABLE \"" + t.Setting.Database.TablePrefix + "group_permission\" " +
		"( \"id\"  serial ," +
		"\"permission_key\" varchar(50) NOT NULL, " +
		"\"group_key\" varchar(50) NOT NULL, " +
		"PRIMARY KEY (\"id\") ," +
		"constraint " + t.Setting.Database.TablePrefix + "user_permission_unique_permission_key_group_key unique(\"permission_key\",\"group_key\")" +
		");"
	if err := t.Db.Exec(sql).Error; err != nil {
		t.Logger.Print("Func initGroupPermission err :" + err.Error())
	}
}

func initUserDefaultValue() {
	initAdminList := []string{"Superadmin"}
	for _, v := range initAdminList {
		var admingroup []Group
		Db.Where("name = ?", "Superadmin").Find(&admingroup)
		var adminuser User
		adminuser.UserGroup = admingroup
		adminuser.Username = v
		if err := Db.FirstOrCreate(&adminuser, map[string]interface{}{"username": v}).Error; err != nil {
			LogPrint("Init Admin user  err :" + err.Error())
		}
	}

}

func initPermissionDefaultValue() {
	//naming rule
	// Code : module.type.code
	// Desc : app.right.app.type.xxx
	PermissionList := [][]string{
		{"system.role.superadmin", "app.right.system.role.superadmin"},
	}
	for _, v := range PermissionList {
		Db.Where(Permission{Code: v[0]}).FirstOrCreate(&Permission{Code: v[0], Name: v[1]})
		Db.Where(Permission{Code: v[0]}).First(&Permission{}).Update(map[string]interface{}{"name": v[1]})
	}
}

func initGroupDefaultValue() {
	GroupDefaultList := [][]interface{}{
		//naming rule
		// Code : module.type.code
		// Desc : app.right.app.type.xxx
		{"Superadmin", "system admin", []string{"system.role.superadmin"}},
	}
	for _, v := range GroupDefaultList {
		var rolepermission []Permission
		name, _ := v[0].(string)
		description, _ := v[1].(string)
		plist, _ := v[2].([]string)
		Db.Where("code in (?)", plist).Find(&rolepermission)
		var group Group

		group.Name = name
		group.Description = description
		group.GroupPermission = rolepermission

		if err := Db.Where(Group{Name: name}).FirstOrCreate(&group).Error; err != nil {
			LogPrint("Func initial group data  err :" + err.Error())
		}
		if err := Db.Model(&group).Updates(map[string]interface{}{"description": description}).Error; err != nil {
			LogPrint("Func update  group   err :" + err.Error())
		}
		if err := Db.Model(&group).Association("GroupPermission").Replace(group.GroupPermission).Error; err != nil {
			LogPrint("Func update  group permission  err :" + err.Error())
		}

	}
}

func (t *Trinity) migrate() {
	t.initEnumtype()
	// initial releationship table
	t.initUserGroup()       //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.initUserPermission()  //因为many2many自动生成的表会继承唯一约束，所以手动建立表
	t.initGroupPermission() //因为many2many自动生成的表会继承唯一约束，所以手动建立表

	MigrateModel(
		ToMigrateDB,
		ToCreatePermission,
		&Migration{},
		&Permission{},
		&AppError{},
		&Group{},
		&User{},
	)
	initPermissionDefaultValue()
	initGroupDefaultValue()
	initUserDefaultValue()
}
