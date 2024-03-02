package types

type ListTableRequest struct {
	ExclusiveTableName string `json:"ExclusiveTableName"`
	Limit              int64  `json:"Limit"`
}

type ListTableResponse struct {
	TableNames             []string `json:"TableNames"`
	LastEvaluatedTableName string   `json:"LastEvaluatedTableName"`
}
