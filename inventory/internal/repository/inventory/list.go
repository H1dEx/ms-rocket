package inventory

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository/converter"
	repoModel "github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f := converter.PartFilterToRepoModel(filter)
	filterSet := BuildFilterSets(f)
	parts := filterParts(r.parts, filterSet)
	result := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		result = append(result, converter.PartToModel(part))
	}
	return result, nil
}

type FilterSets struct {
	Uuids                 map[string]struct{}
	Names                 map[string]struct{}
	Categories            map[repoModel.Category]struct{}
	ManufacturerCountries map[string]struct{}
	Tags                  map[string]struct{}
}

func BuildFilterSets(f *repoModel.PartFilter) FilterSets {
	if f == nil {
		return FilterSets{}
	}

	return FilterSets{
		Uuids:                 filterToMap(f.Uuids),
		Names:                 filterToMap(f.Names),
		Categories:            filterToMap(f.Categories),
		ManufacturerCountries: filterToMap(f.ManufacturerCountries),
		Tags:                  filterToMap(f.Tags),
	}
}

func filterToMap[T comparable](f []T) map[T]struct{} {
	if f == nil {
		return nil
	}
	acc := make(map[T]struct{}, len(f))
	for _, v := range f {
		acc[v] = struct{}{}
	}
	return acc
}

func hasAnyTag(tags []string, f map[string]struct{}) bool {
	for _, t := range tags {
		if _, ok := f[t]; ok {
			return true
		}
	}
	return false
}

func filterParts(parts map[string]repoModel.Part, f FilterSets) []repoModel.Part {
	result := []repoModel.Part{}
	for _, p := range parts {
		if f.Uuids != nil && len(f.Uuids) > 0 {
			if _, ok := f.Uuids[p.Uuid]; !ok {
				continue
			}
		}
		if f.Names != nil && len(f.Names) > 0 {
			if _, ok := f.Names[p.Name]; !ok {
				continue
			}
		}
		if f.Categories != nil && len(f.Categories) > 0 {
			if _, ok := f.Categories[p.Category]; !ok {
				continue
			}
		}
		if f.ManufacturerCountries != nil && len(f.ManufacturerCountries) > 0 {
			if _, ok := f.ManufacturerCountries[p.Uuid]; !ok {
				continue
			}
		}
		if f.Tags != nil && len(f.Tags) > 0 {
			if !hasAnyTag(p.Tags, f.Tags) {
				continue
			}
		}
		result = append(result, p)
	}
	return result
}
