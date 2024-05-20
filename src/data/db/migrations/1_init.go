package migrations

import (
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/data/db"
	"base_structure/src/data/models"
	"base_structure/src/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

func Up1() {
	database := db.GetDb()
	createTables(database)
	createDefaultInformation(database)
}

func Down1() {

}

func createTables(database *gorm.DB) {
	var tables []interface{}

	// User
	tables = addNewTable(database, models.User{}, tables)
	tables = addNewTable(database, models.Role{}, tables)
	tables = addNewTable(database, models.RoleUser{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Migration, "cannot create tables through migration", nil)
		return
	}
	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createDefaultInformation(database *gorm.DB) {
	adminRole := models.Role{Name: constants.AdminRoleName}
	createRoleIfNotExists(database, &adminRole)
	defaultRole := models.Role{Name: constants.DefaultRoleName}
	createRoleIfNotExists(database, &defaultRole)
	u := models.User{
		Username:     constants.DefaultUserName,
		FirstName:    constants.AdminFirstName,
		LastName:     constants.AdminLastName,
		MobileNumber: constants.AdminMobileNumber,
		Email:        constants.AdminEmail,
	}
	pass := constants.AdminPassword
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	createAdminUserIfNotExists(database, &u, adminRole.Id)
}

func createRoleIfNotExists(database *gorm.DB, r *models.Role) {
	exists := 0
	database.Model(&models.Role{}).Select("1").Where("name = ?", r.Name).First(&exists)
	if exists == 0 {
		database.Create(r)
	}
}

func createAdminUserIfNotExists(database *gorm.DB, u *models.User, roleId int) {
	exists := 0
	database.Model(&models.User{}).Select("1").Where("username = ?", u.Username).First(&exists)
	if exists == 0 {
		database.Create(u)
		ru := models.RoleUser{UserId: u.Id, RoleId: roleId}
		database.Create(&ru)
	}
}
