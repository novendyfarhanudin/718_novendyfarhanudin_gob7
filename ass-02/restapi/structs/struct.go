package structs

import "time"

type Order struct {
	ID            int       `json:"order_id" gorm:"primaryKey;autoIncrement:true"`
	Customer_Name string    `json:"customer_name" gorm:"not nill;type:varchar(191)"`
	Ordered_at    time.Time `json:"ordered_at"`
}

type Item struct {
	ID          int    `json:"item_id" gorm:"primaryKey;autoIncrement:true"`
	Item_Code   string `json:"item_code" gorm:"not null;type:varchar(191)"`
	Description string `json:"description" gorm:"not null;type:varchar(191)"`
	Quantity    int    `json:"quantity"`
	Order_Id    int    `json:"Order_Id" gorm:"index"`
	Items       []Item `json:"item" gorm:"foreignKey:Order_Id"`
}
