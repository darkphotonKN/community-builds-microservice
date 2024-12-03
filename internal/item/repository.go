package item

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
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

func (r *ItemRepository) GetItems() (*[]models.Item, error) {
	var items []models.Item

	query := `
	SELECT * FROM items
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

// , itemsCh chan WikiItem, wg *sync.WaitGroup
func getCategoryItem(category string, itemsCh chan WikiItem, wg *sync.WaitGroup) {

	resp, err := http.Get("https://www.poewiki.net/wiki/List_of_unique_" + category)
	if err != nil {
		panic(err)
	}
	defer func() {
		resp.Body.Close()
		wg.Done()
	}()

	// fmt.Println("Response status:", resp.Status)

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".item-table").Each(func(tableIndex int, table *goquery.Selection) {

		var currentItemThs []string

		currentItemThs = append(currentItemThs, "Item")
		table.Find("th").Each(func(thIndex int, item *goquery.Selection) {
			// text := item.Text()
			fmt.Println("item:", item.Text())
			titleAttr, _ := item.Find("abbr").Attr("title")
			// fmt.Println("Href:", attr)

			if strings.Contains(titleAttr, "Required level") {
				// fmt.Println("包含子字串:", "strength")
				currentItemThs = append(currentItemThs, "RequiredLevel")
			}
			if strings.Contains(titleAttr, "Required strength") {
				// fmt.Println("包含子字串:", "strength")
				currentItemThs = append(currentItemThs, "Strength")
			}

			if strings.Contains(titleAttr, "Required dexterity") {
				// fmt.Println("包含子字串:", "dexterity")
				currentItemThs = append(currentItemThs, "Dexterity")
			}

			if strings.Contains(titleAttr, "Required intelligence") {
				// fmt.Println("包含子字串:", "intelligence")
				currentItemThs = append(currentItemThs, "Intelligence")
			}

			if strings.Contains(titleAttr, "Armour") {
				// fmt.Println("包含子字串:", "Armour")
				currentItemThs = append(currentItemThs, "Armour")
			}

			if strings.Contains(titleAttr, "Energy shield") {
				// fmt.Println("包含子字串:", "Energy shield")
				currentItemThs = append(currentItemThs, "EnergyShield")
			}

			if strings.Contains(titleAttr, "Evasion rating") {
				// fmt.Println("包含子字串:", "Evasion rating")
				currentItemThs = append(currentItemThs, "Evasion")
			}

			if strings.Contains(titleAttr, "Chance to block") {
				// fmt.Println("包含子字串:", "Armour")
				currentItemThs = append(currentItemThs, "Block")
			}

			if strings.Contains(titleAttr, "Ward") {
				currentItemThs = append(currentItemThs, "Ward")
			}
			// weapon
			if strings.Contains(titleAttr, "Colour coded damage") {
				currentItemThs = append(currentItemThs, "Damage")
			}
			if strings.Contains(titleAttr, "Attacks per second") {
				currentItemThs = append(currentItemThs, "APS")
			}
			if strings.Contains(titleAttr, "Local weapon critical strike chance") {
				currentItemThs = append(currentItemThs, "Crit")
			}
			if strings.Contains(titleAttr, "physical damage per second") {
				currentItemThs = append(currentItemThs, "pDPS")
			}
			if strings.Contains(titleAttr, "elemental damage") {
				currentItemThs = append(currentItemThs, "eDPS")
			}
			if strings.Contains(titleAttr, "total damage") {
				currentItemThs = append(currentItemThs, "DPS")
			}
			// flask
			if strings.Contains(titleAttr, "Life regenerated over the flask duration") {
				currentItemThs = append(currentItemThs, "Life")
			}

			if strings.Contains(titleAttr, "Mana regenerated over the flask duration") {
				currentItemThs = append(currentItemThs, "Mana")
			}

			if strings.Contains(item.Text(), "Duration") {
				currentItemThs = append(currentItemThs, "Duration")
			}

			if strings.Contains(titleAttr, "Number of charges consumed on use") {
				currentItemThs = append(currentItemThs, "Usage")
			}
			if strings.Contains(titleAttr, "Maximum number of flask charges held") {
				currentItemThs = append(currentItemThs, "Capacity")
			}
			if strings.Contains(item.Text(), "Stats") {
				currentItemThs = append(currentItemThs, "Stats")
			}
			if strings.Contains(item.Text(), "Additionaldrop restrictions") {
				currentItemThs = append(currentItemThs, "Additional")
			}
		})
		// for _, th := range currentItemThs {
		// 	fmt.Println("th", th)
		// }
		table.Find("tbody tr").Each(func(trIndex int, tr *goquery.Selection) {
			wg.Add(1)
			go getItem(currentItemThs, trIndex, tr, itemsCh, wg, category)

		})
	})

	// go func() {
	// 	wg.Wait()
	// 	close(itemsCh) // 所有 goroutine 完成後才關閉 Channel
	// }()

	// for itemCh := range itemsCh {
	// 	items = append(items, itemCh)
	// }

}

func getItem(currentItemThs []string, index int, tr *goquery.Selection, itemsCh chan WikiItem, wg *sync.WaitGroup, category string) {
	defer wg.Done()
	// go getItem()
	myItem := WikiItem{}
	tr.Find("td").Each(func(tdIndex int, td *goquery.Selection) {

		if columnIndex := checkStr(currentItemThs, "Item"); columnIndex == tdIndex {
			aTag := td.Find("a").First()
			myItem.Name = aTag.Text()
			aTagUrl, _ := aTag.Attr("href")
			resp, err := http.Get("https://www.poewiki.net/" + aTagUrl)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				panic(err)
			}

			// first item box info
			itemBox := doc.Find(".item-box.-unique").Children()
			html, _ := itemBox.Html()
			// fmt.Println("itemBox", html)
			lines := strings.Split(html, "<br/>")
			if len(lines) > 1 {
				value, _ := goquery.NewDocumentFromReader(strings.NewReader(lines[1]))
				myItem.Type = value.Text()
			}

			desc := doc.Find(".group.tc.-flavour").First()
			// fmt.Println("desc", desc)
			myItem.Description = desc.Text()
			imgTag := td.Find("img").First()
			imgUrl, _ := imgTag.Attr("src")
			myItem.ImageUrl = imgUrl

			// second item box info
			itemBox = doc.Find("div.item-box.-unique .tc.-value").Last()
			// value = itemBox.Find(".tc.-value").Last()
			html, _ = itemBox.Html()
			// fmt.Println("html", html)
			myItem.Class = html
			// if value == "" {
			// 	fmt.Println("itemBox", itemBox)
			// }

		}
		// fmt.Println("currentItemThs", currentItemThs)
		if columnIndex := checkStr(currentItemThs, "RequiredLevel"); columnIndex == tdIndex {
			myItem.RequiredLevel = td.Text()
		}
		// armor
		if columnIndex := checkStr(currentItemThs, "Strength"); columnIndex == tdIndex {
			myItem.RequiredStrength = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Dexterity"); columnIndex == tdIndex {
			myItem.RequiredDexterity = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Intelligence"); columnIndex == tdIndex {
			myItem.RequiredIntelligence = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Armour"); columnIndex == tdIndex {
			myItem.Armour = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "EnergyShield"); columnIndex == tdIndex {
			myItem.EnergyShield = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Evasion"); columnIndex == tdIndex {
			myItem.Evasion = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Block"); columnIndex == tdIndex {
			myItem.Block = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Ward"); columnIndex == tdIndex {
			myItem.Ward = td.Text()
		}
		// weapon
		if columnIndex := checkStr(currentItemThs, "Damage"); columnIndex == tdIndex {
			myItem.Damage = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "APS"); columnIndex == tdIndex {
			myItem.APS = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Crit"); columnIndex == tdIndex {
			myItem.Crit = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "pDPS"); columnIndex == tdIndex {
			myItem.PDPS = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "eDPS"); columnIndex == tdIndex {
			myItem.EDPS = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "DPS"); columnIndex == tdIndex {
			myItem.DPS = td.Text()
		}

		// flask
		if columnIndex := checkStr(currentItemThs, "Life"); columnIndex == tdIndex {
			myItem.Life = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Mana"); columnIndex == tdIndex {
			myItem.Mana = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Duration"); columnIndex == tdIndex {
			myItem.Duration = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Usage"); columnIndex == tdIndex {
			myItem.Usage = td.Text()
		}
		if columnIndex := checkStr(currentItemThs, "Capacity"); columnIndex == tdIndex {
			myItem.Capacity = td.Text()
		}
		// common
		if columnIndex := checkStr(currentItemThs, "Stats"); columnIndex == tdIndex {
			content, _ := td.Html()
			lines := strings.Split(content, "<br/>")

			// 移除空行並修剪空白
			var cleanedLines []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					lineHtml, _ := goquery.NewDocumentFromReader(strings.NewReader(line))
					cleanedLines = append(cleanedLines, lineHtml.Text())
				}
			}
			myItem.Stats = cleanedLines
		}

		if columnIndex := checkStr(currentItemThs, "Additional"); columnIndex == tdIndex {
			myItem.Additional = td.Text()
		}
		myItem.Category = formatCategory(category)
	})

	if myItem.Name != "" && myItem.ImageUrl != "" {
		// *items = append(*items, myItem)
		itemsCh <- myItem
	}

}

func (r *ItemRepository) GetWikiItems() (*[]WikiItem, error) {
	// var items []models.Item

	// query := `
	// SELECT * FROM items
	// `

	// err := r.DB.Select(&items, query)

	// if err != nil {
	// 	return nil, errorutils.AnalyzeDBErr(err)
	// }

	uniqueCategories := []string{
		// Weapon
		"axes",
		"bows",
		"quivers",
		"claws",
		"daggers",
		"fishing_rods",
		"maces",
		"sceptres",
		"staves",
		"swords",
		"wands",
		// Armor
		"body_armors",
		"boots",
		"gloves",
		"helmets",
		"shields",
		// Jewellery,
		"amulets",
		"belts",
		"rings",
		// Flask
		"life_flasks",
		"mana_flasks",
		"hybrid_flasks",
		"utility_flasks",
	}

	// Process each category
	// formattedCategories := make([]string, len(uniqueCategories))
	// for i, category := range uniqueCategories {
	// 	formattedCategories[i] = formatCategory(category)
	// }

	var wg sync.WaitGroup
	itemsCh := make(chan WikiItem)
	items := []WikiItem{}

	// for _, category := range formattedCategories {
	// 	wg.Add(1)
	// 	go getCategoryItem(category, itemsCh, &wg)
	// }
	for _, category := range uniqueCategories {
		wg.Add(1)
		go getCategoryItem(category, itemsCh, &wg)
	}
	go func() {
		wg.Wait()
		close(itemsCh) // 所有 goroutine 完成後才關閉 Channel
	}()

	for itemCh := range itemsCh {
		items = append(items, itemCh)
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(pq.CopyIn(
		"items",
		"image_url",
		"category",
		"class",
		"name",
		"type",
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
		return nil, err
	}

	for _, item := range items {
		_, err := stmt.Exec(
			item.ImageUrl,
			item.Category,
			item.Class,
			item.Name,
			item.Type,
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
			tx.Rollback()
			return nil, err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return &items, nil
}

// get base items

func getBaseItemEquipType(equipType string, itemsCh chan BaseItem, wg *sync.WaitGroup) {

	// 建立上下文
	ctx, cancelChromedp := chromedp.NewContext(context.Background())
	defer cancelChromedp() // 釋放資源

	// 設定超時
	ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()
	defer wg.Done()

	var pageHTML string

	// 模擬瀏覽器進入網站並抓取內容
	err := chromedp.Run(ctxWithTimeout,
		// 打開指定 URL
		chromedp.Navigate("https://www.pathofexile.com/item-data/"+equipType),
		// 等待網頁載入完成
		chromedp.WaitReady("body"),
		// 抓取完整 HTML
		chromedp.OuterHTML("html", &pageHTML),
	)

	// 處理錯誤
	if err != nil {
		log.Fatalf("Failed to get HTML: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(pageHTML)))
	if err != nil {
		panic(err)
	}

	doc.Find(".layoutBox1.layoutBoxStretch").Each(func(index int, box *goquery.Selection) {
		// boxHtml, _ := box.Html()
		// fmt.Println("boxHtml", boxHtml)
		wg.Add(1)
		go getBaseItemTable(box, itemsCh, wg)
	})

}

func getBaseItemTable(table *goquery.Selection, itemsCh chan BaseItem, wg *sync.WaitGroup) {
	defer wg.Done()

	category := table.Find("h1").First()
	// fmt.Println("category", category.Text())

	table.Find("table.itemDataTable").Each(func(index int, table *goquery.Selection) {

		// tableHtml, _ := table.Html()
		// fmt.Println("tableHtml", tableHtml)
		thList := []string{}
		table.Find("thead tr th").Each(func(index int, th *goquery.Selection) {
			// thHtml, _ := th.Html()
			// fmt.Println("thHtml", thHtml)
			if index == 0 {
				thList = append(thList, "ImageUrl")
			} else {
				thList = append(thList, th.Text())
			}
		})
		// for _, th := range thList {
		// 	fmt.Println("th", th)
		// }
		var baseItem BaseItem
		table.Find("tbody tr").Each(func(trIndex int, tr *goquery.Selection) {
			if trIndex%2 == 0 {
				baseItem = BaseItem{
					Category: category.Text(),
					Type:     category.Text(),
					Class:    category.Text(),
				}
			}
			if trIndex%2 == 0 {
				tr.Find("td").Each(func(tdIndex int, td *goquery.Selection) {
					// fmt.Println("td text", td.Text())
					if columnIndex := checkStr(thList, "ImageUrl"); columnIndex == tdIndex {
						img := td.Find("img").First()
						imgPath, _ := img.Attr("src")
						baseItem.ImageUrl = imgPath
					}
					if columnIndex := checkStr(thList, "Name"); columnIndex == tdIndex {
						baseItem.Name = td.Text()
					}
					if columnIndex := checkStr(thList, "Level"); columnIndex == tdIndex {
						baseItem.RequiredLevel = td.Text()
					}
					if columnIndex := checkStr(thList, "Str"); columnIndex == tdIndex {
						baseItem.RequiredStrength = td.Text()
					}
					if columnIndex := checkStr(thList, "Dex"); columnIndex == tdIndex {
						baseItem.RequiredDexterity = td.Text()
					}
					if columnIndex := checkStr(thList, "Int"); columnIndex == tdIndex {
						baseItem.RequiredIntelligence = td.Text()
					}
					if columnIndex := checkStr(thList, "Damage"); columnIndex == tdIndex {
						baseItem.Damage = td.Text()
					}
					if columnIndex := checkStr(thList, "Critical Chance"); columnIndex == tdIndex {
						baseItem.Crit = td.Text()
					}
					if columnIndex := checkStr(thList, "APS"); columnIndex == tdIndex {
						baseItem.APS = td.Text()
					}
					if columnIndex := checkStr(thList, "DPS"); columnIndex == tdIndex {
						baseItem.DPS = td.Text()
					}

					if columnIndex := checkStr(thList, "Armour"); columnIndex == tdIndex {
						baseItem.Armour = td.Text()
					}
					if columnIndex := checkStr(thList, "Evasion Rating"); columnIndex == tdIndex {
						baseItem.Evasion = td.Text()
					}
					if columnIndex := checkStr(thList, "Energy Shield"); columnIndex == tdIndex {
						baseItem.EnergyShield = td.Text()
					}
					if columnIndex := checkStr(thList, "Ward"); columnIndex == tdIndex {
						baseItem.Ward = td.Text()
					}

				})
			} else {
				td := tr.Find("td").First()
				baseItem.Stats = append(baseItem.Stats, strings.TrimSpace(td.Text()))
				itemsCh <- baseItem
			}
		})
	})

}

func (r *ItemRepository) GetBaseItems() (*[]BaseItem, error) {

	var wg sync.WaitGroup
	itemsCh := make(chan BaseItem)
	items := []BaseItem{}

	wg.Add(2)
	go getBaseItemEquipType("weapon", itemsCh, &wg)
	go getBaseItemEquipType("armour", itemsCh, &wg)

	go func() {
		wg.Wait()
		close(itemsCh) // 所有 goroutine 完成後才關閉 Channel
	}()

	for itemCh := range itemsCh {
		items = append(items, itemCh)
	}

	// store items to db

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(pq.CopyIn(
		"base_items",

		"image_url",
		"category",
		"class",
		"name",
		"type",
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
		"stats",
	))
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		_, err := stmt.Exec(
			item.ImageUrl,
			item.Category,
			item.Class,
			item.Name,
			item.Type,
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
			pq.Array(item.Stats),
		)
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return nil, err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return &items, nil
}

func getModHtml(itemsCh chan ModItem, wg *sync.WaitGroup) {

	// 建立上下文
	ctx, cancelChromedp := chromedp.NewContext(context.Background())
	defer cancelChromedp() // 釋放資源

	// 設定超時
	ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()
	defer wg.Done()

	var pageHTML string

	// 模擬瀏覽器進入網站並抓取內容
	err := chromedp.Run(ctxWithTimeout,
		// 打開指定 URL
		chromedp.Navigate("https://www.pathofexile.com/item-data/mods"),
		// 等待網頁載入完成
		chromedp.WaitReady("body"),
		// 抓取完整 HTML
		chromedp.OuterHTML("html", &pageHTML),
	)

	// 處理錯誤
	if err != nil {
		log.Fatalf("Failed to get HTML: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(pageHTML)))
	if err != nil {
		panic(err)
	}

	table := doc.Find(".layoutBox1.layoutBoxStretch").First()

	table.Find("table.itemDataTable tbody tr").Each(func(index int, tr *goquery.Selection) {
		// boxHtml, _ := box.Html()
		// fmt.Println("boxHtml", boxHtml)
		wg.Add(1)
		go getModItem(tr, itemsCh, wg)
	})
}

func getModItem(tr *goquery.Selection, itemsCh chan ModItem, wg *sync.WaitGroup) {
	defer wg.Done()
	modItem := ModItem{}
	thList := []string{"Affix", "Name", "Level", "Stat", "Tags"}

	tr.Find("td").Each(func(tdIndex int, td *goquery.Selection) {
		if columnIndex := checkStr(thList, "Affix"); columnIndex == tdIndex {
			modItem.Affix = td.Text()
		}
		if columnIndex := checkStr(thList, "Name"); columnIndex == tdIndex {
			text := td.Text()
			text = strings.Replace(text, "of", "", 1)
			text = strings.TrimSpace(text)
			modItem.Name = text
		}
		if columnIndex := checkStr(thList, "Level"); columnIndex == tdIndex {
			modItem.Level = td.Text()
		}
		if columnIndex := checkStr(thList, "Stat"); columnIndex == tdIndex {
			modItem.Stat = td.Text()
		}
		if columnIndex := checkStr(thList, "Tags"); columnIndex == tdIndex {
			modItem.Tags = td.Text()
		}
	})

	itemsCh <- modItem
}

func (r *ItemRepository) GetModItems() (*[]ModItem, error) {

	var wg sync.WaitGroup
	itemsCh := make(chan ModItem)
	items := []ModItem{}

	wg.Add(1)
	go getModHtml(itemsCh, &wg)

	go func() {
		wg.Wait()
		close(itemsCh) // 所有 goroutine 完成後才關閉 Channel
	}()

	for item := range itemsCh {
		items = append(items, item)
	}

	// store items to db

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(pq.CopyIn(
		"mod_items",

		"affix",
		"name",
		"level",
		"stat",
		"tags",
	))
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		_, err := stmt.Exec(
			item.Affix,
			item.Name,
			item.Level,
			item.Stat,
			item.Tags,
		)

		if err != nil {
			stmt.Close()
			tx.Rollback()
			return nil, err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return &items, nil
}
