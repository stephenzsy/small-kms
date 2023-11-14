package models

import "github.com/golang-jwt/jwt/v5"

type (
	applicationByAppIdComposed struct {
		Ref
		ApplicationByAppIdFields
	}
)

type NumericDate = jwt.NumericDate
