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
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/darkphotonKN/community-builds-microservice/items-service/internal/models"
	"github.com/darkphotonKN/community-builds-microservice/items-service/internal/utils/dbutils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo Repository
}

// type Service struct {
// 	pb.UnimplementedItemServiceServer
// 	Repository *ItemRepository
// 	// MqService *mq.MQService
// }

type Repository interface {
	GetItems(slot string) (*[]models.Item, error)
	CreateItem(createItem *CreateItemRequest) error
	CheckBaseItemExist() bool
	CheckUniqueItemExist() bool
	CheckItemModExist() bool
	AddUniqueItems(tx *sqlx.Tx, items *[]models.Item) error
	AddBaseItems(tx *sqlx.Tx, items *[]models.BaseItem) error
	AddItemMods(tx *sqlx.Tx, items *[]models.ItemMod) error
	UpdateItemById(id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error)
}

// mq version
//
//	func NewItemService(mqService *mq.MQService) *FileService {
//		log.Printf("mqService 是否為 nil: %v", mqService == nil)
//		return &ItemService{
//			MqService: mqService,
//		}
//	}
// func NewItemService(repository *ItemRepository) *Service {
// 	return &Service{
// 		Repository: repository,
// 	}
// }

func toPtr[T any](v T) *T {
	return &v
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
func (s *service) GetItemsService(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	// todo: handler
	items, err := s.repo.GetItems(req.Slot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "取得 items 時發生錯誤: %v", err)
	}

	var pbItems []*pb.Item

	for _, item := range *items {
		pbItems = append(pbItems, &pb.Item{
			Id:                   toPtr(item.ID.String()),
			MemberId:             item.MemberID.String(),
			BaseItemId:           item.BaseItemId.String(),
			ImageUrl:             item.ImageUrl,
			Category:             item.Category,
			Class:                item.Class,
			Name:                 item.Name,
			Type:                 item.Type,
			Description:          item.Description,
			UniqueItem:           item.UniqueItem,
			Slot:                 item.Slot,
			RequiredLevel:        toPtr(item.RequiredLevel),
			RequiredStrength:     toPtr(item.RequiredStrength),
			RequiredDexterity:    toPtr(item.RequiredDexterity),
			RequiredIntelligence: toPtr(item.RequiredIntelligence),
			Armour:               toPtr(item.Armour),
			EnergyShield:         toPtr(item.EnergyShield),
			Evasion:              toPtr(item.Evasion),
			Block:                toPtr(item.Block),
			Ward:                 toPtr(item.Ward),
			Damage:               toPtr(item.Damage),
			APS:                  toPtr(item.APS),
			Crit:                 toPtr(item.Crit),
			PDPS:                 toPtr(item.PDPS),
			EDPS:                 toPtr(item.EDPS),
			DPS:                  toPtr(item.DPS),
			Life:                 toPtr(item.Life),
			Mana:                 toPtr(item.Mana),
			Duration:             toPtr(item.Duration),
			Usage:                toPtr(item.Usage),
			Capacity:             toPtr(item.Capacity),
			Additional:           toPtr(item.Additional),
			Stats:                item.Stats,
			Implicit:             item.Implicit,
			CreatedAt:            toPtr(item.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:            toPtr(item.UpdatedAt.Format(time.RFC3339)),
		})
	}

	return &pb.GetItemsResponse{
		Message: "成功取得items",
		Items:   pbItems,
	}, nil
}
func (s *service) CreateItemService(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {

	err := s.repo.CreateItem(&CreateItemRequest{
		Category: req.Category,
		Class:    req.Class,
		Type:     req.Type,
		Name:     req.Name,
		ImageURL: req.ImageURL,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *service) UpdateItemService(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "無效的 ID 格式: %v", err)
	}

	_, err = s.repo.UpdateItemById(id, UpdateItemReq{
		Id:       id,
		Category: req.Category,
		Class:    req.Class,
		Type:     req.Type,
		Name:     req.Name,
		ImageURL: req.ImageURL,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) InitCrawling(db *sqlx.DB) error {

	s.CrawlingAndAddUniqueItemsService(db)
	s.CrawlingAndAddBaseItemsService(db)
	s.CrawlingAndAddItemModsService(db)
	return nil
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
			titleSpanAttr, _ := item.Find("span").Attr("title")
			titleAbbrAttr, _ := item.Find("abbr").Attr("title")
			// fmt.Println("Href:", attr)

			if strings.Contains(titleAbbrAttr, "Required level") || strings.Contains(titleSpanAttr, "Required level") {
				// fmt.Println("包含子字串:", "strength")
				currentItemThs = append(currentItemThs, "RequiredLevel")
			}
			if strings.Contains(titleAbbrAttr, "Required strength") || strings.Contains(titleSpanAttr, "Required strength") {
				// fmt.Println("包含子字串:", "strength")
				currentItemThs = append(currentItemThs, "Strength")
			}

			if strings.Contains(titleAbbrAttr, "Required dexterity") || strings.Contains(titleSpanAttr, "Required dexterity") {
				// fmt.Println("包含子字串:", "dexterity")
				currentItemThs = append(currentItemThs, "Dexterity")
			}

			if strings.Contains(titleAbbrAttr, "Required intelligence") || strings.Contains(titleSpanAttr, "Required intelligence") {
				// fmt.Println("包含子字串:", "intelligence")
				currentItemThs = append(currentItemThs, "Intelligence")
			}

			if strings.Contains(titleAbbrAttr, "Armour") || strings.Contains(titleSpanAttr, "Armour") {
				// fmt.Println("包含子字串:", "Armour")
				currentItemThs = append(currentItemThs, "Armour")
			}

			if strings.Contains(titleAbbrAttr, "Energy shield") || strings.Contains(titleSpanAttr, "Energy shield") {
				// fmt.Println("包含子字串:", "Energy shield")
				currentItemThs = append(currentItemThs, "EnergyShield")
			}

			if strings.Contains(titleAbbrAttr, "Evasion rating") || strings.Contains(titleSpanAttr, "Evasion rating") {
				// fmt.Println("包含子字串:", "Evasion rating")
				currentItemThs = append(currentItemThs, "Evasion")
			}

			if strings.Contains(titleAbbrAttr, "Chance to block") || strings.Contains(titleSpanAttr, "Chance to block") {
				// fmt.Println("包含子字串:", "Armour")
				currentItemThs = append(currentItemThs, "Block")
			}

			if strings.Contains(titleAbbrAttr, "Ward") || strings.Contains(titleSpanAttr, "Ward") {
				currentItemThs = append(currentItemThs, "Ward")
			}
			// weapon
			if strings.Contains(item.Text(), "Damage") {
				currentItemThs = append(currentItemThs, "Damage")
			}
			if strings.Contains(titleAbbrAttr, "Attacks per second") || strings.Contains(titleSpanAttr, "Attacks per second") {
				currentItemThs = append(currentItemThs, "APS")
			}
			if strings.Contains(titleAbbrAttr, "Local weapon critical strike chance") || strings.Contains(titleSpanAttr, "Local weapon critical strike chance") {
				currentItemThs = append(currentItemThs, "Crit")
			}
			if strings.Contains(titleAbbrAttr, "physical damage per second") || strings.Contains(titleSpanAttr, "physical damage per second") {
				currentItemThs = append(currentItemThs, "pDPS")
			}
			if strings.Contains(titleAbbrAttr, "elemental damage") || strings.Contains(titleSpanAttr, "elemental damage") {
				currentItemThs = append(currentItemThs, "eDPS")
			}
			if strings.Contains(titleAbbrAttr, "Damage per second from all damage types") || strings.Contains(titleSpanAttr, "Damage per second from all damage types") {
				currentItemThs = append(currentItemThs, "DPS")
			}
			// flask
			if strings.Contains(titleAbbrAttr, "Life regenerated over the flask duration") || strings.Contains(titleSpanAttr, "Life regenerated over the flask duration") {
				currentItemThs = append(currentItemThs, "Life")
			}

			if strings.Contains(titleAbbrAttr, "Mana regenerated over the flask duration") || strings.Contains(titleSpanAttr, "Mana regenerated over the flask duration") {
				currentItemThs = append(currentItemThs, "Mana")
			}

			if strings.Contains(item.Text(), "Duration") || strings.Contains(titleAbbrAttr, "Flask effect duration") || strings.Contains(titleSpanAttr, "Flask effect duration") {
				currentItemThs = append(currentItemThs, "Duration")
			}

			if strings.Contains(titleAbbrAttr, "Number of charges consumed on use") || strings.Contains(titleSpanAttr, "Number of charges consumed on use") {
				currentItemThs = append(currentItemThs, "Usage")
			}
			if strings.Contains(titleAbbrAttr, "Maximum number of flask charges held") || strings.Contains(titleSpanAttr, "Maximum number of flask charges held") {
				currentItemThs = append(currentItemThs, "Capacity")
			}
			if strings.Contains(item.Text(), "Stats") {
				currentItemThs = append(currentItemThs, "Stats")
			}
			if strings.Contains(item.Text(), "Additionaldrop restrictions") {
				currentItemThs = append(currentItemThs, "Additional")
			}
		})
		fmt.Println("currentItemThs:", currentItemThs)
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
			// fmt.Println("td.Text()", td.Text())
			content, _ := td.Html()
			// fmt.Println("content", content)
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

func (s *service) CrawlingAndAddUniqueItemsService(db *sqlx.DB) error {
	return dbutils.ExecTx(db, func(tx *sqlx.Tx) error {

		isExist := s.repo.CheckUniqueItemExist()
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

		err := s.repo.AddUniqueItems(tx, &items)
		if err != nil {
			return err
		}
		return nil
	})

}

// base items
func (s *service) CrawlingAndAddBaseItemsService(db *sqlx.DB) error {
	return dbutils.ExecTx(db, func(tx *sqlx.Tx) error {

		isExist := s.repo.CheckBaseItemExist()
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
		err := s.repo.AddBaseItems(tx, &items)
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
func (s *service) CrawlingAndAddItemModsService(db *sqlx.DB) error {
	return dbutils.ExecTx(db, func(tx *sqlx.Tx) error {

		isExist := s.repo.CheckItemModExist()
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

		err := s.repo.AddItemMods(tx, &items)
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
