package constants

import "github.com/darkphotonKN/community-builds/internal/class"

// --- Classes ---

// -- base --
import "github.com/google/uuid"

var DefaultClasses []class.CreateDefaultClass = []class.CreateDefaultClass{
	{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:        "Warrior",
		Description: "Brutal monster wielding melee weapons.",
		ImageURL:    "Placeholder.",
	},
	{
		ID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Name:        "Sorceror",
		Description: "Master of elemental and arcane magic.",
		ImageURL:    "Placeholder.",
	},
	{
		ID:          uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Name:        "Witch",
		Description: "Dark caster who deals in curses and chaos.",
		ImageURL:    "Placeholder.",
	},
	{
		ID:          uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		Name:        "Monk",
		Description: "A disciplined fighter using martial arts.",
		ImageURL:    "Placeholder.",
	},
	{
		ID:          uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		Name:        "Ranger",
		Description: "Skilled archer with a deep connection to nature.",
		ImageURL:    "Placeholder.",
	},
	{
		ID:          uuid.MustParse("66666666-6666-6666-6666-666666666666"),
		Name:        "Mercenary",
		Description: "Versatile fighter who masters various weapons.",
		ImageURL:    "Placeholder.",
	},
}

// -- ascendancy --
