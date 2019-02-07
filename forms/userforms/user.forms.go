package forms

// LoginForm struct a form to login a Tyr User.
type LoginForm struct {
	Email    string `bson:"email" json:"email" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required"`
}

// RegisterForm struct a form for register a Tyr User.
type RegisterForm struct {
	Email                string `bson:"email" json:"email" binding:"required"`
	Password             string `bson:"password" json:"password" binding:"required"`
	PasswordConfirmation string `bson:"passwordConfirmation" json:"passwordConfirmation" binding:"required"`
	First                string `bson:"firstName" json:"firstName" binding:"required"`
	Last                 string `bson:"lastName" json:"lastName" binding:"required"`
}