package db

import "github.com/jinzhu/copier"

// FromDomainGeneric copies fields from the domain entity to the model.
// T — type of the domain entity, M — type of the ORM model.
// The function returns the model and an error if the copy failed.
func FromDomainGeneric[T any, M any](domain T) (M, error) {
	var model M
	err := copier.Copy(&model, &domain)
	return model, err
}
