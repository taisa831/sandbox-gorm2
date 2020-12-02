package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	db2 "github.com/taisa831/sandbox-gorm2/db"
	"github.com/taisa831/sandbox-gorm2/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

var db *gorm.DB

func main() {
	var err error
	db, err = db2.ConnWithLogger()
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()


	CreateCompany()

	// Major Features
	//

	ctx := context.Background()
	FindWithContext(ctx)

	CreateUserBatch()

	PrepareStmt()

	DryRun()

	JoinPreload()

	// 動かない
	//FindToMap()

	CreateFromMap()

	FindInBatches()

	NestedTransaction()

	SavePointRollbackTo()

	NamedArgument()

	GroupConditions()

	SubQuery()

	Upsert()

	Locking()

	OptimizerIndexCommentHints()

	// Field permissions
	//type User struct {
	//	Name string `gorm:"<-:create"` // allow read and create
	//	Name string `gorm:"<-:update"` // allow read and update
	//	Name string `gorm:"<-"`        // allow read and write (create and update)
	//	Name string `gorm:"->:false;<-:create"` // createonly
	//	Name string `gorm:"->"` // readonly
	//	Name string `gorm:"-"`  // ignored
	//}

	DataTypes()

	SmartSelect()

	AssociationsBatchMode()

	// 想定通りに動かない
	DeleteAssociationsWhenDeleting()
}

func CreateCompany() {
	user := model.Company{
		Name: "company name",
	}
	db.Create(&user)
}

func FindWithContext(ctx context.Context) {
	var user model.User
	db.WithContext(ctx).Find(&user)
}

func CreateUserBatch() {
	var users = []model.User{
		{Name: "test", Address: "address", Age: 20, CompanyID: 1},
		{Name: "test2", Address: "address2", Age: 30, CompanyID: 1},
		{Name: "test3", Address: "address3", Age: 40, CompanyID: 1},
	}
	db.Create(&users)
}

func PrepareStmt() {
	var user model.User
	tx := db.Session(&gorm.Session{
		PrepareStmt: true,
	})
	tx.First(&user)
}

func DryRun() {
	var user model.User
	stmt := db.Session(&gorm.Session{DryRun: true}).Find(&user, 1).Statement
	fmt.Println(stmt.SQL.String())
}

// https://github.com/go-gorm/gorm/issues/1436
func JoinPreload() {
	var users []model.User
	db.Joins("Company").Find(&users, "users.id IN ?", []int{1, 2})
	b, _ := json.Marshal(&users)
	fmt.Println(string(b))
}

//func FindToMap() {
//	var result map[string]interface{}
//	db.Model(&model.User{}).First(&result, "id = ?", 1)
//}

func CreateFromMap() {
	db.Model(&model.User{}).Create(map[string]interface{}{"Name": "test", "Address": "address", "Age": 50, "CompanyID": 1})

	datas := []map[string]interface{}{
		{"Name": "test", "Address": "address1", "Age": 60, "CompanyID": 1},
		{"name": "test", "Address": "address2", "Age": 70, "CompanyID": 1},
	}
	db.Model(&model.User{}).Create(datas)
}

func FindInBatches() {
	var users []model.User
	err := db.Where("id >= ?", 1).FindInBatches(&users, 2, func(tx *gorm.DB, batch int) error {
		for _, user := range users {
			fmt.Println(user.ID)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error)
	}
}

func NestedTransaction() {
	db.Transaction(func(tx *gorm.DB) error {
		user1 := model.User{
			Name:      "name",
			Address:   "address",
			Age:       10,
			CompanyID: 1,
		}
		tx.Create(&user1)

		tx.Transaction(func(tx2 *gorm.DB) error {
			user2 := model.User{
				Name:      "name",
				Address:   "address",
				Age:       10,
				CompanyID: 1,
			}
			tx.Create(&user2)
			return errors.New("rollback user2") // rollback user2
		})

		tx.Transaction(func(tx2 *gorm.DB) error {
			user3 := model.User{
				Name:      "name",
				Address:   "address",
				Age:       10,
				CompanyID: 1,
			}
			tx.Create(&user3)
			return nil
		})
		return nil // commit user1 and user3
	})
}

func SavePointRollbackTo() {
	tx := db.Begin()

	user1 := model.User{
		Name:      "name1",
		Address:   "address1",
		Age:       10,
		CompanyID: 1,
	}
	tx.Create(&user1)

	tx.SavePoint("sp1")

	user2 := model.User{
		Name:      "name2",
		Address:   "address2",
		Age:       10,
		CompanyID: 1,
	}
	tx.Create(&user2)

	tx.RollbackTo("sp1") // rollback user2

	tx.Commit() // commit user1
}

func NamedArgument() {
	var user model.User
	db.Where("name1 = @name OR name2 = @name", sql.Named("name", "test")).Find(&user)

	var user2 model.User
	db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "test"}).First(&user2)

	db.Raw(
		"SELECT * FROM users WHERE name1 = @name OR name2 = @name2 OR name3 = @name",
		sql.Named("name", "test"), sql.Named("name2", "test"),
	).Find(&model.User{})

	db.Exec(
		"UPDATE users SET name1 = @name, name2 = @name2, name3 = @name",
		map[string]interface{}{"name": "test", "name2": "test"},
	)
}

func GroupConditions() {
	var user model.User
	db.Where(
		db.Where("name = ?", "test").Where(db.Where("address = ?", "address").Or("company = ?", 1)),
	).Or(
		db.Where("name = ?", "test2").Where("address = ?", "address"),
	).Find(&user)
}

func SubQuery() {
	// Where SubQuery
	var user model.User
	db.Where("age > (?)", db.Table("users").Select("AVG(age)")).Find(&user)

	// From SubQuery
	var user2 model.User
	db.Table("(?) as u", db.Model(&model.User{}).Select("name", "age")).Where("age = ?", 20).Find(&user2)

	// Update SubQuery
	db.Model(&model.User{}).Update(
		"name", db.Model(&model.Company{}).Select("name").Where("companies.id = users.company_id"),
	)
}

func Upsert() {
	var users []model.User
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"name": "test", "address": "address", "age": 20, "company_id": 1}),
	}).Create(&users)

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "address", "age", "company_id"}),
	}).Create(&users)
}

func Locking() {
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&model.User{})

	db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Find(&model.User{})
}

func OptimizerIndexCommentHints() {
	// Optimizer Hints
	db.Clauses(hints.New("hint")).Find(&model.User{})

	// Index Hints
	db.Clauses(hints.UseIndex("idx_user_name")).Find(&model.User{})

	// Comment Hints
	db.Clauses(hints.Comment("select", "master")).Find(&model.User{})
}

func DataTypes() {
	db.Create(&model.User{
		Name:       "user",
		CompanyID:  1,
		Address:    "address",
		Age:        30,
		Attributes: datatypes.JSON([]byte(`{"name": "user", "age": 20, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
	})

	var user model.User
	db.First(&user, datatypes.JSONQuery("attributes").HasKey("name"))
	db.First(&user, datatypes.JSONQuery("attributes").HasKey("orgs", "orga"))
}

type APIUser struct {
	ID   int
	Name string
}

func SmartSelect() {
	db.Model(&model.User{}).Limit(10).Find(&APIUser{})
}

func AssociationsBatchMode() {
	userA := model.User{ID: 1}
	userB := model.User{ID: 3}
	users := []model.User{userA, userB}
	var creditCard model.CreditCard

	db.Model(&users).Association("CreditCard").Find(&creditCard)

	db.Model(&users).Association("CreditCard").Delete(&userA)

	db.Model(&users).Association("CreditCard").Count()
}

func DeleteAssociationsWhenDeleting() {
	// delete user's account when deleting user
	user := model.User{
		ID: 2,
	}
	db.Select("Company").Delete(&user)
}
