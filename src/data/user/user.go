package user

import (
	"fmt"

	"github.com/uoregon-libraries/student-course-integrator/src/data/course"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
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

// Find looks for the user with the given login, returning an unauthorized user
// if the login isn't found or has no courses
func Find(login string) (*User, error) {
	if login == "" {
		return &User{}, nil
	}

	var u = &User{Login: login}
	return u, u.readDB()
}

// String serializes the user's information for display
func (u *User) String() string {
	var login = u.Login
	if u.Login == "" {
		login = "N/A"
	}
	return fmt.Sprintf("%s - %s - %v", login, u.IP, u.Authorized)
}

// readDB grabs all courses associated with this person (if any) and attaches
// course instances to the user.  If no courses are found, the user is
// considered unauthorized.  On any database error, the error is returned and
// the user is unchanged.
func (u *User) readDB() error {
	var rows, err = global.DB.Query(`
		SELECT c.canvas_id, c.description
		FROM faculty_courses fc
		JOIN courses c ON (fc.course_id = c.id)
		WHERE fc.login = ?`, u.Login)
	if err != nil {
		return err
	}

	var courses []*course.Course
	for rows.Next() {
		var crn, desc string
		var err = rows.Scan(&crn, &desc)
		if err != nil {
			return err
		}
		courses = append(courses, &course.Course{CRN: crn, Description: desc})
	}

	u.Courses = courses
	u.Authorized = len(u.Courses) > 0

	return nil
}
