package product

import (
	"atorgold/database"
	"atorgold/dto"
	"atorgold/models"
	"atorgold/response"
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DocumentRepository struct {
	dbPool *pgxpool.Pool
}

func CreateMainDocument(c fiber.Ctx) error {
	reqBody := new(dto.MainDocumentCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	MainDocumentCreateRequestToModel := func (req *dto.MainDocumentCreateRequest)(*models.MainDocument, error) {
		mainDocument := new(models.MainDocument)
		mainID, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		mainDocument = &models.MainDocument{
			ID: mainID,
			MainTitle: req.MainTitle,
			Date: time.Now(),
		}
		return mainDocument, nil
	}


	documentModel, err := MainDocumentCreateRequestToModel(reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert document create request to model", err, nil, fiber.StatusBadRequest)
	}

	Insert := func (ctx context.Context, q *DocumentRepository, documentModel *models.MainDocument) (*models.MainDocument, error) {
		query := `INSERT INTO doc_main (id, title, status,  date) VALUES ($1, $2, $3, $4) RETURNING id, title, status, date`
		queryRow := q.dbPool.QueryRow(ctx, query, documentModel.ID, documentModel.MainTitle, documentModel.Status, documentModel.Date)
		err := queryRow.Scan(&documentModel.ID, &documentModel.MainTitle, &documentModel.Status, &documentModel.Date)
		if err != nil {
			return nil, err
		}
		return documentModel, nil
	}

	AddMainDocument := func (ctx context.Context, document *models.MainDocument)(*models.MainDocument, error) {
		repo := &DocumentRepository{dbPool: database.DBPool} // `db` değişkeni daha önce tanımlı olmalı
		return Insert(ctx, repo, document)
	}

	document, err := AddMainDocument(c.Context(), documentModel)
	if err != nil {
		return response.Error_Response(c, "error while trying to create document", err, nil, fiber.StatusBadRequest)
	}

	zap.S().Info("Document Created Successfully! Document:", document)
	return response.Success_Response(c, documentModel.ID, "document created successfully",  fiber.StatusOK)
}

func CreateSubDocument(c fiber.Ctx) error {
	reqBody := new(dto.SubDocumentCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	SubDocumentCreateRequestToModel := func (req *dto.SubDocumentCreateRequest)(*models.SubDocument, error) {
		subDocument := new(models.SubDocument)
		ID, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		mainID, err := uuid.Parse(req.MainID)
		if err != nil{
			return nil, err
		}
		subDocument = &models.SubDocument{
			ID: ID,
			MainID: mainID,
			SubTitle: req.SubTitle,
			ProductCode: req.ProductCode,
			SubMessage: req.SubMessage,
			Asset: req.Asset,
			Date: time.Now(),
		}
		return subDocument, nil
	}

	documentModel, err := SubDocumentCreateRequestToModel(reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert document create request to model", err, nil, fiber.StatusBadRequest)
	}

	Insert := func (ctx context.Context, q *DocumentRepository, documentModel *models.SubDocument) (*models.SubDocument, error) {
		query := `INSERT INTO doc_sub (id, main_id, sub_title, product_code, sub_message, asset, status, date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, main_id, sub_title, product_code, sub_message, asset, status, date`
		queryRow := q.dbPool.QueryRow(ctx, query, documentModel.ID, documentModel.MainID, documentModel.SubTitle, documentModel.ProductCode, documentModel.SubMessage, documentModel.Asset, documentModel.Status, documentModel.Date)
		err := queryRow.Scan(&documentModel.ID, &documentModel.MainID, &documentModel.SubTitle, &documentModel.ProductCode, &documentModel.SubMessage, &documentModel.Asset, &documentModel.Status, &documentModel.Date)
		if err != nil {
			return nil, err
		}
		return documentModel, nil
	}

	AddSubDocument := func (ctx context.Context, document *models.SubDocument)(*models.SubDocument, error) {
		repo := &DocumentRepository{dbPool: database.DBPool} // `db` değişkeni daha önce tanımlı olmalı
		return Insert(ctx, repo, document)
	}

	document, err := AddSubDocument(c.Context(), documentModel)
	if err != nil {
		return response.Error_Response(c, "error while trying to create document", err, nil, fiber.StatusBadRequest)
	}

	zap.S().Info("Document Created Successfully! Document:", document)
	return response.Success_Response(c, documentModel.ID, "document created successfully", fiber.StatusOK)
}

func CreateContentDocument(c fiber.Ctx) error {
	reqBody := new(dto.ContentDocumentCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	ContentDocumentCreateRequestToModel := func (req *dto.ContentDocumentCreateRequest)(*models.ContentDocument, error){
		document := new(models.ContentDocument)
		ID, err := uuid.NewV7()
		if err != nil{
			return nil, err
		}
		subID, err := uuid.Parse(req.SubID)
		if err != nil {
			return nil, err
		}
		document = &models.ContentDocument{
			ID: ID,
			SubID: subID,
			ColText: req.ColText,
			JewCare: req.JewCare,
			Date: time.Now(),
		}
		return document, nil
	}


	documentModel, err := ContentDocumentCreateRequestToModel(reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert document create request to model", err, nil, fiber.StatusBadRequest)
	}

	Insert := func (ctx context.Context, q *DocumentRepository, documentModel *models.ContentDocument) (*models.ContentDocument, error) {
		query := `INSERT INTO doc_content (id, sub_id, about_collection, jewellery_care, status, date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, sub_id, about_collection, jewellery_care, status, date`
		queryRow := q.dbPool.QueryRow(ctx, query, documentModel.ID, documentModel.SubID, documentModel.ColText, documentModel.JewCare, documentModel.Status, documentModel.Date)
		err := queryRow.Scan(&documentModel.ID, &documentModel.SubID, &documentModel.ColText, &documentModel.JewCare, &documentModel.Status, &documentModel.Date)
		if err != nil {
			return nil, err
		}
		return documentModel, nil
	}

	AddContentDocument := func (ctx context.Context, document *models.ContentDocument)(*models.ContentDocument, error) {
		repo := &DocumentRepository{dbPool: database.DBPool} // `db` değişkeni daha önce tanımlı olmalı
		return Insert(ctx, repo, document)
	}

	document, err := AddContentDocument(c.Context(), documentModel)
	if err != nil {
		return response.Error_Response(c, "error while trying to create document", err, nil, fiber.StatusBadRequest)
	}

	zap.S().Info("Document Created Successfully! Document:", document)
	return response.Success_Response(c, nil, "document created successfully", fiber.StatusOK)
}

func GetAllDocuments(c fiber.Ctx) error {

	GetAllMain := func (ctx context.Context, q *DocumentRepository) ([]*models.MainDocument, error) {
		query := `SELECT * FROM doc_main`
		queryRows, err := q.dbPool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
	
		var documents []*models.MainDocument
		for queryRows.Next() {
			document := new(models.MainDocument)
			err = queryRows.Scan(&document.ID, &document.MainTitle, &document.Position, &document.Status, &document.Date)
			if err != nil {
				return nil, err
			}
			documents = append(documents, document)
		}
		return documents, nil
	}

	GetAllSub := func (ctx context.Context, q *DocumentRepository) ([]*models.SubDocument, error) {
		query := `SELECT * FROM doc_sub`
		queryRows, err := q.dbPool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
	
		var documents []*models.SubDocument
		for queryRows.Next() {
			document := new(models.SubDocument)
			err = queryRows.Scan(&document.ID, &document.MainID, &document.ProductCode, &document.SubMessage, &document.Asset, &document.Position, &document.Status, &document.Date)
			if err != nil {
				return nil, err
			}
			documents = append(documents, document)
		}
		return documents, nil
	}
	
	GetAllContent := func (ctx context.Context, q *DocumentRepository) ([]*models.ContentDocument, error) {
		query := `SELECT * FROM doc_content`
		queryRows, err := q.dbPool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
	
		var documents []*models.ContentDocument
		for queryRows.Next() {
			document := new(models.ContentDocument)
			err = queryRows.Scan(&document.ID, &document.SubID, &document.ColText, &document.JewCare, &document.Position, &document.Status, &document.Date)
			if err != nil {
				return nil, err
			}
			documents = append(documents, document)
		}
		return documents, nil
	}


	GetAllDocuments := func (ctx context.Context)([]*models.MainDocument, error) {
		mainDocuments, err := GetAllMain(ctx, &DocumentRepository{})
		if err != nil {
			return nil, err
		}
	
		subDocuments, err := GetAllSub(ctx, &DocumentRepository{})
		if err != nil {
			return nil, err
		}
		contentDocuments, err := GetAllContent(ctx,  &DocumentRepository{})
		if err != nil{
			return nil, err
		}
		subDocMap := make(map[uuid.UUID]*models.SubDocument)
		for _, subDocument := range subDocuments {
			subDocMap[subDocument.ID] = subDocument
			subDocument.ContentDocuments = []*models.ContentDocument{}
		}
	
		for _, contentDocument := range contentDocuments {
			if subDocument, ok := subDocMap[contentDocument.SubID]; ok {
				subDocument.ContentDocuments = append(subDocument.ContentDocuments, contentDocument)
			}
		}
	
		mainDocMap := make(map[uuid.UUID]*models.MainDocument)
		for _, mainDocument := range mainDocuments {
			mainDocMap[mainDocument.ID] = mainDocument
			mainDocument.SubDocuments = []*models.SubDocument{}
		}
	
		for _, subDocument := range subDocuments {
			if mainDocument, ok := mainDocMap[subDocument.MainID]; ok {
				mainDocument.SubDocuments = append(mainDocument.SubDocuments, subDocument)
			}
		}
	
		return mainDocuments, nil
	
	}
	
	documents, err := GetAllDocuments(c.Context())
	if err != nil {
		return response.Error_Response(c, "error while trying to get all documents", err, nil, fiber.StatusBadRequest)
	}

	return response.Success_Response(c, documents, "documents fetched successfully", fiber.StatusOK)
}

func GetAllDocumentsByJoin(c fiber.Ctx) error {

	GetAllDocumentsByJoin:= func (ctx context.Context, q *DocumentRepository) ([]*models.MainDocument, error) {
		rows, err := q.dbPool.Query(ctx, `
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
				return nil, err
			}
	
			if mainDocMap[mainDocument.ID] == nil {
				mainDocument.SubDocuments = []*models.SubDocument{}
				mainDocMap[mainDocument.ID] = &mainDocument
			}
	
			if subDocument.ID != uuid.Nil {
				if subDocMap[subDocument.ID] == nil {
					subDocument.ContentDocuments = []*models.ContentDocument{}
					subDocMap[subDocument.ID] = &subDocument
					mainDocMap[mainDocument.ID].SubDocuments = append(mainDocMap[mainDocument.ID].SubDocuments, &subDocument)
				}
	
				if contentDocument.ID != uuid.Nil {
					subDocMap[subDocument.ID].ContentDocuments = append(subDocMap[subDocument.ID].ContentDocuments, &contentDocument)
				}
			}
		}
	
		mainDocuments := make([]*models.MainDocument, 0, len(mainDocMap))
		for _, mainDoc := range mainDocMap {
			mainDocuments = append(mainDocuments, mainDoc)
		}
	
		return mainDocuments, nil
	}

	GetAllDocumentsWithMainDocument := func (ctx context.Context) ([]*models.MainDocument, error) {
		return GetAllDocumentsByJoin(ctx, &DocumentRepository{})
	}

	documents, err := GetAllDocumentsWithMainDocument(c.Context())
	if err != nil {
		return response.Error_Response(c, "error while trying to get all documents", err, nil, fiber.StatusBadRequest)
	}

	return response.Success_Response(c, documents, "documents fetched successfully", fiber.StatusOK)
}


