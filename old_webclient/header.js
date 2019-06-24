if (isLoggedIn()) {
	$('#logout-link').show();
} else {
	$('#login-link').show();
}
$('#logout-link').click(function() {
	logout()
		.then(function() {
			window.location.href = '/';
		})
		.catch(function() {
			alert("Oops! There was an error logging out; you are still logged in. Please try again later.");
		});
});
