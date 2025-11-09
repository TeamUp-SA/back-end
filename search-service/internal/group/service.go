package group

import (
	"context"
	"strings"
)

// Service orchestrates search specific business logic.
type Service interface {
	Search(ctx context.Context, filter SearchFilter) ([]Group, error)
}

type service struct {
	client       Client
	defaultLimit int64
	maxLimit     int64
}

// NewService creates a new search service instance.
func NewService(client Client, defaultLimit, maxLimit int64) Service {
	if defaultLimit <= 0 {
		defaultLimit = 20
	}
	if maxLimit <= 0 {
		maxLimit = 100
	}
	return &service{
		client:       client,
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

func (s *service) Search(ctx context.Context, filter SearchFilter) ([]Group, error) {
	if filter.Limit <= 0 {
		filter.Limit = s.defaultLimit
	}
	if filter.Limit > s.maxLimit {
		filter.Limit = s.maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	includeClosed := true
	if filter.IncludeClosed != nil {
		includeClosed = *filter.IncludeClosed
	}

	groups, err := s.client.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	matching := filterGroups(groups, filter, includeClosed)
	start := clamp(int(filter.Offset), 0, len(matching))
	end := clamp(start+int(filter.Limit), start, len(matching))

	return matching[start:end], nil
}

func filterGroups(groups []Group, filter SearchFilter, includeClosed bool) []Group {
	titleQuery := strings.ToLower(strings.TrimSpace(filter.Title))
	dateQuery := strings.ToLower(strings.TrimSpace(filter.DateQuery))
	tagLookup := buildTagLookup(filter.Tags)

	var result []Group
	for _, g := range groups {
		if !includeClosed && g.Closed {
			continue
		}
		if titleQuery != "" && !strings.Contains(strings.ToLower(g.Title), titleQuery) {
			continue
		}
		if dateQuery != "" && !strings.Contains(strings.ToLower(g.Date), dateQuery) {
			continue
		}
		if len(tagLookup) > 0 && !groupMatchesTags(g.Tags, tagLookup) {
			continue
		}
		result = append(result, g)
	}
	return result
}

func buildTagLookup(tags []GroupTag) map[GroupTag]struct{} {
	if len(tags) == 0 {
		return nil
	}
	lookup := make(map[GroupTag]struct{}, len(tags))
	for _, tag := range tags {
		lookup[tag] = struct{}{}
	}
	return lookup
}

func groupMatchesTags(groupTags []GroupTag, lookup map[GroupTag]struct{}) bool {
	for _, tag := range groupTags {
		if _, ok := lookup[tag]; ok {
			return true
		}
	}
	return false
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
