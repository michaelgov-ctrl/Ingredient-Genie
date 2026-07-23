var dayOfWeek = new Date().toLocaleString('en-us', {  weekday: 'long' })
if (dayOfWeek == "Thursday") {
	footerText = "I never could get the hang of Thursdays"
} else {
	footerText = "Have a nice " + dayOfWeek
}

// document.getElementsByTagName("footer")[0].getElementsByTagName("div")[0].innerHTML = footerText;
document.getElementsByTagName("footer")[0].innerHTML = footerText; //change content

var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}