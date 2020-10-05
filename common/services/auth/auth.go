package auth

// IsAdmin identifies admin status
func IsAdmin(discordID string) bool {
	admins := []string{"192906671167635457", "691353925093294161"}

	for _, id := range admins {
		if id == discordID {
			return true
		}
	}

	return false
}
