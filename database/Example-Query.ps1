#Requires -Modules PSSQLite

Import-Module PSSQLite

$ErrorActionPreference = 'Stop'

$DataSource = "./meals.sqlite"

$Meals = Invoke-SqliteQuery -DataSource $DataSource -Query "SELECT * FROM Meal"
$Ingredients = Invoke-SqliteQuery -DataSource $DataSource -Query "SELECT * FROM Ingredient"
$MealIngredients = Invoke-SqliteQuery -DataSource $DataSource -Query "SELECT * FROM MealIngredient"
