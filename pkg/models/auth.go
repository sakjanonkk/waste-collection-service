package models

// LoginRequest DTO for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse DTO for login response
type LoginResponse struct {
	Token string `json:"token"`
	Staff Staff  `json:"staff"`
}

// AuthUser context key for authenticated user
type AuthUser struct {
	StaffID uint        `json:"staff_id"`
	Email   string      `json:"email"`
	Role    StaffRole   `json:"role"`
	Status  StaffStatus `json:"status"`
}
