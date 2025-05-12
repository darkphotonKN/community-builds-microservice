package rating

// --- Validation Maps ---

// -- Rating Category Types --
var validCategoryTypes = map[string]bool{
	"endgame":   true,
	"bossing":   true,
	"speedfarm": true,
	"creative":  true,
	"fun":       true,
}

// --- Validation Helpers ---
func IsValidCategoryType(cagtegoryType string) bool {
	return validCategoryTypes[cagtegoryType]
}
