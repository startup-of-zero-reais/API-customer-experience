package domain

const (
	UserCreated = Event("user_created")
	UserUpdated = Event("user_updated")
	UserDeleted = Event("user_deleted")

	SessionStarted = Event("session_started")
	SessionEnded   = Event("session_ended")

	RequestPasswordRecover = Event("request_password_recover")
	PasswordRecovered      = Event("password_recovered")

	FavoriteAdded   = Event("favorite_added")
	FavoriteRemoved = Event("favorite_removed")
	FavoriteQuery   = Event("favorite_query")
)
