package common

type Token struct {
	Symbol  string `gorm:"column:symbol;type:varchar(255);"`
	Amount  string `gorm:"column:amount;type:varchar(255);"`
	Decimal uint   `gorm:"column:decimal;type:integer;"`
}

type Operation struct {
	ProjectName string `gorm:"column:project_name;type:varchar(255)"`
	Url         string `gorm:"column:url;type:varchar(2048)"`
	Token
}

type Intent struct {
	Description string `gorm:"column:description;type:varchar(255)"`
}
