package entity

// user entity
type User struct {
	CommonFields
	Username string `gorm:"uniqueIndex;not null" json:"username" validate:"required"`
	Email    string `gorm:"not null" json:"email" validate:"required,email"`
	// password hash
	Password string `gorm:"not null" json:"-" validate:"required"`
	Bio      string `gorm:"size:1024" json:"bio" `
	Image    string `json:"image"`
	Admin    bool   `gorm:"default:false" json:"admin"`
	// "many to many" association
	// https://gorm.io/docs/many_to_many.html
	Roles []Role `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"roles"`
}

// used only for password validation, since you cannot validate a password hash
type Password struct {
	Value string `validate:"required,min=8,max=64"`
}

// user role: a collection of permissions
type Role struct {
	CommonFieldsName
	Description string       `gorm:"size:1024" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"permissions"`
}

// user permission: a path and actions allowed on that path
type Permission struct {
	CommonFieldsName
	Description  string       `gorm:"size:1024" json:"description"`
	Path         string       `gorm:"not null" json:"path" validate:"required"`
	Capabilities []Capability `gorm:"not null;many2many:permission_capabilities;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"capabilities"`
}

// cabability: allowed action
type Capability struct {
	CommonFieldsName
	Description string `gorm:"size:1024" json:"description"`
}
