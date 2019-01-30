package roles

var Roles = [2]string{"GE", "Grader"}

func validateRole(role string) bool {

  for _, element := range Roles {
    if role == element{
      return true
    }
  }
  return false
}