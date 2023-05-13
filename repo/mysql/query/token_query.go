package query

const (
	USER_TOKENS_INSERT = `
		INSERT INTO user_tokens (token, jwt_id, created_at, expired_at, is_used, invalidated, user_id)
		VALUES (:token, :jwt_id, :created_at, :expired_at, :is_used, :invalidated, :user_id)
	`
	
	USER_TOKENS_GET = `
		SELECT token, jwt_id, created_at, expired_at, is_used, invalidated, user_id FROM user_tokens
	`

	USER_TOKENS_UPDATE = `
		UPDATE user_tokens 
		SET jwt_id = :jwt_id, created_at = :created_at, expired_at = :expired_at, is_used = :is_used, invalidated = :invalidated, user_id = :user_id
		WHERE token = :token
	`
)