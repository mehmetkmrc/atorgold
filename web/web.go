package web

import (
	"atorgold/database"
	"atorgold/models"
	"atorgold/response"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var year = time.Now().Year()

func LoginPage(c fiber.Ctx) error {
	path := "login"
	return c.Render(path, fiber.Map{
		"Title": "Login",
	},"layouts/back-layout")
}

func NotFoundPage(c fiber.Ctx) error {
    // 404 durum kodunu ayarla
    c.Status(fiber.StatusNotFound)
    
    path := "404"
    return c.Render(path, fiber.Map{
        "Title": "Sayfa Bulunamadı - 404",
        "RequestedURL": c.OriginalURL(), // İsteğin yapıldığı URL'i template'e gönder
    }, "layouts/main")
}

func IndexPage(c fiber.Ctx) error {

	// Yardımcı fonksiyon
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	GetAllDocumentsWithMainDocument := func(ctx context.Context) ([]*models.MainDocument, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT 
				dm.id as main_id, 
				dm.title as main_title, 
				dm.position as main_position, 
				dm.status as main_status, 
				dm.date as main_date,
				ds.id as sub_id,
				ds.main_id as sub_main_id,
				ds.sub_title,
				ds.product_code,
				ds.sub_message, 
				ds.asset as sub_asset, 
				ds.position as sub_position, 
				ds.status as sub_status, 
				ds.date as sub_date,
				dc.id as content_id, 
				dc.sub_id as content_sub_id,
				dc.about_collection,
				dc.jewellery_care, 
				dc.position as content_position, 
				dc.status as content_status, 
				dc.date as content_date
			FROM 
				doc_main dm
			LEFT JOIN 
				doc_sub ds ON dm.id = ds.main_id
			LEFT JOIN 
				doc_content dc ON ds.id = dc.sub_id
			ORDER BY 
				dm.position, ds.position, dc.position;
		`)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}
		defer rows.Close()

		mainDocMap := make(map[uuid.UUID]*models.MainDocument)
		subDocMap := make(map[uuid.UUID]*models.SubDocument)

		for rows.Next() {
			var mainDocument models.MainDocument
			var subDocument models.SubDocument
			var contentDocument models.ContentDocument

			err := rows.Scan(
				&mainDocument.ID, &mainDocument.MainTitle, &mainDocument.Position, &mainDocument.Status, &mainDocument.Date,
				&subDocument.ID, &subDocument.MainID, &subDocument.SubTitle, &subDocument.ProductCode, &subDocument.SubMessage, &subDocument.Asset, &subDocument.Position, &subDocument.Status, &subDocument.Date,
				&contentDocument.ID, &contentDocument.SubID, &contentDocument.ColText, &contentDocument.JewCare, &contentDocument.Position, &contentDocument.Status, &contentDocument.Date,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			if _, ok := mainDocMap[mainDocument.ID]; !ok {
				mainDocument.SubDocuments = []*models.SubDocument{}
				mainDocMap[mainDocument.ID] = &mainDocument
			}

			if subDocument.ID != uuid.Nil {
				if _, ok := subDocMap[subDocument.ID]; !ok {
					subDocument.ContentDocuments = []*models.ContentDocument{}
					subDocMap[subDocument.ID] = &subDocument
					mainDocMap[mainDocument.ID].SubDocuments = append(mainDocMap[mainDocument.ID].SubDocuments, &subDocument)
				}

				if contentDocument.ID != uuid.Nil {
					subDocMap[subDocument.ID].ContentDocuments = append(subDocMap[subDocument.ID].ContentDocuments, &contentDocument)
				}
			}
		}

		if err := rows.Err(); err != nil {
			fmt.Println("Satır işleme hatası: ", err)
			return nil, err
		}

		mainDocuments := make([]*models.MainDocument, 0, len(mainDocMap))
		for _, mainDoc := range mainDocMap {
			mainDocuments = append(mainDocuments, mainDoc)
		}
		return mainDocuments, nil
	}

	encodeBase64 := func(data []byte) string {
		encoded := base64.StdEncoding.EncodeToString(data)
		// Base64 encoded stringi HTML uyumluluğu için temizle
		encoded = strings.ReplaceAll(encoded, "+", "%2B")
		encoded = strings.ReplaceAll(encoded, "=", "%3D")
		return encoded
	}

	const maxCarouselProducts = 6

	allDocuments, err := GetAllDocumentsWithMainDocument(c.Context())
	if err != nil {
		return err
	}

	var carouselProducts []interface{}
	for _, document := range allDocuments {
		for _, subDoc := range document.SubDocuments {
			var assets []map[string]interface{}
			if len(subDoc.Asset) > 0 {
				for i := 0; i < min(2, len(subDoc.Asset)); i++ {
					asset := subDoc.Asset[i]
					encodedAsset := "data:" + http.DetectContentType(asset) + ";base64," + encodeBase64(asset)
					assets = append(assets, map[string]interface{}{
						"url":     template.URL(encodedAsset),
						"main_id": subDoc.MainID,
					})
				}
			}

			product := map[string]interface{}{
				"id":           subDoc.ID,
				"main_id":      subDoc.MainID,
				"title":        subDoc.SubTitle,
				"category":     document.MainTitle,
				"images":       assets,
				"product_code": subDoc.ProductCode,
			}
			if len(carouselProducts) < maxCarouselProducts {
				carouselProducts = append(carouselProducts, product)
			}
		}
	}

	path := "home"
	return c.Render(path, fiber.Map{
		"Title":            "Ator Gold - Anasayfa",
		"CarouselProducts": carouselProducts, // Carousel ürünlerini gönder
	}, "layouts/main")
}

func AboutPage(c fiber.Ctx) error {
	path := "about-us"
	return c.Render(path, fiber.Map{
		"Title": "Hakkımızda",
	},"layouts/main")
}

func ContactsPage(c fiber.Ctx) error {
	path := "contacts"
	return c.Render(path, fiber.Map{
		"Title": "İletişim",
	},"layouts/main")
}

func ProductsListPage(c fiber.Ctx) error {
	categoryFilter := c.Query("category")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "24"))


	GetAllDocumentsByJoin := func(ctx context.Context) ([]*models.MainDocument, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT 
				dm.id as main_id, 
				dm.title as main_title, 
				dm.position as main_position, 
				dm.status as main_status, 
				dm.date as main_date,
				ds.id as sub_id,
				ds.main_id as sub_main_id,
				ds.sub_title,
				ds.product_code,
				ds.sub_message, 
				ds.asset as sub_asset, 
				ds.position as sub_position, 
				ds.status as sub_status, 
				ds.date as sub_date,
				dc.id as content_id, 
				dc.sub_id as content_sub_id,
				dc.about_collection,
				dc.jewellery_care, 
				dc.position as content_position, 
				dc.status as content_status, 
				dc.date as content_date
			FROM 
				doc_main dm
			LEFT JOIN 
				doc_sub ds ON dm.id = ds.main_id
			LEFT JOIN 
				doc_content dc ON ds.id = dc.sub_id
			ORDER BY 
				dm.position, ds.position, dc.position;
		`)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}
		defer rows.Close()


		mainDocMap := make(map[uuid.UUID]*models.MainDocument)
		subDocMap := make(map[uuid.UUID]*models.SubDocument)
		
		for rows.Next() {
			var mainDocument models.MainDocument
			var subDocument models.SubDocument
			var contentDocument models.ContentDocument
		
			err := rows.Scan(
				&mainDocument.ID, &mainDocument.MainTitle, &mainDocument.Position, &mainDocument.Status, &mainDocument.Date,
				&subDocument.ID, &subDocument.MainID, &subDocument.SubTitle, &subDocument.ProductCode, &subDocument.SubMessage, &subDocument.Asset, &subDocument.Position, &subDocument.Status, &subDocument.Date,
				&contentDocument.ID, &contentDocument.SubID, &contentDocument.ColText, &contentDocument.JewCare, &contentDocument.Position, &contentDocument.Status, &contentDocument.Date,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			if _, ok := mainDocMap[mainDocument.ID]; !ok {
				mainDocument.SubDocuments = []*models.SubDocument{}
				mainDocMap[mainDocument.ID] = &mainDocument
			}
		
			if subDocument.ID != uuid.Nil {
				if _, ok := subDocMap[subDocument.ID]; !ok {
					subDocument.ContentDocuments = []*models.ContentDocument{}
					subDocMap[subDocument.ID] = &subDocument
					mainDocMap[mainDocument.ID].SubDocuments = append(mainDocMap[mainDocument.ID].SubDocuments, &subDocument)
				}
		
				if contentDocument.ID != uuid.Nil {
					subDocMap[subDocument.ID].ContentDocuments = append(subDocMap[subDocument.ID].ContentDocuments, &contentDocument)
				}
			}
		}
		
		if err := rows.Err(); err != nil {
			fmt.Println("Satır işleme hatası: ", err)
			return nil, err
		}

		mainDocuments := make([]*models.MainDocument, 0, len(mainDocMap))
		for _, mainDoc := range mainDocMap {
			mainDocuments = append(mainDocuments, mainDoc)
		}
	
		return mainDocuments, nil
	}


	GetAllDocumentsWithMainDocument := func (ctx context.Context) ([]*models.MainDocument, error) {
		return GetAllDocumentsByJoin(ctx)
	}	

	encodeBase64 := func(data []byte) string {
		encoded := base64.StdEncoding.EncodeToString(data)
		// Base64 encoded stringi HTML uyumluluğu için temizle
		encoded = strings.ReplaceAll(encoded, "+", "%2B")
		encoded = strings.ReplaceAll(encoded, "=", "%3D")
		return encoded
	}
	

	 // Fonksiyonu çağırıp, sonuçlarını al
	 allDocuments, err := GetAllDocumentsWithMainDocument(c.Context())
	 if err != nil {
		return response.Error_Response(c, "error while trying to get all documents", err, nil, fiber.StatusBadRequest)
	 }

	var filteredDocuments []interface{}
	var totalDocumentsCount int // Toplam ürün sayısını saklayacak değişken
	for _, document := range allDocuments {
		var subDocs []interface{}
		for _, subDoc := range document.SubDocuments {
			if categoryFilter != "" && subDoc.SubTitle != categoryFilter {
				continue
			}

			var assets []map[string]interface{}
			for _, asset := range subDoc.Asset {
				encodedAsset := "data:" + http.DetectContentType(asset) + ";base64," + encodeBase64(asset)
				//assets = append(assets, template.URL(encodedAsset))
				assets = append(assets, map[string]interface{}{
					"url": template.URL(encodedAsset),
					"main_id": subDoc.MainID,
				})
			}

			subDocument := map[string]interface{}{
				"id":           subDoc.ID,
				"main_id":      subDoc.MainID,
				"sub_title":    subDoc.SubTitle,
				"product_code": subDoc.ProductCode,
				"asset":        assets,
			}
			subDocs = append(subDocs, subDocument)
		}

		if len(subDocs) > 0 {
			documentMap := map[string]interface{}{
				"id":           document.ID,
				"title":        document.MainTitle,
				"SubDocuments": subDocs,
			}
			filteredDocuments = append(filteredDocuments, documentMap)
		}
	}
	
	totalDocumentsCount = len(filteredDocuments) // Toplam ürün sayısını kaydet
    
    // Offset ve limit'e göre filtreleme
	if len(filteredDocuments) > 0 {
		start := offset
		end := offset + limit
		if start >= len(filteredDocuments) {
			filteredDocuments = []interface{}{}
		} else if end > len(filteredDocuments) {
			end = len(filteredDocuments)
		}
		filteredDocuments = filteredDocuments[start:end]
	}

	if c.Query("offset") != "" { // Eğer offset parametresi varsa sadece ürünleri döndür
		response := make([]map[string]interface{}, 0)
		for _, doc := range filteredDocuments {
			if docMap, ok := doc.(map[string]interface{}); ok {
				response = append(response, docMap)
			}
		}
    	remainingItems := totalDocumentsCount - (offset + len(response))
        return c.JSON(fiber.Map{
            "items": response,
            "hasMore": remainingItems > 0,
        })
	}

	path := "product-list"
	return c.Render(path, fiber.Map{
		"Title":            "Ürünler",
		"FilteredDocuments": filteredDocuments,
		"ActiveCategory":   categoryFilter,
		"TotalDocumentsCount": totalDocumentsCount, // Toplam ürün sayısını şablona gönder
	}, "layouts/main")
}

func ProductSinglePage(c fiber.Ctx) error {
	mainIDStr := c.Params("main_id")
	mainID, err := uuid.Parse(mainIDStr)
	if err != nil {
		return err
	}

	GetAllDocumentsWithMainDocument := func(ctx context.Context) ([]*models.MainDocument, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT 
				dm.id as main_id, 
				dm.title as main_title, 
				dm.position as main_position, 
				dm.status as main_status, 
				dm.date as main_date,
				ds.id as sub_id,
				ds.main_id as sub_main_id,
				ds.sub_title,
				ds.product_code,
				ds.sub_message, 
				ds.asset as sub_asset, 
				ds.position as sub_position, 
				ds.status as sub_status, 
				ds.date as sub_date,
				dc.id as content_id, 
				dc.sub_id as content_sub_id,
				dc.about_collection,
				dc.jewellery_care, 
				dc.position as content_position, 
				dc.status as content_status, 
				dc.date as content_date
			FROM 
				doc_main dm
			LEFT JOIN 
				doc_sub ds ON dm.id = ds.main_id
			LEFT JOIN 
				doc_content dc ON ds.id = dc.sub_id
			ORDER BY 
				dm.position, ds.position, dc.position;
		`)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}
		defer rows.Close()

		mainDocMap := make(map[uuid.UUID]*models.MainDocument)
		subDocMap := make(map[uuid.UUID]*models.SubDocument)

		for rows.Next() {
			var mainDocument models.MainDocument
			var subDocument models.SubDocument
			var contentDocument models.ContentDocument

			err := rows.Scan(
				&mainDocument.ID, &mainDocument.MainTitle, &mainDocument.Position, &mainDocument.Status, &mainDocument.Date,
				&subDocument.ID, &subDocument.MainID, &subDocument.SubTitle, &subDocument.ProductCode, &subDocument.SubMessage, &subDocument.Asset, &subDocument.Position, &subDocument.Status, &subDocument.Date,
				&contentDocument.ID, &contentDocument.SubID, &contentDocument.ColText, &contentDocument.JewCare, &contentDocument.Position, &contentDocument.Status, &contentDocument.Date,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			if _, ok := mainDocMap[mainDocument.ID]; !ok {
				mainDocument.SubDocuments = []*models.SubDocument{}
				mainDocMap[mainDocument.ID] = &mainDocument
			}

			if subDocument.ID != uuid.Nil {
				if _, ok := subDocMap[subDocument.ID]; !ok {
					subDocument.ContentDocuments = []*models.ContentDocument{}
					subDocMap[subDocument.ID] = &subDocument
					mainDocMap[mainDocument.ID].SubDocuments = append(mainDocMap[mainDocument.ID].SubDocuments, &subDocument)
				}

				if contentDocument.ID != uuid.Nil {
					subDocMap[subDocument.ID].ContentDocuments = append(subDocMap[subDocument.ID].ContentDocuments, &contentDocument)
				}
			}
		}

		if err := rows.Err(); err != nil {
			fmt.Println("Satır işleme hatası: ", err)
			return nil, err
		}

		mainDocuments := make([]*models.MainDocument, 0, len(mainDocMap))
		for _, mainDoc := range mainDocMap {
			mainDocuments = append(mainDocuments, mainDoc)
		}
		return mainDocuments, nil
	}

	encodeBase64 := func(data []byte) string {
		encoded := base64.StdEncoding.EncodeToString(data)
		// Base64 encoded stringi HTML uyumluluğu için temizle
		encoded = strings.ReplaceAll(encoded, "+", "%2B")
		encoded = strings.ReplaceAll(encoded, "=", "%3D")
		return encoded
	}

	// Fonksiyonu çağırıp, sonuçlarını al
	documents, err := GetAllDocumentsWithMainDocument(c.Context())
	if err != nil {
		return err
	}

	var mainDocs []interface{}
	var title string                               // Title için değişken
	var relatedSubDocuments []interface{} // İlgili ürünler için değişken

	for _, document := range documents {
		mainDoc := map[string]interface{}{
			"id":    document.ID,
			"title": document.MainTitle,
		}
		mainDocs = append(mainDocs, mainDoc)
	}

	var filteredDocuments []interface{}
	//Sadece belirtilen main_id'ye ait belgeleri filtrele
	for _, document := range documents {
		if document.ID == mainID {
			title = strings.ToUpper(document.MainTitle) // Title'ı burada al
			mainDoc := map[string]interface{}{
				"id":       document.ID,
				"title":    document.MainTitle,
				"status":   document.Status,
				"position": document.Position,
				"date":     document.Date,
			}

			var subDocs []interface{}
			for _, subDoc := range document.SubDocuments {
				var assets []template.URL
				for _, asset := range subDoc.Asset {
					encodedAsset := "data:" + http.DetectContentType(asset) + ";base64," + encodeBase64(asset)
					assets = append(assets, template.URL(encodedAsset))
				}

				subDocument := map[string]interface{}{
					"id":           subDoc.ID,
					"main_id":      subDoc.MainID,
					"sub_title":    subDoc.SubTitle,
					"product_code": subDoc.ProductCode,
					"sub_message":  subDoc.SubMessage,
					"asset":        assets,
					"status":       subDoc.Status,
					"date":         subDoc.Date,
				}

				var contentDocs []interface{}
				for _, contentDoc := range subDoc.ContentDocuments {
					contentDocument := map[string]interface{}{
						"id":             contentDoc.ID,
						"sub_id":         contentDoc.SubID,
						"about_collection": contentDoc.ColText,
						"jewellery_care": contentDoc.JewCare,
						"position":       contentDoc.Position,
						"status":         contentDoc.Status,
						"date":           contentDoc.Date,
					}
					contentDocs = append(contentDocs, contentDocument)
				}
				subDocument["ContentDocuments"] = contentDocs
				subDocs = append(subDocs, subDocument)

				// İlgili alt ürünleri bul ve ekle
				for _, relatedDocument := range documents {
					if relatedDocument.ID != mainID {
						for _, relatedSubDoc := range relatedDocument.SubDocuments {
							if relatedSubDoc.SubTitle == subDoc.SubTitle {
								var relatedAssets []template.URL
								for _, asset := range relatedSubDoc.Asset {
									encodedAsset := "data:" + http.DetectContentType(asset) + ";base64," + encodeBase64(asset)
									relatedAssets = append(relatedAssets, template.URL(encodedAsset))
								}
								relatedSubDocument := map[string]interface{}{
									"id":           relatedSubDoc.ID,
									"main_id":      relatedSubDoc.MainID,
									"title":        relatedDocument.MainTitle,
									"sub_title":    relatedSubDoc.SubTitle,
									"product_code": relatedSubDoc.ProductCode,
									"sub_message":  relatedSubDoc.SubMessage,
									"asset":        relatedAssets,
									"status":       relatedSubDoc.Status,
									"date":         relatedSubDoc.Date,
								}
								relatedSubDocuments = append(relatedSubDocuments, relatedSubDocument)
								if len(relatedSubDocuments) >= 5 {
									break // İstenilen sayıya ulaşıldığında döngüden çık
								}
							}
						}
					}
					if len(relatedSubDocuments) >= 5 {
						break
					}
				}
			}
			mainDoc["SubDocuments"] = subDocs
			filteredDocuments = append(filteredDocuments, mainDoc)
			break // Sadece ilgili main_id yi al
		}
	}
	path := "product-single"
	return c.Render(path, fiber.Map{
		"Title":             title,
		"Year":              year,
		"FilteredDocuments": filteredDocuments,
		"AllDocuments":      mainDocs,
		"RelatedSubDocuments": relatedSubDocuments,
	}, "layouts/main")
}

func AddProductPage(c fiber.Ctx) error {
	path := "add-product"
	return c.Render(path, fiber.Map{
		"Title": "Ürün Ekle",
	},"layouts/main")
}

func UploadHandler(c fiber.Ctx) error {
	GetAllDocumentsWithMainDocument := func(ctx context.Context) ([]*models.MainDocument, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT 
				dm.id as main_id, 
				dm.title as main_title, 
				dm.position as main_position, 
				dm.status as main_status, 
				dm.date as main_date,
				ds.id as sub_id,
				ds.main_id as sub_main_id,
				ds.sub_title,
				ds.product_code,
				ds.sub_message, 
				ds.asset as sub_asset, 
				ds.position as sub_position, 
				ds.status as sub_status, 
				ds.date as sub_date,
				dc.id as content_id, 
				dc.sub_id as content_sub_id,
				dc.about_collection,
				dc.jewellery_care, 
				dc.position as content_position, 
				dc.status as content_status, 
				dc.date as content_date
			FROM 
				doc_main dm
			LEFT JOIN 
				doc_sub ds ON dm.id = ds.main_id
			LEFT JOIN 
				doc_content dc ON ds.id = dc.sub_id
			ORDER BY 
				dm.position, ds.position, dc.position;
		`)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}
		defer rows.Close()

		mainDocMap := make(map[uuid.UUID]*models.MainDocument)
		subDocMap := make(map[uuid.UUID]*models.SubDocument)

		for rows.Next() {
			var mainDocument models.MainDocument
			var subDocument models.SubDocument
			var contentDocument models.ContentDocument

			err := rows.Scan(
				&mainDocument.ID, &mainDocument.MainTitle, &mainDocument.Position, &mainDocument.Status, &mainDocument.Date,
				&subDocument.ID, &subDocument.MainID, &subDocument.SubTitle, &subDocument.ProductCode, &subDocument.SubMessage, &subDocument.Asset, &subDocument.Position, &subDocument.Status, &subDocument.Date,
				&contentDocument.ID, &contentDocument.SubID, &contentDocument.ColText, &contentDocument.JewCare, &contentDocument.Position, &contentDocument.Status, &contentDocument.Date,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			if _, ok := mainDocMap[mainDocument.ID]; !ok {
				mainDocument.SubDocuments = []*models.SubDocument{}
				mainDocMap[mainDocument.ID] = &mainDocument
			}

			if subDocument.ID != uuid.Nil {
				if _, ok := subDocMap[subDocument.ID]; !ok {
					subDocument.ContentDocuments = []*models.ContentDocument{}
					subDocMap[subDocument.ID] = &subDocument
					mainDocMap[mainDocument.ID].SubDocuments = append(mainDocMap[mainDocument.ID].SubDocuments, &subDocument)
				}

				if contentDocument.ID != uuid.Nil {
					subDocMap[subDocument.ID].ContentDocuments = append(subDocMap[subDocument.ID].ContentDocuments, &contentDocument)
				}
			}
		}

		if err := rows.Err(); err != nil {
			fmt.Println("Satır işleme hatası: ", err)
			return nil, err
		}

		mainDocuments := make([]*models.MainDocument, 0, len(mainDocMap))
		for _, mainDoc := range mainDocMap {
			mainDocuments = append(mainDocuments, mainDoc)
		}
		return mainDocuments, nil
	}


	// Fonksiyonu çağırıp, sonuçlarını al
	allDocuments, err := GetAllDocumentsWithMainDocument(c.Context())
	if err != nil {
		return err
	}

	var mainDocs []interface{}
	for _, document := range allDocuments {
		mainDoc := map[string]interface{}{
			"id":    document.ID,
			"title": document.MainTitle,
		}
		mainDocs = append(mainDocs, mainDoc)
	}
	return c.Render("upload", fiber.Map{
		"PageTitle":  "Upload Page",
		"Title":      "Welcome to Otovinn App!",
		"Year":       year,
		"AllDocuments": mainDocs,
	}, "layouts/main")
}