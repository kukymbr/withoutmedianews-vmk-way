package newsportal

import (
	"context"
	"reflect"
	"strings"

	"apisrv/pkg/db"
	"github.com/go-playground/validator/v10"
)

const (
	customValidationCategory = "category"
)

func NewValidator(repo db.NewsRepo) *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = validate.RegisterValidationCtx(customValidationCategory, func(ctx context.Context, fl validator.FieldLevel) bool {
		return validateCategoryExists(ctx, fl, repo)
	})

	return validate
}

func validateCategoryExists(ctx context.Context, fl validator.FieldLevel, repo db.NewsRepo) bool {
	id := int(fl.Field().Int())

	cat, err := repo.CategoryByID(ctx, id)

	return err == nil && cat != nil
}
