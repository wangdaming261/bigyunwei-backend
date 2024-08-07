package common

const (
	GIN_CTX_CONFIG_LOGGER = "gin_logger"
	GIN_CTX_CONFIG_CONFIG = "gin_config"
	GIN_CTX_JWT_CLAIM     = "jwt_claim"
	GIN_CTX_JWT_USER_NAME = "jwt_user_name"
)

var (
	COMMON_SHOW_MAP = map[string]bool{
		"1": true,
		"0": true,
	}
)
