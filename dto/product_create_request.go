package dto



type MainDocumentCreateRequest struct {
	MainTitle string `json:"main_title"`
}

type SubDocumentCreateRequest struct {
	MainID      string `json:"main_id"`
	SubTitle    string   `json:"sub_title"`
	ProductCode string   `json:"product_code"`
	SubMessage  string   `json:"sub_message"` //Ürün Özellikleri
	Asset       [][]byte `json:"asset"`
}

type ContentDocumentCreateRequest struct {
	SubID   string `json:"sub_id"`
	ColText string `json:"about_collection"`
	JewCare string `json:"jewellery_care"`
}