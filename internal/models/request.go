package models

type SearchPersonRequest struct {
	Query string `xml:"Query"`
}
type Body struct {

	SearchPerson   *SearchPersonRequest   `xml:"SearchPerson,omitempty"`
}