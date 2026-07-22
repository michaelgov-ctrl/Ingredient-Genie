$body = @{
	"ingredients" = @(
		"Garlic",
		"Red Onions",
		"Vegetable Oil",
		"Lime"
	)
	"filters" = @{
		"page" = 1
		"pageSize" = 10
		"sort" = "-ratio"
	}
} | ConvertTo-Json
$resp = irm http://localhost:4269/v1/meals/search -Method POST -Body $body -ContentType "application/json"

$resp.meals[0]