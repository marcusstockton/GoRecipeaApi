package chef

import "gorm.io/gorm"

/// Will act as the user model - all of the users will be chefs, and they will have recipes associated with them
type Chef struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
}
