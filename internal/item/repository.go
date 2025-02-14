package item

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/darkphotonKN/community-builds/internal/utils/errorutils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ItemRepository struct {
	DB *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{
		DB: db,
	}
}

func (r *ItemRepository) CreateItem(createItemReq CreateItemRequest) error {
	query := `
		INSERT INTO items(name, category, class, type, image_url)
		VALUES(:name, :category, :class, :type,  :image_url)
	`

	_, err := r.DB.NamedExec(query, createItemReq)

	fmt.Print("Error when creating item:", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *ItemRepository) GetItems(slot string) (*[]models.Item, error) {
	var items []models.Item

	query := `
	SELECT 
		id,
		image_url, 
		name, 
		category, 
		type, 
		slot, 
		unique_item, 
		class, 
		stats,
		required_level,
		required_intelligence,
		required_strength,
		required_dexterity,
		damage,
		crit,
		aps,
		dps,
		implicit,
		armour,
		evasion,
		energy_shield,
		ward,
		COALESCE(description, '') AS description
	FROM items
	`
	var err error
	if slot != "" {
		query = query + ` WHERE items.slot = $1`
		formatSlot := strings.ToUpper(string(slot[0])) + slot[1:]
		fmt.Println("formatSlot", formatSlot)
		err = r.DB.Select(&items, query, formatSlot)
	} else {
		err = r.DB.Select(&items, query)
	}

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &items, nil
}

func (r *ItemRepository) GetBaseItems() (*[]models.BaseItem, error) {
	var items []models.BaseItem

	query := `
	SELECT * FROM base_items
	`

	err := r.DB.Select(&items, query)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &items, nil
}

func (r *ItemRepository) GetItemMods() (*[]models.ItemMod, error) {
	var items []models.ItemMod

	query := `
	SELECT * FROM item_mods
	`

	err := r.DB.Select(&items, query)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &items, nil
}

func (r *ItemRepository) UpdateItemById(id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error) {
	var item models.Item

	query := `
	UPDATE items
	SET name = :name,
		category = :category,
		type = :type,
		class = :class,
		img_url = :img_url,
	WHERE user_id = :user_id AND id = :id
	RETURNING *;
	`

	params := map[string]interface{}{
		"id":       id,
		"name":     updateItemReq.Name,
		"type":     updateItemReq.Type,
		"category": updateItemReq.Category,
		"class":    updateItemReq.Class,
		"img_url":  updateItemReq.ImageURL,
	}

	rows, err := r.DB.NamedQuery(query, params)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	// loop through and check next table row exists
	if rows.Next() {
		// map the row data to our item struct
		err := rows.StructScan(&item)

		if err != nil {
			return nil, err
		}
	} else {
		// no results found
		return nil, sql.ErrNoRows
	}

	return &item, nil
}

func formatCategory(category string) string {
	// Replace underscores with spaces
	category = strings.ReplaceAll(category, "_", " ")

	// Use cases.Title for title-casing each word
	caser := cases.Title(language.English)
	words := strings.Fields(category)
	for i, word := range words {
		words[i] = caser.String(word)
	}

	// Join the words back into a single string
	return strings.Join(words, " ")
}

func checkStr(items []string, text string) int {

	strIndex := -1
	for index, item := range items {
		// fmt.Println("item:", item, "text:", text)

		if strings.Contains(item, text) {
			// fmt.Printf("Found a match in: %s\n", item)
			strIndex = index
			break
		}
	}
	return strIndex
}

// checking unique items is exist
func (r *ItemRepository) CheckUniqueItemExist() bool {
	query := `SELECT EXISTS (SELECT 1 FROM items WHERE unique_item = true LIMIT 1 )`

	var exists bool
	err := r.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking item existence: %v\n", err)
		return false
	}
	return exists
}

// checking base items is exist
func (r *ItemRepository) CheckBaseItemExist() bool {
	query := `SELECT EXISTS (SELECT 1 FROM base_items LIMIT 1 )`

	var exists bool
	err := r.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking base item existence: %v\n", err)
		return false
	}
	return exists
}

// checking base items is exist
func (r *ItemRepository) CheckItemModExist() bool {
	query := `SELECT EXISTS (SELECT 1 FROM item_mods LIMIT 1 )`

	var exists bool
	err := r.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking item mod existence: %v\n", err)
		return false
	}
	return exists
}

func (r *ItemRepository) AddUniqueItems(tx *sqlx.Tx, items *[]models.Item) error {

	generateCustomUUID := func(baseUUID string, sequence int) (*uuid.UUID, error) {

		// Format sequence as 4 digits
		suffix := fmt.Sprintf("%04d", sequence)

		// Replace the last 4 digits of the base UUID
		customUUID := baseUUID[:len(baseUUID)-4] + suffix
		parseUuid, err := uuid.Parse(customUUID)
		if err != nil {
			return nil, err
		}
		return &parseUuid, nil
	}

	stmt, err := tx.Prepare(pq.CopyIn(
		"items",
		"id",
		"image_url",
		"category",
		"class",
		"name",
		"type",
		"unique_item",
		"slot",
		"description",
		"required_level",
		"required_strength",
		"required_dexterity",
		"required_intelligence",
		// armor
		"armour",
		"energy_shield",
		"evasion",
		"block",
		"ward",
		// weapon
		"damage",
		"aps",
		"crit",
		"pdps",
		"edps",
		"dps",
		// flask
		"life",
		"mana",
		"duration",
		"usage",
		"capacity",
		// common
		"stats",
		"additional",
	))
	if err != nil {
		return err
	}
	// fixed uuid
	baseUUID := "11111111-1111-1111-1111-111111110000"
	for index, item := range *items {
		uuid, err := generateCustomUUID(baseUUID, index)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			uuid,
			item.ImageUrl,
			item.Category,
			item.Class,
			item.Name,
			item.Type,
			item.UniqueItem,
			item.Slot,
			item.Description,
			item.RequiredLevel,
			item.RequiredStrength,
			item.RequiredDexterity,
			item.RequiredIntelligence,
			// armor
			item.Armour,
			item.EnergyShield,
			item.Evasion,
			item.Block,
			item.Ward,
			// weapon
			item.Damage,
			item.APS,
			item.Crit,
			item.PDPS,
			item.EDPS,
			item.DPS,
			// flask
			item.Life,
			item.Mana,
			item.Duration,
			item.Usage,
			item.Capacity,
			// common
			pq.Array(item.Stats),
			item.Additional,
		)
		if err != nil {
			stmt.Close()
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) AddBaseItems(tx *sqlx.Tx, items *[]models.BaseItem) error {

	generateCustomUUID := func(baseUUID string, sequence int) (*uuid.UUID, error) {

		// Format sequence as 4 digits
		suffix := fmt.Sprintf("%04d", sequence)

		// Replace the last 4 digits of the base UUID
		customUUID := baseUUID[:len(baseUUID)-4] + suffix
		parseUuid, err := uuid.Parse(customUUID)
		if err != nil {
			return nil, err
		}
		return &parseUuid, nil
	}

	stmt, err := tx.Prepare(pq.CopyIn(
		"base_items",
		"id",
		"image_url",
		"category",
		"class",
		"name",
		"type",
		"equip_type",
		"is_two_hands",
		"slot",
		"required_level",
		"required_strength",
		"required_dexterity",
		"required_intelligence",
		// armor
		"armour",
		"energy_shield",
		"evasion",
		"ward",
		// weapon
		"damage",
		"aps",
		"crit",
		"dps",
		// common
		"implicit",
	))
	if err != nil {
		return err
	}

	// fixed uuid
	baseUUID := "11111111-1111-1111-1111-111111120000"
	for index, item := range *items {
		uuid, err := generateCustomUUID(baseUUID, index)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			uuid,
			item.ImageUrl,
			item.Category,
			item.Class,
			item.Name,
			item.Type,
			item.EquipType,
			item.IsTwoHands,
			item.Slot,
			item.RequiredLevel,
			item.RequiredStrength,
			item.RequiredDexterity,
			item.RequiredIntelligence,
			// armor
			item.Armour,
			item.EnergyShield,
			item.Evasion,
			item.Ward,
			// weapon
			item.Damage,
			item.APS,
			item.Crit,
			item.DPS,
			// common
			pq.Array(item.Implicit),
		)
		if err != nil {
			stmt.Close()
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) GetBaseItemById(id uuid.UUID) (*models.BaseItem, error) {
	fmt.Println("base id", id)
	var baseItem models.BaseItem
	query := `
    SELECT *
    FROM base_items 
    WHERE id = $1
    `
	// var rawImplicits []byte
	err := r.DB.Get(&baseItem, query, id)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	// if err := json.Unmarshal(rawImplicits, &baseItem.Implicit); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal implicit: %w", err)
	// }
	return &baseItem, nil
}

func (r *ItemRepository) AddItemMods(tx *sqlx.Tx, items *[]models.ItemMod) error {

	generateCustomUUID := func(baseUUID string, sequence int) (*uuid.UUID, error) {

		// Format sequence as 4 digits
		suffix := fmt.Sprintf("%04d", sequence)

		// Replace the last 4 digits of the base UUID
		customUUID := baseUUID[:len(baseUUID)-4] + suffix
		parseUuid, err := uuid.Parse(customUUID)
		if err != nil {
			return nil, err
		}
		return &parseUuid, nil
	}

	stmt, err := tx.Prepare(pq.CopyIn(
		"item_mods",
		"id",
		"affix",
		"name",
		"level",
		"stat",
		"tags",
	))
	if err != nil {
		return err
	}

	// fixed uuid
	baseUUID := "11111111-1111-1111-1111-111111130000"
	for index, item := range *items {
		uuid, err := generateCustomUUID(baseUUID, index)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			uuid,
			item.Affix,
			item.Name,
			item.Level,
			item.Stat,
			item.Tags,
		)

		if err != nil {
			stmt.Close()
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) CreateRareItem(id uuid.UUID, createRareItemReq CreateRareItemReq) (*uuid.UUID, error) {
	baseItem, err := r.GetBaseItemById(createRareItemReq.BaseItemId)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	query := `
		INSERT INTO items(
		base_item_id,
		image_url, 
		name, 
		category, 
		type, 
		slot, 
		unique_item, 
		class, 
		stats,
		required_level,
		required_intelligence,
		required_strength,
		required_dexterity,
		damage,
		crit,
		aps,
		dps,
		implicit,
		armour,
		evasion,
		energy_shield,
		ward
		)
		VALUES(
		:base_item_id,
		:image_url, 
		:name, 
		:category, 
		:type, 
		:slot, 
		:unique_item, 
		:class, 
		:stats,
		:required_level,
		:required_intelligence,
		:required_strength,
		:required_dexterity,
		:damage,
		:crit,
		:aps,
		:dps,
		:implicit,
		:armour,
		:evasion,
		:energy_shield,
		:ward
		)
		RETURNING id;
	`

	payload := map[string]interface{}{
		"base_item_id":          createRareItemReq.BaseItemId,
		"image_url":             baseItem.ImageUrl,
		"name":                  createRareItemReq.Name,
		"category":              baseItem.Category,
		"type":                  baseItem.Type,
		"slot":                  baseItem.Slot,
		"unique_item":           false,
		"class":                 baseItem.Class,
		"stats":                 pq.StringArray(createRareItemReq.Stats),
		"required_level":        baseItem.RequiredLevel,
		"required_intelligence": baseItem.RequiredIntelligence,
		"required_strength":     baseItem.RequiredStrength,
		"required_dexterity":    baseItem.RequiredDexterity,
		"damage":                baseItem.Damage,
		"crit":                  baseItem.Crit,
		"aps":                   baseItem.APS,
		"dps":                   baseItem.DPS,
		"implicit":              pq.StringArray(baseItem.Implicit),
		"armour":                baseItem.Armour,
		"evasion":               baseItem.Evasion,
		"energy_shield":         baseItem.EnergyShield,
		"ward":                  baseItem.Ward,
	}
	rows, createErr := r.DB.NamedQuery(query, payload)
	var lastInsertID uuid.UUID

	if rows.Next() {
		err := rows.Scan(&lastInsertID)

		if err != nil {
			return nil, err
		}
	}

	if createErr != nil {
		return nil, errorutils.AnalyzeDBErr(createErr)
	}

	return &lastInsertID, nil
}

func (r *ItemRepository) CreateRareItemToList(id uuid.UUID, createRareItemReq CreateRareItemReq) (*uuid.UUID, error) {
	baseItem, err := r.GetBaseItemById(createRareItemReq.BaseItemId)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	query := `
		INSERT INTO items(
		member_id, 
		base_item_id,
		image_url, 
		name, 
		category, 
		type, 
		slot, 
		unique_item, 
		class, 
		stats,
		required_level,
		required_intelligence,
		required_strength,
		required_dexterity,
		damage,
		crit,
		aps,
		dps,
		implicit,
		armour,
		evasion,
		energy_shield,
		ward
		)
		VALUES(
		:member_id, 
		:base_item_id,
		:image_url, 
		:name, 
		:category, 
		:type, 
		:slot, 
		:unique_item, 
		:class, 
		:stats,
		:required_level,
		:required_intelligence,
		:required_strength,
		:required_dexterity,
		:damage,
		:crit,
		:aps,
		:dps,
		:implicit,
		:armour,
		:evasion,
		:energy_shield,
		:ward
		)
		RETURNING id;
	`

	payload := map[string]interface{}{
		"member_id":             id,
		"base_item_id":          createRareItemReq.BaseItemId,
		"image_url":             baseItem.ImageUrl,
		"name":                  createRareItemReq.Name,
		"category":              baseItem.Category,
		"type":                  baseItem.Type,
		"slot":                  baseItem.Slot,
		"unique_item":           false,
		"class":                 baseItem.Class,
		"stats":                 pq.StringArray(createRareItemReq.Stats),
		"required_level":        baseItem.RequiredLevel,
		"required_intelligence": baseItem.RequiredIntelligence,
		"required_strength":     baseItem.RequiredStrength,
		"required_dexterity":    baseItem.RequiredDexterity,
		"damage":                baseItem.Damage,
		"crit":                  baseItem.Crit,
		"aps":                   baseItem.APS,
		"dps":                   baseItem.DPS,
		"implicit":              pq.StringArray(baseItem.Implicit),
		"armour":                baseItem.Armour,
		"evasion":               baseItem.Evasion,
		"energy_shield":         baseItem.EnergyShield,
		"ward":                  baseItem.Ward,
	}
	rows, createErr := r.DB.NamedQuery(query, payload)
	var lastInsertID uuid.UUID

	if rows.Next() {
		err := rows.Scan(&lastInsertID)

		if err != nil {
			return nil, err
		}
	}

	if createErr != nil {
		return nil, errorutils.AnalyzeDBErr(createErr)
	}

	return &lastInsertID, nil
}
