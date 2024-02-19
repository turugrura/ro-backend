package repository

type RefreshToken struct {
	Id        string `bson:"_id,omitempty"`
	UserId    string `bson:"user_id"`
	Count     uint32 `bson:"count"`
	UserAgent string `bson:"user_agent"`
	CreatedAt string `bson:"created_at"`
	UpdatedAt string `bson:"updated_at"`
}

type PatchRefreshToken struct {
	UserId    string `bson:"user_id,omitempty"`
	Count     uint32 `bson:"count,omitempty"`
	UserAgent string `bson:"user_agent,omitempty"`
	CreatedAt string `bson:"created_at,omitempty"`
	UpdatedAt string `bson:"updated_at,omitempty"`
}

type CreateRefreshTokenInput struct {
	UserId    string
	UserAgent string
	Count     uint32
}

type UpdateRefreshTokenInput struct {
	Id    string
	Count uint32
}

type RefreshTokenRepository interface {
	GetRefreshTokenById(string) (*RefreshToken, error)
	CreateRefreshToken(CreateRefreshTokenInput) (*RefreshToken, error)
	UpdateRefreshToken(UpdateRefreshTokenInput) (*RefreshToken, error)
	DeleteRefreshTokenByUserId(string) error
}
