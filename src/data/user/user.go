package user

import (
	"fmt"

	"github.com/uoregon-libraries/student-course-integrator/src/data/course"
)

// User holds data for people who have been authenticated
type User struct {
	Login   string
	Courses []*course.Course

	// Authorized users are faculty we have seen in our feed, and must have CRNs associated with them
	Authorized bool
	// IP is obviously not persistent; it just gives us extra information in case we need it
	IP string
}

// Find looks for the user with the given login, returning an empty structure
// if the given user is not found
func Find(login string) *User {
	return &User{}
}

// String serializes the user's information for display
func (u *User) String() string {
	var login = u.Login
	if u.Login == "" {
		login = "N/A"
	}
	return fmt.Sprintf("%s - %s - %v", login, u.IP, u.Authorized)
}
