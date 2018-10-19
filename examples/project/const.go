package main

// Redis database keys.
const (
	idSetsAccount = "idsets:Account"

	hashAccount = "hash:Account"

	basicInfoAccountPattern = "basic:info:Account:%s"
)

// Redis Basic Info Account properties
const (
	username = "username"

	password = "password"

	key = "key"
)


// Error pattern constants
const (
	missingPropertyBasicInfoAccountPatternError = "Key %s has no valid basic info account"

	redisConnectionFailedPatternError = "Unable to open redis connection %s"

	failedToLoadKeysPatternError = "Failed to load keys because %s"

	failedToLoadAccountsPatternError = "Failed to load accounts because %s"

	failedToLoginAccountsPatternError = "Failed to login %s accounts"
)