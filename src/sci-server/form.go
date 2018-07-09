package sciserver

import (
	"errors"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/data/course"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
	"github.com/uoregon-libraries/student-course-integrator/src/person"
)

// form holds the submission data: crn and duckid
type form struct {
	CRN     string
	DuckID  string
	User    *user.User
	Course  *course.Course
	Student *person.Person
	Confirm string
	errors  []error
}

func (r *response) getForm() (f form, err error) {
	f.CRN = r.req.PostFormValue("crn")
	f.DuckID = r.req.PostFormValue("duckid")
	f.Confirm = r.req.PostFormValue("confirm")
	f.User = r.user

	if f.DuckID == "" {
		f.errors = append(f.errors, errors.New("duckid must be filled out"))
	}
	if f.CRN == "" {
		f.errors = append(f.errors, errors.New("a course must be chosen"))
	}

	// if we have a missing field, we don't bother with further validation
	if len(f.errors) > 0 {
		return f, err
	}

	// Make sure the duckid is valid and represents a GTF
	f.Student, err = person.FindByDuckID(f.DuckID)
	if err != nil {
		return f, err
	}
	if f.Student == nil {
		f.errors = append(f.errors, errors.New("duckid doesn't match a known student"))
	} else if !f.Student.IsGTF() {
		f.errors = append(f.errors, errors.New(f.Student.DisplayName+" is not a GTF"))
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
