package common

import "time"

type Property struct {
	Id           			uint      `json:"id"`
	Bedrooms     			int       `json:"bedrooms"`
	Size         			float64   `json:"size"`
	Garage       			int       `json:"garage"`
	Bathroom     			int       `json:"bathroom"`
	Floor        			string    `json:"floor"`
	Neighborhood 			string    `json:"neighborhood"`
	City         			string    `json:"city"`
	State        			string    `json:"state"`
	Complement   			string    `json:"complement"`
	UsageType    			string    `json:"usage_type"`
	SizeUnit     			string    `json:"size_unit"`
	Latitude     			float64   `json:"latitude"`
	Longitude    			float64   `json:"longitude"`
	StreetNumber 			int       `json:"street_number"`
	StreetName   			string    `json:"street_name"`
	PostalCode   			string    `json:"postal_code"`
	CreatedAt    			time.Time `json:"created_at"`
	UpdatedAt    			time.Time `json:"updated_at"`
	Iptu				 			float64   `json:"iptu"`
	CondominiumValue	float64   `json:"condominium"`
}

type Round struct {
	Id             uint      `json:"id"`
	Discount       float64   `json:"discount"`
	IncrementValue float64   `json:"increment_value"`
	MinPrice       float64   `json:"min_price"`
	RoundNumber    int       `json:"round_number"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	AuctionId      uint      `json:"auction_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Auction struct {
	Id                  uint      `json:"id"`
	Title               string    `json:"title"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ExternalID          string    `json:"external_id"`
	ExternalUrl         string    `json:"external_url"`
	Origin              string    `json:"origin"`
	AuctioneerComission float64   `json:"auctioneer_comission"`
	AuctioneerViews     int       `json:"auctioneer_views"`
	PriceSold           float64   `json:"price_sold"`
	QualifiedUsers      int       `json:"qualified_users"`
	InternalStatus      string    `json:"internal_status"`
	Description         string    `json:"description"`
	ViewsCount          int       `json:"views_count"`
	BidsCount           int       `json:"bids_count"`
	PropertyID          uint      `json:"property_id"`
	AppraisalValue      float64   `json:"appraisal_value"`
	AuctionnerName      string    `json:"auctionner_name"`
	Ocupation           string    `json:"ocupation"`
	ProccessNumber      string    `json:"proccess_number"`
	CurrentMinBid       float64   `json:"current_min_bid"`
}

type Announcement struct {
	Id          uint      `json:"id"`
	Description string    `json:"description"`
	ExternalUrl string    `json:"external_url"`
	Origin      string    `json:"origin"`
	SalePrice   float64   `json:"sale_price"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExternalID  string    `json:"external_id"`
	PropertyID  uint      `json:"property_id"`
}

type Attachment struct {
	Url  string
	Name string
}

type DataToInsertAuction struct {
	Auction     *Auction
	Property    *Property
	Rounds      *[]Round
	Images      []Attachment
	Attachments []Attachment
}
