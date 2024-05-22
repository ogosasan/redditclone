package vote

type Vote struct {
	UserID string `json:"user"`
	Value  int    `json:"vote"`
}

func CreateVote(userID string, value int) *Vote {
	return &Vote{
		UserID: userID,
		Value:  value,
	}
}
