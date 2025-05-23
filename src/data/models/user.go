package models

type User struct {
	BaseModel
	Username     string `gorm:"type:string;size:20;not null;unique"`
	FirstName    string `gorm:"type:string;size:15;null"`
	LastName     string `gorm:"type:string;size:25;null"`
	MobileNumber string `gorm:"type:string;size:11;null;unique;default:null"`
	Email        string `gorm:"type:string;size:64;null;unique;default:null"`
	Password     string `gorm:"type:string;size:64;not null"`
	Enabled      bool   `gorm:"default:true"`
	RoleUsers    *[]RoleUser
}

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:30;not null;unique"`
	RoleUsers *[]RoleUser
}

type RoleUser struct {
	BaseModel
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	RoleId uint
	UserId uint
}
