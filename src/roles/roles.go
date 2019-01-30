package roles

// Roles lists the roles
var Roles = []string{"GE", "Grader"}

func validateRole(role string) bool {

	for _, element := range Roles {
		if role == element {
			return true
		}
	}
	return false
}
