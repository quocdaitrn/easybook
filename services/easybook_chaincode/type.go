package easybook_chaincode

// Hotel stores information of a hotel
type Hotel struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	IsActive      bool            `json:"isActive"`
	Rating        float32         `json:"rating"`
	ServiceLevels []*ServiceLevel `json:"serviceLevels"`
}

// ServiceLevel stores information of a service level in a hotel
type ServiceLevel struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	IsUsed           bool         `json:"isUsed"`
	SatisfactionRate float32      `json:"satisfactionRate"`
	RuleAbidingRate  float32      `json:"ruleAbidingRate"`
	HotelID          string       `json:"hotelId"`
	Agreements       []*Agreement `json:"agreements"`
}

// Agreement stores information of a agreement of a level in a hotel
type Agreement struct {
	ID                          string `json:"id"`
	IsApplied                   bool   `json:"isApplied"`
	TotalFeedbacks              uint   `json:"totalFeedbacks"`
	TotalUnfulfilledCommitments uint   `json:"totalUnfulfilledCommitments"`
	IsAppliedPenalty            bool   `json:"isAppliedPenalty"`
	TotalCompensations          uint   `json:"totalCompensations"`
	TotalNoCompensations        uint   `json:"totalNoCompensations"`
	ServiceLevelID              string `json:"serviceLevelId"`
	HotelID                     string `json:"hotelId"`
}
