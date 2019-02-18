document.addEventListener("DOMContentLoaded", function() {
  // Don't validate when clicking the "go back" button
  var goBack = document.getElementById("submit-go-back");
  goBack.addEventListener("click", function() {
    document.getElementById("association-form-confirm").setAttribute("novalidate", "novalidate");
  });

  var agree = document.getElementById("graderReqMet");
  var msg = "You must agree to the terms for assigning a Grader to a course";
  agree.setCustomValidity(msg);
  agree.addEventListener("change", function() {
    this.setCustomValidity(this.validity.valueMissing ? msg : "");
  }, false);

}, false);
