package sciserver

import (
	"errors"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/data/course"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
	"github.com/uoregon-libraries/student-course-integrator/src/person"
	"github.com/uoregon-libraries/student-course-integrator/src/roles"
)

// form holds the submission data as well as derived fields which are loaded
// based on the user filling out the form and the form fields themselves
type form struct {
	// Actual form fields: CRN and DuckID
	CRN             string // CRN holds the submitted dropdown value for the selected course, e.g., "201704.X3159"
	DuckID          string // DuckID holds the submitted duckid of the user being added, e.g., "echjeremy"
	Confirm         string // Confirm is "1" when the form is re-submitted after confirmation of the Agent by name
	Role            string // Role holds the submitted dropdown value for the selected role, e.g., "TA"
	GraderConfirmed string // GraderConfirmed is set only if Grader and faculty clicks graderReqMet

	// Derived fields
	User   *user.User     // User gets populated with the faculty member who is logged in
	Course *course.Course // Course gets the course by looking up the faculty member + CRN
	Agent  *person.Person // Agent is the person being assigned a role, after looking up the DuckID

	// Internal data
	errors []error // errors will be populated with anything preventing the form submission, e.g., a bad duckid
}

func (r *responder) getForm() (f *form, err error) {
	f = r.vars.Form
	f.errors = []error{errors.New("this process is no longer available")}
	return f, err
}

func (f form) errorString() string {
	var strs = make([]string, len(f.errors))
	for i, err := range f.errors {
		strs[i] = err.Error()
	}
	return strings.Join(strs, ", ")
}

func (f form) IsGrader() bool {
	if f.Role == roles.Grader {
		return true
	}
	return false
}
