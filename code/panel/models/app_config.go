package models

type AppConfig struct {
	ID      int64  `json:"id"`
	Keyword string `json:"keyword"`
	Value   string `json:"value"`
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (AppConfig) TableName() string {
	return "app_config"
}

// PutUser updates the given user
func SaveAppConfig(a *AppConfig) error {
	return db.Save(a).Error
}

// GetConfigValueByKey returns the value of the given key
func GetConfigValueByKey(keyword string) (string, error) {
	var a AppConfig
	err := db.Where("keyword = ?", keyword).First(&a).Error
	return a.Value, err
}
