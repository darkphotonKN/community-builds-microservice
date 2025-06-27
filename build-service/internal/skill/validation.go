package skill

// --- Validation Maps ---

// -- Skill Types --
var validTypes = map[string]bool{
	"active":  true,
	"support": true,
}

// --- Validation Helpers ---
func IsValidType(skillType string) bool {
	return validTypes[skillType]
}
