package security

/**
 * Enumération des droits
 */
const (
	NONE       = 0
	READ       = 1 << 0
	WRITE_OWN  = 1 << 1
	WRITE_ALL  = 1 << 2
	SUPER_USER = 1 << 3
	USER       = READ | WRITE_OWN
	MOD        = USER | WRITE_ALL
	ADMIN      = MOD | SUPER_USER
)
