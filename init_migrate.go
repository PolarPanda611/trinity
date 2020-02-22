package trinity

import "fmt"

// ToMigrateDB migrate db function
func (t *Trinity) migrateModel(m interface{}) {
	if err := t.db.AutoMigrate(m).Error; err != nil {
		t.logger.Print(fmt.Sprintf("[info] Migrate Model %v Error : %v  ", GetTypeName(m, true), err))
	}
}

// ToCreatePermission create permission func
func (t *Trinity) migratePermission(pSlice []string) {
	p := Permission{Code: pSlice[0], Name: pSlice[1]}
	p.CreateOrInitPermission(t)
}

// Migrate to migrate model
func (t *Trinity) Migrate(modelToMigrate ...interface{}) {
	for _, v := range modelToMigrate {
		modelName := GetTypeName(v, true)
		t.migrateModel(v)
		PermissionList := [][]string{
			{"system.add." + modelName, "app.right.system.add." + modelName},
			{"system.view." + modelName, "app.right.system.view." + modelName},
			{"system.edit." + modelName, "app.right.system.edit." + modelName},
			{"system.delete." + modelName, "app.right.system.delete." + modelName},
		}
		for _, v := range PermissionList {
			t.migratePermission(v)
		}
	}
}
