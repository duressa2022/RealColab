package domain

type LoginRequest struct {
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type ChangePassword struct {
	Password    string `json:"password" bson:"password"`
	NewPassword string `json:"newPassword" bson:"newpassword"`
}


