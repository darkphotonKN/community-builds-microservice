package item

// --- Validation Maps ---

// -- Item Categories --
var validCategories = map[string]bool{
	"Flask":             true,
	"Other":             true,
	"Jewellery":         true,
	"One Handed Weapon": true,
	"Two Handed Weapon": true,
	"Off-hand":          true,
	"Armor":             true,
}

// -- Item Classes --
var validClasses = map[string]bool{
	// Flasks
	"Life Flasks":    true,
	"Mana Flasks":    true,
	"Hybrid Flasks":  true,
	"Utility Flasks": true,

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

// -- Item Types --
var validTypes = map[string]bool{
	"Glorious Plate": true,
	"Vaal Axe":       true,
}

// --- Validation Helpers ---
func IsValidCategory(category string) bool {
	return validCategories[category]
}

func IsValidClass(class string) bool {
	return validClasses[class]
}

func IsValidType(itemType string) bool {
	return validTypes[itemType]
}
