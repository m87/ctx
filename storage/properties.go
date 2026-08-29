package storage


type Properties struct {
	Id string `gorm:"primaryKey"`
	Theme string 
  FirstDay string
	Timezone string
}

