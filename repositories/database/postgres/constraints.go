package postgres

const (
	// UNIQUE CONSTRAINTS
	UniqueUsersEmail    = "users_email_key"
	UniqueUsersUsername = "users_username_key"

	// FOREIGN KEYS
	FKeyPostsUserId = "posts_user_id_fkey"
)
