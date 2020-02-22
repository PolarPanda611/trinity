package trinity

import "github.com/jinzhu/gorm"

//QueryByPreload handling preload
func QueryByPreload(PreloadList map[string]func(db *gorm.DB) *gorm.DB) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		if len(PreloadList) > 0 {
			for k, v := range PreloadList {
				if v == nil {
					db = db.Preload(k)
				} else {
					db = db.Preload(k, v)
				}

			}
		}
		return db
	}
}
