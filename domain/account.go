package domain

type Account struct {
	Id        string  `json:"id" bson:"_id,omitempty"`
	Name      string  `json:"Name" bson:"Name"`
	Cpf       string  `json:"cpf" bson:"cpf"`
	Ballance  float64 `json:"ballance" bson:"ballance"`
	CreatedAt string  `json:"created_at" bson:"created_at"`
}
