package models

// CreateUserRequest represents the payload to create a new user.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest represents the payload to update an existing user.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UserResponse is returned to API consumers.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age,omitempty"`
}



