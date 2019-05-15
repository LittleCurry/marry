package vm

type TemplateRes struct {
	OriginalTotal int    `bson:"original_total" json:"original_total"`
	Total         int    `bson:"total" json:"total"`
	TemplateName  string `bson:"template_name" json:"template_name"`
	Url           string `bson:"url" json:"url"`
	IsPay         bool   `bson:"is_pay" json:"is_pay"`
}
