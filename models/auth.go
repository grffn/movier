package models

//RegisterModel - Registration View Model
type RegisterModel struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//LoginModel - Login View Model
type LoginModel struct {
	UserID   string `json:"userid" binding:"required"`
	Password string `json:"password" binding:"required"`
}
