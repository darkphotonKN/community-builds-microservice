package item

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/darkphotonKN/community-builds/internal/utils/dbutils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemService struct {
	Repo *ItemRepository
}

func NewItemService(repo *ItemRepository) *ItemService {
	return &ItemService{
		Repo: repo,
	}
}

func (s *ItemService) CreateItemService(createItemReq CreateItemRequest) error {
	return s.Repo.CreateItem(createItemReq)
}

func (s *ItemService) AddItemToBuildService(memberId uuid.UUID, item CreateItemRequest) error {
	return s.Repo.CreateItem(item)
}

func (s *ItemService) GetItemsService(slot string) (*[]models.Item, error) {
	return s.Repo.GetItems(slot)
}

func (s *ItemService) GetBaseItemsService() (*[]models.BaseItem, error) {
	return s.Repo.GetBaseItems()
}

func (s *ItemService) GetItemModsService() (*[]models.ItemMod, error) {
	return s.Repo.GetItemMods()
}

func (s *ItemService) UpdateItemsService(id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error) {
	return s.Repo.UpdateItemById(id, updateItemReq)
}

func (s *ItemService) GetBaseItemByIdService(id uuid.UUID) (*models.BaseItem, error) {
	return s.Repo.GetBaseItemById(id)
}

func (s *ItemService) CreateRareItemService(id uuid.UUID, createRareItemReq CreateRareItemReq) error {
	return s.Repo.CreateRareItem(id, createRareItemReq)
}

func getCategoryItem(category string, itemsCh chan models.Item, wg *sync.WaitGroup) {

	resp, err := http.Get("https://www.poewiki.net/wiki/List_of_unique_" + category)
	if err != nil {
		panic(err)
	}
	defer func() {
		resp.Body.Close()
		wg.Done()
	}()

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
		table.Find("tbody tr").Each(func(trIndex int, tr *goquery.Selection) {
			wg.Add(1)
			go getItem(currentItemThs, trIndex, tr, itemsCh, wg, category)

		})
	})
}

func getItem(currentItemThs []string, index int, tr *goquery.Selection, itemsCh chan models.Item, wg *sync.WaitGroup, category string) {
	defer wg.Done()
	myItem := models.Item{}
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
			myItem.ImageUrl = "https://www.poewiki.net" + imgUrl

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
		slot := "Weapon"
		switch category {
		case "shields":
			slot = "Shield"
		case "boots":
			slot = "Boots"
		case "gloves":
			slot = "Gloves"
		case "helmets":
			slot = "Helmet"
		case "amulets":
			slot = "Amulet"
		case "body_armours":
			slot = "Body Armour"
		case "belts":
			slot = "Belt"
		case "rings":
			slot = "Rings"
		case "life_flasks":
		case "mana_flasks":
		case "hybrid_flasks":
		case "utility_flasks":
			slot = "Flask"
		default:
			slot = "Weapon"
		}

		myItem.Slot = slot
		myItem.Category = formatCategory(category)
		myItem.UniqueItem = true

	})

	if myItem.Name != "" && myItem.ImageUrl != "" {
		// *items = append(*items, myItem)
		itemsCh <- myItem
	}

}

func (s *ItemService) CrawlingAndAddUniqueItemsService() error {
	return dbutils.ExecTx(s.Repo.DB, func(tx *sqlx.Tx) error {

		isExist := s.Repo.CheckUniqueItemExist()
		if isExist {
			fmt.Println("Items already exist, skipping unique item crawling.")
			return nil
		}
		fmt.Println("Items already exist, skipping unique item crawling.")
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
			"body_armours",
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

		var wg sync.WaitGroup
		itemsCh := make(chan models.Item)
		items := []models.Item{}

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

		err := s.Repo.AddUniqueItems(tx, &items)
		if err != nil {
			return err
		}
		return nil
	})

}

// base items
func (s *ItemService) CrawlingAndAddBaseItemsService() error {
	return dbutils.ExecTx(s.Repo.DB, func(tx *sqlx.Tx) error {

		isExist := s.Repo.CheckBaseItemExist()
		if isExist {
			fmt.Println("base Items already exist, skipping base item crawling.")
			return nil
		}

		var wg sync.WaitGroup
		itemsCh := make(chan models.BaseItem)
		items := []models.BaseItem{}

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
		err := s.Repo.AddBaseItems(tx, &items)
		if err != nil {
			return err
		}
		return nil
	})
}

// get base items
func getBaseItemEquipType(equipType string, itemsCh chan models.BaseItem, wg *sync.WaitGroup) {

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
		go getBaseItemTable(equipType, box, itemsCh, wg)
	})

}

func getBaseItemTable(equipType string, table *goquery.Selection, itemsCh chan models.BaseItem, wg *sync.WaitGroup) {
	defer wg.Done()

	slot := ""
	category := table.Find("h1").First()
	if equipType == "weapon" {
		slot = "Weapon"
		// slot = "One Hand"
		// if strings.Contains(category.Text(), "Two Hand") {
		// 	slot = "Two Hand"
		// }
	}
	if equipType == "armour" {
		if strings.Contains(category.Text(), "Body Armour") {
			slot = "Body Armour"
		}
		if strings.Contains(category.Text(), "Helmets") {
			slot = "Helmet"
		}
		if strings.Contains(category.Text(), "Boots") {
			slot = "Boots"
		}
		if strings.Contains(category.Text(), "Gloves") {
			slot = "Gloves"
		}
		if strings.Contains(category.Text(), "Shields") {
			slot = "Shield"
		}
	}
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
		var baseItem models.BaseItem
		table.Find("tbody tr").Each(func(trIndex int, tr *goquery.Selection) {
			if trIndex%2 == 0 {
				baseItem = models.BaseItem{
					Category:   category.Text(),
					Type:       category.Text(),
					Class:      category.Text(),
					EquipType:  equipType,
					IsTwoHands: strings.Contains(category.Text(), "Two Hand"),
					Slot:       slot,
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
				baseItem.Implicit = append(baseItem.Implicit, strings.TrimSpace(td.Text()))
				itemsCh <- baseItem
			}
		})
	})

}

// base items
func (s *ItemService) CrawlingAndAddItemModsService() error {
	return dbutils.ExecTx(s.Repo.DB, func(tx *sqlx.Tx) error {

		isExist := s.Repo.CheckItemModExist()
		if isExist {
			fmt.Println("Item mods already exist, skipping item mods crawling.")
			return nil
		}

		var wg sync.WaitGroup
		itemsCh := make(chan models.ItemMod)
		items := []models.ItemMod{}

		wg.Add(1)
		go getModHtml(itemsCh, &wg)

		go func() {
			wg.Wait()
			close(itemsCh) // 所有 goroutine 完成後才關閉 Channel
		}()

		for itemCh := range itemsCh {
			items = append(items, itemCh)
		}

		err := s.Repo.AddItemMods(tx, &items)
		if err != nil {
			return err
		}
		return nil
	})
}

func getModHtml(itemsCh chan models.ItemMod, wg *sync.WaitGroup) {

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
		go getItemMod(tr, itemsCh, wg)
	})
}

func getItemMod(tr *goquery.Selection, itemsCh chan models.ItemMod, wg *sync.WaitGroup) {
	defer wg.Done()
	itemMod := models.ItemMod{}
	thList := []string{"Affix", "Name", "Level", "Stat", "Tags"}

	tr.Find("td").Each(func(tdIndex int, td *goquery.Selection) {
		if columnIndex := checkStr(thList, "Affix"); columnIndex == tdIndex {
			itemMod.Affix = td.Text()
		}
		if columnIndex := checkStr(thList, "Name"); columnIndex == tdIndex {
			text := td.Text()
			text = strings.Replace(text, "of", "", 1)
			text = strings.TrimSpace(text)
			itemMod.Name = text
		}
		if columnIndex := checkStr(thList, "Level"); columnIndex == tdIndex {
			itemMod.Level = td.Text()
		}
		if columnIndex := checkStr(thList, "Stat"); columnIndex == tdIndex {
			itemMod.Stat = td.Text()
		}
		if columnIndex := checkStr(thList, "Tags"); columnIndex == tdIndex {
			itemMod.Tags = td.Text()
		}
	})

	itemsCh <- itemMod
}
