#Requires -Modules PSSQLite

Import-Module PSSQLite

$ErrorActionPreference = 'Stop'

function ConvertTo-TrimmedOrNull {
    param (
        [Parameter()]
        [AllowNull()]
        [object] $Value
    )

    if ($null -eq $Value) {
        return $null
    }

    $text = ([string]$Value).Trim()

    if ([string]::IsNullOrWhiteSpace($text)) {
        return $null
    }

    return $text
}

function New-PreparedSQLiteCommand {
    param (
        [Parameter(Mandatory)]
        $Connection,

        [Parameter(Mandatory)]
        $Transaction,

        [Parameter(Mandatory)]
        [string] $CommandText,

        [Parameter(Mandatory)]
        [hashtable] $ParameterTypes
    )

    $command = $Connection.CreateCommand()
    $command.Transaction = $Transaction
    $command.CommandText = $CommandText

    foreach ($name in $ParameterTypes.Keys) {
        $parameter = $command.CreateParameter()
        $parameter.ParameterName = "@$name"
        $parameter.DbType = $ParameterTypes[$name]

        [void]$command.Parameters.Add($parameter)
    }

    $command.Prepare()

    return $command
}

function Set-SQLiteCommandValues {
    param (
        [Parameter(Mandatory)]
        $Command,

        [Parameter(Mandatory)]
        [hashtable] $Values
    )

    foreach ($name in $Values.Keys) {
        $value = $Values[$name]

        $databaseValue = if ($null -eq $value) {
            [DBNull]::Value
        }
        else {
            $value
        }

        $Command.Parameters["@$name"].Value = $databaseValue
    }
}

$url = "https://www.themealdb.com/api/json/v1/1"

$meals = 'a'..'z' | ForEach-Object {
    (Invoke-RestMethod "$url/search.php?f=$_").meals
} | Where-Object { $null -ne $_ }

$normalizedMeals = foreach ($meal in $meals) {
    $ingredients = for ($i = 1; $i -le 20; $i++) {
        $ingredientProperty = "strIngredient$i"
        $measureProperty    = "strMeasure$i"

        $ingredientName = ConvertTo-TrimmedOrNull -Value $meal.$ingredientProperty

        if ($null -eq $ingredientName) {
            continue
        }

        [pscustomobject]@{
            Position       = $i
            Name           = $ingredientName
            NormalizedName = (
                $ingredientName -replace '\s+', ' '
            ).ToLowerInvariant()
            MeasureText    = ConvertTo-TrimmedOrNull -Value $meal.$measureProperty
        }
    }

    [pscustomobject]@{
        ExternalMealId = ConvertTo-TrimmedOrNull -Value $meal.idMeal
        Name           = ConvertTo-TrimmedOrNull -Value $meal.strMeal
        AlternateName  = ConvertTo-TrimmedOrNull -Value $meal.strMealAlternate
        Category       = ConvertTo-TrimmedOrNull -Value $meal.strCategory
        Area           = ConvertTo-TrimmedOrNull -Value $meal.strArea
        Country        = ConvertTo-TrimmedOrNull -Value $meal.strCountry
        Instructions   = ConvertTo-TrimmedOrNull -Value $meal.strInstructions
        ThumbnailUrl   = ConvertTo-TrimmedOrNull -Value $meal.strMealThumb
        YoutubeUrl     = ConvertTo-TrimmedOrNull -Value $meal.strYoutube
        SourceUrl      = ConvertTo-TrimmedOrNull -Value $meal.strSource
        Ingredients    = @($ingredients)
    }
}

$databasePath = Join-Path $PSScriptRoot "meals.sqlite"

$connection = New-SQLiteConnection -DataSource $databasePath
$transaction = $null
$commands = @()

try {
    $schema = @'
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Meal (
    MealId          INTEGER PRIMARY KEY,
    ExternalMealId  TEXT NOT NULL UNIQUE,
    Name            TEXT NOT NULL,
    AlternateName   TEXT,
    Category        TEXT,
    Area            TEXT,
    Country         TEXT,
    Instructions    TEXT,
    ThumbnailUrl    TEXT,
    YoutubeUrl      TEXT,
    SourceUrl       TEXT
);

CREATE TABLE IF NOT EXISTS Ingredient (
    IngredientId    INTEGER PRIMARY KEY,
    Name            TEXT NOT NULL,
    NormalizedName  TEXT NOT NULL COLLATE NOCASE UNIQUE
);

CREATE TABLE IF NOT EXISTS MealIngredient (
    MealId          INTEGER NOT NULL,
    IngredientId    INTEGER NOT NULL,
    Position        INTEGER NOT NULL CHECK (Position BETWEEN 1 AND 20),
    MeasureText     TEXT,

    PRIMARY KEY (MealId, Position),

    FOREIGN KEY (MealId)
        REFERENCES Meal (MealId)
        ON DELETE CASCADE,

    FOREIGN KEY (IngredientId)
        REFERENCES Ingredient (IngredientId)
);

CREATE INDEX IF NOT EXISTS IX_MealIngredient_IngredientId
    ON MealIngredient (IngredientId);
'@

    $schemaQueryParams = @{
        SQLiteConnection = $connection
        Query            = $schema
        ErrorAction      = "Stop"
    }

    Invoke-SqliteQuery @schemaQueryParams

    $transaction = $connection.BeginTransaction()

    $mealParameterTypes = @{
        ExternalMealId = [System.Data.DbType]::String
        Name           = [System.Data.DbType]::String
        AlternateName  = [System.Data.DbType]::String
        Category       = [System.Data.DbType]::String
        Area           = [System.Data.DbType]::String
        Country        = [System.Data.DbType]::String
        Instructions   = [System.Data.DbType]::String
        ThumbnailUrl   = [System.Data.DbType]::String
        YoutubeUrl     = [System.Data.DbType]::String
        SourceUrl      = [System.Data.DbType]::String
    }

    $insertMealSql = @'
INSERT OR IGNORE INTO Meal (
    ExternalMealId,
    Name,
    AlternateName,
    Category,
    Area,
    Country,
    Instructions,
    ThumbnailUrl,
    YoutubeUrl,
    SourceUrl
)
VALUES (
    @ExternalMealId,
    @Name,
    @AlternateName,
    @Category,
    @Area,
    @Country,
    @Instructions,
    @ThumbnailUrl,
    @YoutubeUrl,
    @SourceUrl
);
'@

    $insertMealCommandParams = @{
        Connection     = $connection
        Transaction    = $transaction
        CommandText    = $insertMealSql
        ParameterTypes = $mealParameterTypes
    }

    $insertMealCommand = New-PreparedSQLiteCommand @insertMealCommandParams
    $commands += $insertMealCommand

    $updateMealSql = @'
UPDATE Meal
SET
    Name          = @Name,
    AlternateName = @AlternateName,
    Category      = @Category,
    Area          = @Area,
    Country       = @Country,
    Instructions  = @Instructions,
    ThumbnailUrl  = @ThumbnailUrl,
    YoutubeUrl    = @YoutubeUrl,
    SourceUrl     = @SourceUrl
WHERE ExternalMealId = @ExternalMealId;
'@

    $updateMealCommandParams = @{
        Connection     = $connection
        Transaction    = $transaction
        CommandText    = $updateMealSql
        ParameterTypes = $mealParameterTypes
    }

    $updateMealCommand = New-PreparedSQLiteCommand @updateMealCommandParams
    $commands += $updateMealCommand

    $selectMealIdSql = @'
SELECT MealId
FROM Meal
WHERE ExternalMealId = @ExternalMealId;
'@

    $selectMealIdCommandParams = @{
        Connection  = $connection
        Transaction = $transaction
        CommandText = $selectMealIdSql
        ParameterTypes = @{
            ExternalMealId = [System.Data.DbType]::String
        }
    }

    $selectMealIdCommand = New-PreparedSQLiteCommand @selectMealIdCommandParams
    $commands += $selectMealIdCommand

    $insertIngredientSql = @'
INSERT OR IGNORE INTO Ingredient (
    Name,
    NormalizedName
)
VALUES (
    @Name,
    @NormalizedName
);
'@

    $insertIngredientCommandParams = @{
        Connection  = $connection
        Transaction = $transaction
        CommandText = $insertIngredientSql
        ParameterTypes = @{
            Name           = [System.Data.DbType]::String
            NormalizedName = [System.Data.DbType]::String
        }
    }

    $insertIngredientCommand = New-PreparedSQLiteCommand @insertIngredientCommandParams
    $commands += $insertIngredientCommand

    $selectIngredientIdSql = @'
SELECT IngredientId
FROM Ingredient
WHERE NormalizedName = @NormalizedName;
'@

    $selectIngredientIdCommandParams = @{
        Connection  = $connection
        Transaction = $transaction
        CommandText = $selectIngredientIdSql
        ParameterTypes = @{
            NormalizedName = [System.Data.DbType]::String
        }
    }

    $selectIngredientIdCommand = New-PreparedSQLiteCommand @selectIngredientIdCommandParams
    $commands += $selectIngredientIdCommand

    $deleteMealIngredientsSql = @'
DELETE FROM MealIngredient
WHERE MealId = @MealId;
'@

    $deleteMealIngredientsCommandParams = @{
        Connection  = $connection
        Transaction = $transaction
        CommandText = $deleteMealIngredientsSql
        ParameterTypes = @{
            MealId = [System.Data.DbType]::Int64
        }
    }

    $deleteMealIngredientsCommand =
        New-PreparedSQLiteCommand @deleteMealIngredientsCommandParams

    $commands += $deleteMealIngredientsCommand

    $insertMealIngredientSql = @'
INSERT INTO MealIngredient (
    MealId,
    IngredientId,
    Position,
    MeasureText
)
VALUES (
    @MealId,
    @IngredientId,
    @Position,
    @MeasureText
);
'@

    $insertMealIngredientCommandParams = @{
        Connection  = $connection
        Transaction = $transaction
        CommandText = $insertMealIngredientSql
        ParameterTypes = @{
            MealId       = [System.Data.DbType]::Int64
            IngredientId = [System.Data.DbType]::Int64
            Position     = [System.Data.DbType]::Int32
            MeasureText  = [System.Data.DbType]::String
        }
    }

    $insertMealIngredientCommand =
        New-PreparedSQLiteCommand @insertMealIngredientCommandParams

    $commands += $insertMealIngredientCommand

    foreach ($meal in $normalizedMeals) {
        if (
            [string]::IsNullOrWhiteSpace($meal.ExternalMealId) -or
            [string]::IsNullOrWhiteSpace($meal.Name)
        ) {
            Write-Warning "Skipping a meal without an ID or name."
            continue
        }

        $mealValues = @{
            ExternalMealId = $meal.ExternalMealId
            Name           = $meal.Name
            AlternateName  = $meal.AlternateName
            Category       = $meal.Category
            Area           = $meal.Area
            Country        = $meal.Country
            Instructions   = $meal.Instructions
            ThumbnailUrl   = $meal.ThumbnailUrl
            YoutubeUrl     = $meal.YoutubeUrl
            SourceUrl      = $meal.SourceUrl
        }

        $setInsertMealValuesParams = @{
            Command = $insertMealCommand
            Values  = $mealValues
        }

        Set-SQLiteCommandValues @setInsertMealValuesParams
        [void]$insertMealCommand.ExecuteNonQuery()

        $setUpdateMealValuesParams = @{
            Command = $updateMealCommand
            Values  = $mealValues
        }

        Set-SQLiteCommandValues @setUpdateMealValuesParams
        [void]$updateMealCommand.ExecuteNonQuery()

        $setSelectMealIdValuesParams = @{
            Command = $selectMealIdCommand
            Values = @{
                ExternalMealId = $meal.ExternalMealId
            }
        }

        Set-SQLiteCommandValues @setSelectMealIdValuesParams

        $mealId = $selectMealIdCommand.ExecuteScalar()

        if ($null -eq $mealId -or $mealId -is [DBNull]) {
            throw "Could not obtain the database ID for meal '$($meal.Name)'."
        }

        $setDeleteMealIngredientsValuesParams = @{
            Command = $deleteMealIngredientsCommand
            Values = @{
                MealId = [long]$mealId
            }
        }

        Set-SQLiteCommandValues @setDeleteMealIngredientsValuesParams
        [void]$deleteMealIngredientsCommand.ExecuteNonQuery()

        foreach ($ingredient in $meal.Ingredients) {
            $setInsertIngredientValuesParams = @{
                Command = $insertIngredientCommand
                Values = @{
                    Name           = $ingredient.Name
                    NormalizedName = $ingredient.NormalizedName
                }
            }

            Set-SQLiteCommandValues @setInsertIngredientValuesParams
            [void]$insertIngredientCommand.ExecuteNonQuery()

            $setSelectIngredientIdValuesParams = @{
                Command = $selectIngredientIdCommand
                Values = @{
                    NormalizedName = $ingredient.NormalizedName
                }
            }

            Set-SQLiteCommandValues @setSelectIngredientIdValuesParams

            $ingredientId = $selectIngredientIdCommand.ExecuteScalar()

            if ($null -eq $ingredientId -or $ingredientId -is [DBNull]) {
                throw "Could not obtain the ID for ingredient '$($ingredient.Name)'."
            }

            $setInsertMealIngredientValuesParams = @{
                Command = $insertMealIngredientCommand
                Values = @{
                    MealId       = [long]$mealId
                    IngredientId = [long]$ingredientId
                    Position     = [int]$ingredient.Position
                    MeasureText  = $ingredient.MeasureText
                }
            }

            Set-SQLiteCommandValues @setInsertMealIngredientValuesParams
            [void]$insertMealIngredientCommand.ExecuteNonQuery()
        }
    }

    $transaction.Commit()

    $summarySql = @'
SELECT
    (SELECT COUNT(*) FROM Meal) AS MealCount,
    (SELECT COUNT(*) FROM Ingredient) AS IngredientCount,
    (SELECT COUNT(*) FROM MealIngredient) AS MealIngredientCount;
'@

    $summaryQueryParams = @{
        SQLiteConnection = $connection
        Query            = $summarySql
        ErrorAction      = "Stop"
    }

    $summary = Invoke-SqliteQuery @summaryQueryParams
    $summary | Format-List
}
catch {
    if ($null -ne $transaction) {
        try {
            $transaction.Rollback()
        }
        catch {
            Write-Warning "The transaction could not be rolled back: $_"
        }
    }

    throw
}
finally {
    foreach ($command in $commands) {
        if ($null -ne $command) {
            $command.Dispose()
        }
    }

    if ($null -ne $transaction) {
        $transaction.Dispose()
    }

    if ($null -ne $connection) {
        $connection.Close()
        $connection.Dispose()
    }
}