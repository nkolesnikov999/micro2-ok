package part

import (
	"context"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	allParts, err := s.partRepository.ListParts(ctx)
	if err != nil {
		return nil, err
	}

	// Если фильтр пустой или все поля пустые, возвращаем все детали
	if len(filter.Uuids) == 0 && len(filter.Names) == 0 && len(filter.Categories) == 0 && len(filter.ManufacturerCountries) == 0 && len(filter.Tags) == 0 {
		return allParts, nil
	}

	// Создаем set'ы для O(1) проверки (OR внутри одного поля)
	uuidsSet := makeStringSet(filter.Uuids)
	namesSet := makeStringSet(filter.Names)
	categoriesSet := makeCategorySet(filter.Categories)
	countriesSet := makeStringSet(filter.ManufacturerCountries)
	tagsSet := makeStringSet(filter.Tags)

	// AND между разными полями фильтра
	parts := make([]model.Part, 0, len(allParts))
	for _, part := range allParts {
		if uuidsSet != nil {
			if _, ok := uuidsSet[part.Uuid]; !ok {
				continue
			}
		}
		if namesSet != nil {
			if _, ok := namesSet[part.Name]; !ok {
				continue
			}
		}
		if categoriesSet != nil {
			if _, ok := categoriesSet[part.Category]; !ok {
				continue
			}
		}
		if countriesSet != nil {
			country := ""
			if part.Manufacturer != nil {
				country = part.Manufacturer.Country
			}
			if _, ok := countriesSet[country]; !ok {
				continue
			}
		}
		if tagsSet != nil {
			if !hasAnyTag(part.Tags, tagsSet) {
				continue
			}
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func makeStringSet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func makeCategorySet(values []model.Category) map[model.Category]struct{} {
	if len(values) == 0 {
		return nil
	}
	set := make(map[model.Category]struct{}, len(values))
	for _, category := range values {
		set[category] = struct{}{}
	}
	return set
}

func hasAnyTag(partTags []string, wanted map[string]struct{}) bool {
	for _, tag := range partTags {
		if _, ok := wanted[tag]; ok {
			return true
		}
	}
	return false
}
