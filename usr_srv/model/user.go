package model

import "time"

type BaseModel struct {
	ID        int32      `gorm:"primarykey"`
	CreatedAt *time.Time `gorm:"column:add_time"`
	UpdatedAt *time.Time `gorm:"column:update_time"`
	DeletedAt *time.Time `gorm:"column:delete_time"`
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"column:mobile;index;unique;type:varchar(11);not null"`
	Password string     `gorm:"column:password;type:varchar(100);not null"`
	NickName string     `gorm:"column:nick_name;type:varchar(20)"`
	BirthDay *time.Time `gorm:"column:birth_day;type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}
