package trinity

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//Migration model Migration
type Migration struct {
	Simpmodel
	Seq    int    `json:"seq"  `
	Name   string `json:"name" gorm:"type:varchar(100);"`
	Status bool   `json:"status" gorm:"default:FALSE"`
	Error  string `json:"error" `
}

func runMigrationFile(seq int, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var migError string
	for scanner.Scan() {
		sql := scanner.Text()
		if err := GlobalTrinity.db.Exec(sql).Error; err != nil {
			migError += err.Error()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	mg := Migration{
		Seq:    seq,
		Name:   filepath,
		Status: true,
		Error:  migError,
	}
	GlobalTrinity.db.Create(&mg)
	return nil
}

// RunMigration func to run the migration
// scan the migration file under static/migrations
func RunMigration() {
	var migrationError error
	migrationsDirPath := filepath.Join(rootPath, GlobalTrinity.setting.GetMigrationPath())
	fileInfoList, err := ioutil.ReadDir(migrationsDirPath)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(fileInfoList, func(i, j int) bool {
		indexA, _ := strconv.Atoi(strings.Split(fileInfoList[i].Name(), "_")[0])
		indexB, _ := strconv.Atoi(strings.Split(fileInfoList[j].Name(), "_")[0])
		return indexA < indexB
	})
	var currentMigSeq int
	row := GlobalTrinity.db.Table(GlobalTrinity.setting.GetTablePrefix()+"migration").Where("status = ?", true).Select("MAX(seq)").Row()
	row.Scan(&currentMigSeq)
	for i := range fileInfoList {
		//0_filexxxx.sql
		if len(strings.Split(fileInfoList[i].Name(), ".sql")) < 2 {
			// not a .sql file , break
			fmt.Println(fileInfoList[i].Name() + " is not  .sql file , skip, please use 1_xxx.sql ,2_xxx.sql ")
		}
		seq, err := strconv.Atoi(strings.Split(fileInfoList[i].Name(), "_")[0])
		if err != nil {
			fmt.Println(fileInfoList[i].Name() + " don't have seq number,  skip , please use 1_xxx.sql ,2_xxx.sql")
		}
		if seq <= currentMigSeq {
			fmt.Println(fileInfoList[i].Name() + " already executed , skip ")
			// already executed
			continue
		}
		fmt.Println(fileInfoList[i].Name() + " start migration !")
		migrationError = runMigrationFile(seq, filepath.Join(migrationsDirPath, fileInfoList[i].Name()))
		if migrationError != nil {
			fmt.Println(fileInfoList[i].Name() + " excuting error , " + migrationError.Error())
		}
		fmt.Println(fileInfoList[i].Name() + " end migration !")
	}

	fmt.Println("run all migrations successfully")
}
