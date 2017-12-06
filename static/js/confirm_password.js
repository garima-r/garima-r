function check_pass() {
	console.log("test");
    if (document.getElementById('password').value == document.getElementById('confirm_password').value) {
        	document.getElementById('signupbutton1').disabled = false;
    		document.getElementById("error_message").innerText = "";
    } else {
        	document.getElementById('signupbutton1').disabled = true;
        	document.getElementById("error_message").innerText = "Your Passwords and Confirm Passwords are not matching.";
    }
}
