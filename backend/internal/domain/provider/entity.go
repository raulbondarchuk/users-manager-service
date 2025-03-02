package provider

type Provider struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Desc string
}
