package entity

import "time"

type PromotionCatalog struct {
	ID                   uint64    `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	FromDate             time.Time `json:"from_date"`
	ToDate               time.Time `json:"to_date"`
	Status               int       `json:"status"`
	ConditionsSerialized string    `json:"conditions_serialized"`
	ActionsSerialized    string    `json:"actions_serialized"`
	StopRulesProcessing  int       `json:"stop_rules_processing"`
	SortOrder            int       `json:"sort_order"`
	SimpleAction         string    `json:"simple_action"`
	DiscountAmount       float64   `json:"discount_amount"`
	FdsID                int       `json:"fds_id"`
	ConditionsSQL        string    `json:"conditions_sql"`
	PrincipalData        string    `json:"principal_data"`
	FdsRule              string    `json:"fds_rule"`
	ConditionsElastic    string    `json:"conditions_elastic"`
}
