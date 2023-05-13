package query

const (
	USER_LIST = `
		SELECT id, user_type, email, password, mobile, last_login, is_active, 
			created_by, created_at, updated_by, updated_at
		FROM users
	`
	USER_CREATE = `
		INSERT INTO users (user_type, email, password, mobile, is_active, created_by, created_at, updated_by, updated_at)
		VALUES (:user_type, :email, :password, :mobile, :is_active, :created_by, :created_at, :updated_by, :updated_at)
	`
	USER_UPDATE = `
		UPDATE users SET user_type = :user_type, email = :email, password = :password, last_login = :last_login, is_active = :is_active, 
			created_by = :created_by, created_at = :created_at, updated_by = :updated_by, updated_at = :updated_at, 
			mobile = :mobile
		WHERE id = :id
	`
	USER_TOKEN = `
		SELECT u.id, u.user_type, u.email, u.password, u.mobile, u.last_login, u.is_active, 
			u.created_by, u.created_at, u.updated_by, u.updated_at
		FROM users u, user_tokens ut
	`
	USER_UPDATE_PASSWORD = `
		UPDATE users SET password = :password WHERE id = :id
	`
)
