package tables

import "time"

type LastUse struct {
	Count    uint64    `gorm:"unique;->;<-;default:0"`
	LastTime time.Time `gorm:"unique;->;<-;autoCreateTime"`
}