package validation

// Valid categories for items
var validCategories = map[string]bool{
	"Flasks":            true,
	"Other":             true,
	"Jewellery":         true,
	"One Handed Weapon": true,
	"Two Handed Weapon": true,
	"Off-hand":          true,
	"Armor":             true,
}

// Valid classes for items
var validClasses = map[string]bool{
	// Flasks
	"Life Flasks":   true,
	"Mana Flasks":   true,
	"Hybrid Flasks": true,

	// Jewellery
	"Amulets": true,
	"Rings":   true,
	"Belts":   true,

	// One Handed Weapon
	"Claws":                     true,
	"Daggers":                   true,
	"Wands":                     true,
	"One Hand Swords":           true,
	"Thrusting One Hand Swords": true,
	"One Hand Axes":             true,
	"One Hand Maces":            true,
	"Sceptres":                  true,

	// Two Handed Weapon
	"Bows":            true,
	"Staves":          true,
	"Two Hand Swords": true,
	"Two Hand Axes":   true,
	"Two Hand Maces":  true,

	// Off-hand
	"Quivers": true,
	"Shields": true,

	// Armor
	"Gloves":       true,
	"Boots":        true,
	"Body Armours": true,
	"Helmets":      true,
}

// Valid types for items
var validTypes = map[string]bool{
	"Glorious Plate": true,
	"Gemstone":       true,
}

// Helper functions to check validity
func IsValidCategory(category string) bool {
	return validCategories[category]
}

func IsValidClass(class string) bool {
	return validClasses[class]
}

func IsValidType(itemType string) bool {
	return validTypes[itemType]
}
