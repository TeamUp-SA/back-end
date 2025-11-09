package graphql

import (
	"time"

	"search-service/internal/group"

	gql "github.com/graphql-go/graphql"
)

// SchemaBuilder wires the GraphQL schema with the underlying services.
type SchemaBuilder struct {
	GroupService group.Service
}

// Build creates the executable GraphQL schema.
func (b SchemaBuilder) Build() (gql.Schema, error) {
	groupTagEnum := newGroupTagEnum()
	groupType := newGroupType(groupTagEnum)

	rootQuery := gql.NewObject(gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{
			"searchGroups": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(groupType))),
				Args: gql.FieldConfigArgument{
					"title": &gql.ArgumentConfig{
						Type:        gql.String,
						Description: "Case-insensitive substring match on the group title",
					},
					"tags": &gql.ArgumentConfig{
						Type:        gql.NewList(groupTagEnum),
						Description: "Filter groups that match any of these tags",
					},
					"date": &gql.ArgumentConfig{
						Type:        gql.String,
						Description: "Case-insensitive match against the date field (supports partial values)",
					},
					"includeClosed": &gql.ArgumentConfig{
						Type:        gql.Boolean,
						Description: "When false, returns only open groups",
					},
					"limit": &gql.ArgumentConfig{
						Type:        gql.Int,
						Description: "Maximum number of groups to return",
					},
					"offset": &gql.ArgumentConfig{
						Type:        gql.Int,
						Description: "Number of groups to skip (for pagination)",
					},
				},
				Resolve: b.searchGroupsResolver,
			},
			"health": &gql.Field{
				Type: gql.NewNonNull(gql.String),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return "ok", nil
				},
			},
		},
	})

	return gql.NewSchema(gql.SchemaConfig{
		Query: rootQuery,
	})
}

func (b SchemaBuilder) searchGroupsResolver(p gql.ResolveParams) (interface{}, error) {
	filter := group.SearchFilter{
		Title:     getStringArg(p, "title"),
		DateQuery: getStringArg(p, "date"),
		Limit:     getInt64Arg(p, "limit"),
		Offset:    getInt64Arg(p, "offset"),
	}

	if tags, ok := p.Args["tags"]; ok {
		filter.Tags = convertTags(tags)
	}

	if includeClosed, exists := p.Args["includeClosed"]; exists {
		if val, ok := includeClosed.(bool); ok {
			filter.IncludeClosed = &val
		}
	}

	results, err := b.GroupService.Search(p.Context, filter)
	if err != nil {
		return nil, err
	}

	return adaptGroups(results), nil
}

// GraphQLGroup flattens upstream group-service payloads into GraphQL-friendly data.
type GraphQLGroup struct {
	ID          string
	Title       string
	Description string
	OwnerID     string
	Members     []string
	Tags        []group.GroupTag
	Closed      bool
	Date        string
	CreatedAt   string
}

func newGroupType(tagEnum *gql.Enum) *gql.Object {
	return gql.NewObject(gql.ObjectConfig{
		Name: "Group",
		Fields: gql.Fields{
			"id": &gql.Field{
				Type: gql.NewNonNull(gql.ID),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).ID, nil
				},
			},
			"title": &gql.Field{
				Type: gql.NewNonNull(gql.String),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).Title, nil
				},
			},
			"description": &gql.Field{
				Type: gql.NewNonNull(gql.String),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).Description, nil
				},
			},
			"ownerID": &gql.Field{
				Type: gql.NewNonNull(gql.ID),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).OwnerID, nil
				},
			},
			"members": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.ID))),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).Members, nil
				},
			},
			"tags": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(tagEnum))),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).Tags, nil
				},
			},
			"closed": &gql.Field{
				Type: gql.NewNonNull(gql.Boolean),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).Closed, nil
				},
			},
			"date": &gql.Field{
				Type: gql.String,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).Date, nil
				},
			},
			"createdAt": &gql.Field{
				Type: gql.String,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return p.Source.(GraphQLGroup).CreatedAt, nil
				},
			},
		},
	})
}

func newGroupTagEnum() *gql.Enum {
	return gql.NewEnum(gql.EnumConfig{
		Name: "GroupTag",
		Values: gql.EnumValueConfigMap{
			"STUDY":           {Value: group.GroupTagStudy},
			"PROJECT":         {Value: group.GroupTagProject},
			"HACKATHON":       {Value: group.GroupTagHackathon},
			"CASECOMPETITION": {Value: group.GroupTagCaseCompetition},
		},
	})
}

func adaptGroups(groups []group.Group) []GraphQLGroup {
	items := make([]GraphQLGroup, 0, len(groups))
	for _, g := range groups {
		createdAt := ""
		if !g.CreatedAt.IsZero() {
			createdAt = g.CreatedAt.Format(time.RFC3339)
		}
		items = append(items, GraphQLGroup{
			ID:          g.GroupID,
			Title:       g.Title,
			Description: g.Description,
			OwnerID:     g.OwnerID,
			Members:     g.Members,
			Tags:        g.Tags,
			Closed:      g.Closed,
			Date:        g.Date,
			CreatedAt:   createdAt,
		})
	}
	return items
}

func convertTags(value interface{}) []group.GroupTag {
	switch t := value.(type) {
	case []interface{}:
		res := make([]group.GroupTag, 0, len(t))
		for _, raw := range t {
			switch v := raw.(type) {
			case group.GroupTag:
				res = append(res, v)
			case string:
				if parsed, err := group.ParseGroupTag(v); err == nil {
					res = append(res, parsed)
				}
			}
		}
		return res
	case []group.GroupTag:
		return t
	}
	return nil
}

func getStringArg(p gql.ResolveParams, key string) string {
	if value, ok := p.Args[key]; ok {
		if s, ok := value.(string); ok {
			return s
		}
	}
	return ""
}

func getInt64Arg(p gql.ResolveParams, key string) int64 {
	if value, ok := p.Args[key]; ok {
		switch v := value.(type) {
		case int:
			return int64(v)
		case int64:
			return v
		case float64:
			return int64(v)
		}
	}
	return 0
}
