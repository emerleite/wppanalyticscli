package models

// TemplateListResponse represents the response from the message templates API
type TemplateListResponse struct {
	Data   []MessageTemplate `json:"data"`
	Paging *Paging           `json:"paging,omitempty"`
}

// MessageTemplate represents a WhatsApp message template
type MessageTemplate struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	Language          string               `json:"language"`
	Status            string               `json:"status"`
	Category          string               `json:"category"`
	Components        []TemplateComponent  `json:"components,omitempty"`
	QualityScore      *QualityScore        `json:"quality_score,omitempty"`
	PreviousCategory  string               `json:"previous_category,omitempty"`
	RejectedReason    string               `json:"rejected_reason,omitempty"`
}

// TemplateComponent represents a component of a template
type TemplateComponent struct {
	Type       string                `json:"type"`
	Format     string                `json:"format,omitempty"`
	Text       string                `json:"text,omitempty"`
	URL        string                `json:"url,omitempty"`
	Buttons    []TemplateButton      `json:"buttons,omitempty"`
	Parameters []TemplateParameter   `json:"parameters,omitempty"`
	Example    *TemplateExample      `json:"example,omitempty"`
}

// TemplateButton represents a button in a template
type TemplateButton struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	URL         string `json:"url,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Example     []string `json:"example,omitempty"`
}

// TemplateParameter represents a parameter in a template
type TemplateParameter struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Default string `json:"default,omitempty"`
}

// TemplateExample represents examples in a template
type TemplateExample struct {
	HeaderText []string `json:"header_text,omitempty"`
	BodyText   [][]string `json:"body_text,omitempty"`
	HeaderHandle []string `json:"header_handle,omitempty"`
}

// QualityScore represents the quality score of a template
type QualityScore struct {
	Score   string   `json:"score"`
	Date    int64    `json:"date"`
	Reasons []string `json:"reasons,omitempty"`
}