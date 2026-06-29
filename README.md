# Ingredient-Genie
Ingredient Genie is a web application that helps users discover recipes based on the ingredients they already have at home. Users enter the ingredients currently available and the application generates recipes that use those ingredients. The goal is to reduce food waste, save money on groceries, and make meal planning faster and more convenient.


# Project Overview

Ingredient-Genie will be a browser-based web application consisting of a Python backend API and a frontend built with HTML, CSS, and JavaScript.

The application will have a JSON REST API backend that accepts an arbitrary length list of ingredients from the user, processes the request, and returns recipes that best match the available ingredients.

The project follows a client-server architecture where the frontend communicates with the backend through REST API endpoints. A front end will be developed however the REST API will be directly accessible.

# Functional Requirements

## User Input

The front end application shall allow users to:

* Enter a variable length list of ingredients.
    * we must agree here in the spec on things like if a null list returns all recipes, etc..
* Edit ingredient list.
* Submit ingredients to the backend.
* Consume backend JSON responses.

## Recipe Search

The backend shall:

* Accept ingredient lists through a REST API.
* Search a recipe database or recipe dataset.
* Return recipes containing those ingredients as a JSON response.
* Rank recipes by the number of matching ingredients.

## Recipe Display

The frontend must display:

* Recipe name
* Ingredient list
* Cooking instructions
* How many and which ingredients the recipe uses

## Error Handling

The frontend and backend must:

* Validate input.
* Handle internal errors and return valid HTTP status codes.

# Technical Requirements

## Frontend

Technologies:

* HTML
* CSS
* JavaScript

Responsibilities:

* User interface
* Ingredient entry
* API communication
* Display recipe results

## Backend

Technologies:

* Python
* Flask???
* REST API

Responsibilities:

* API endpoints
* Recipe search logic
* Input validation
* JSON responses

## Data Format

Frontend sends:

```json
{
  "ingredients": [
    "chicken",
    "rice",
    "broccoli"
  ]
}
```

Backend returns:

```json
[
  {
    "title": "Chicken Fried Rice",
    "ingredients": [
      "Chicken",
      "Rice",
      "Broccoli",
      "Soy Sauce"
    ],
    "instructions": "...",
    "missingIngredients": [
      "Soy Sauce"
    ]
  }
]
```

# Project Work Breakdown

### API Development

Responsibilities

* Design REST API
* Create endpoints
* Request validation
* JSON responses
* API testing
* Documentation
* Recipe dataset management
* Ingredient matching algorithm
* Search optimization
* Data processing
* Backend testing

## Frontend

### User Interface

Responsibilities

* Page layout
* Ingredient input interface
* Search button
* Results display
* Responsive styling
* API integration
* User experience improvements

# Development Workflow

Our primary means for version control will be git with this github repository.
Main is the primary working branch.

Each developer will clone this repository, work via feature branches, push changes to their feature branch to Github, for merge within Github. Developers will communicate once their changes are merged into Main so that other developers may pull changes into their local repo. 

### To create a feature branch for development

`git clone git@github.com:TylerDoc/Ingredient-Genie.git ; cd Ingredient-Genie` -- run only once
`git checkout -b "my-feature-branch"`
    ... make your changes ...
`git add .`
`git commit -m "mildly descriptive comment"`
`git push -u origin my-feature-branch`
    ... open a pull request in the github UI so that we are able to track significant changes ...

### Pull changes that were pulled into the main branch

    ... you are currently on your working branch
    make sure to git add/commit/push or stash the local changes for your branch ...
`git checkout main`
`git pull origin main`
`git checkout my-feature-branch`
`git rebase main` -- resolve any conflicts
