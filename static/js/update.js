// Validating Empty Field
function check_empty() {
if (document.getElementById('Fname').value == "" || document.getElementById('Lname').value == "" || document.getElementById('Dob').value == "") {
alert("Fill All Fields !");
} else {
document.getElementById('updateform').submit();
alert("Form Submitted Successfully..."); // try to redirect to profile page
}
}
//Function To Display Popup
function div_show() {
	document.getElementById('updateprofile').style.display = "block";
}
//Function to Hide Popup
function div_hide(){
	document.getElementById('updateprofile').style.display = "none";
}
