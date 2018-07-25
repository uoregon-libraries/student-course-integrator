# Grab dependencies - anything where the first part of the path is *.???
go list -f '{{join .Deps "\n"}}' | \
  grep "^[^/]\+\.[^/]\{3\}" | \
  grep -v "uoregon-libraries/student-course-integrator" | xargs -i% go get %

go install ./src/...
