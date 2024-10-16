package tokenizers

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (tokenizer Mock) GenerateToken(username string, userID int64) (string, error) {
	//TODO implement me
	panic("implement me")
}
