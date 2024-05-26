package appError

const (
	ErrUnverifiedEmail             = "Email is unverified"
	ErrEmptyEmail                  = "Email is empty"
	ErrNotTimeForRefreshToken      = "token is not valid yet"
	ErrUserNotFound                = "User not found"
	ErrForbidden                   = "forbidden"
	ErrUnAuthentication            = "unAuthentication"
	ErrNotMyPreset                 = "not my preset"
	ErrCannotUpdatePublishedPreset = "cannot update published preset"
	ErrCannotTagUnpublished        = "cannot tag unpublished preset"
	ErrUserInactive                = "user is inactive"
	ErrInvalidPresetInput          = "invalid input"
	ErrStoreNotFound               = "store not found"
	ErrBadInput                    = "bad Request"
)
