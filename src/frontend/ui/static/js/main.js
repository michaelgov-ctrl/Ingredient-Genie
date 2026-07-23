// footer
var dayOfWeek = new Date().toLocaleString('en-us', {  weekday: 'long' })
if (dayOfWeek == "Thursday") {
	footerText = "I never could get the hang of Thursdays"
} else {
	footerText = "Have a nice " + dayOfWeek
}

// document.getElementsByTagName("footer")[0].getElementsByTagName("div")[0].innerHTML = footerText;
document.getElementsByTagName("footer")[0].innerHTML = footerText; //change content

// navigation bar
var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

// meal search
const ingredientInputs = document.getElementById("ingredient-inputs");
const addIngredientButton = document.getElementById("add-ingredient");

addIngredientButton.addEventListener("click", () => {
	const index = ingredientInputs.children.length;

	const container = document.createElement("div");
	container.className = "ingredient-input";

	const label = document.createElement("label");
	label.htmlFor = `ingredient-${index}`;

	const input = document.createElement("input");
	input.type = "text";
	input.id = `ingredient-${index}`;
	input.name = "ingredients";
	input.placeholder = "e.g. Garlic";

	container.appendChild(label);
	container.appendChild(input);

	ingredientInputs.appendChild(container);
});