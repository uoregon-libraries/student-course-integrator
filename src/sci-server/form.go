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
	Role            string // Role holds the submitted dropdown value for the selected role, e.g., "GE"
	GraderConfirmed string // GraderConfirmed is set only if Grader and faculty clicks graderReqMet

	// Derived fields
	User           *user.User     // User gets populated with the faculty member who is logged in
	Course         *course.Course // Course gets the course by looking up the faculty member + CRN
	Agent          *person.Person // Agent is the person being assigned a role, after looking up the DuckID
	GraderConfReqd string         // GraderConfReqd is "1" when role == Grader; for displaying the confirm checkbox

	// Internal data
	errors []error // errors will be populated with anything preventing the form submission, e.g., a bad duckid
}

func (r *responder) getForm() (f *form, err error) {
	f = r.vars.Form
	f.CRN = r.req.PostFormValue("crn")
	f.DuckID = r.req.PostFormValue("duckid")
	f.Confirm = r.req.PostFormValue("confirm")
	f.User = r.vars.User
	f.Role = r.req.PostFormValue("role")
	f.GraderConfirmed = r.req.PostFormValue("graderReqMet")

	if f.DuckID == "" {
		f.errors = append(f.errors, errors.New("duckid must be filled out"))
	}
	if f.CRN == "" {
		f.errors = append(f.errors, errors.New("a course must be chosen"))
	}
	if f.Role == "" || !roles.IsValid(f.Role) {
		f.errors = append(f.errors, errors.New("a valid role must be chosen"))
	}
	// if we have a missing field, we don't bother with further validation
	if len(f.errors) > 0 {
		return f, err
	}

	// Find will handle either a duckid or a bannerid and return a ldap-person record if valid.
	// Make sure the record represents someone in the system
	f.Agent, err = person.Find(f.DuckID)
	if err != nil {
		return f, err
	}
	if f.Agent == nil {
		f.errors = append(f.errors, errors.New("nobody with this ID exists"))
	} else {
		switch f.Role {
		case roles.GE:
			if !f.Agent.CanBeGE() {
				f.errors = append(f.errors, errors.New(f.Agent.DisplayName+" is currently not classified as a GE"))
			}
		case roles.Grader:
			f.GraderConfReqd = "1"
			if f.Confirm == "1" && f.GraderConfirmed != "1" {
				f.errors = append(f.errors, errors.New(f.Agent.DisplayName+" must meet the Grader requirement"))
			}
		}
	}

	// Make sure the logged-in user is allowed to assign people to this crn
	f.Course = f.User.FindCourse(f.CRN)
	if f.Course == nil {
		f.errors = append(f.errors, errors.New("the chosen course is invalid"))
	}

	return f, err
}

func (f form) errorString() string {
	var strs = make([]string, len(f.errors))
	for i, err := range f.errors {
		strs[i] = err.Error()
	}
	return strings.Join(strs, ", ")
}
