package roles

// Constant role values so we can avoid the hard-coded strings all over the app
const (
	TA     = "TA"
	Grader = "Grader"
)

// Roles contains the roles that can be assigned by users of the app.
var Roles = []string{TA, Grader}

// IsValid checks that the value being passed around is valid.
func IsValid(role string) bool {

	for _, element := range Roles {
		if role == element {
			return true
		}
	}
	return false
}
