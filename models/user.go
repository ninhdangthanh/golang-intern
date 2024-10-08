package models

type UserModel struct {
	ID       uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string         `gorm:"type:varchar(100);not null" json:"password"`
	Products []ProductModel `gorm:"foreignKey:UserID" json:"-"`
}

func (UserModel) TableName() string {
	return "users"
}
