package group

import (
	"fmt"
	"strings"
	"time"
)

// GroupTag mirrors the tag representation used in the app-service.
type GroupTag int32

const (
	GroupTagStudy GroupTag = iota + 1
	GroupTagProject
	GroupTagHackathon
	GroupTagCaseCompetition
)

var tagToName = map[GroupTag]string{
	GroupTagStudy:           "STUDY",
	GroupTagProject:         "PROJECT",
	GroupTagHackathon:       "HACKATHON",
	GroupTagCaseCompetition: "CASECOMPETITION",
}

var nameToTag = map[string]GroupTag{
	"STUDY":           GroupTagStudy,
	"PROJECT":         GroupTagProject,
	"HACKATHON":       GroupTagHackathon,
	"CASECOMPETITION": GroupTagCaseCompetition,
}

// String yields the canonical GraphQL representation for the tag.
func (t GroupTag) String() string {
	if name, ok := tagToName[t]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN_%d", t)
}

// ParseGroupTag converts a GraphQL enum literal into an internal tag.
func ParseGroupTag(value string) (GroupTag, error) {
	tag, ok := nameToTag[strings.ToUpper(value)]
	if !ok {
		return 0, fmt.Errorf("unsupported group tag %q", value)
	}
	return tag, nil
}

type Group struct {
	GroupID     string     `json:"groupID"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	OwnerID     string     `json:"ownerID"`
	Members     []string   `json:"members"`
	Tags        []GroupTag `json:"tags"`
	Closed      bool       `json:"closed"`
	Date        string     `json:"date"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// SearchFilter captures the supported filter knobs for the search resolver.
type SearchFilter struct {
	Title         string
	Tags          []GroupTag
	DateQuery     string
	IncludeClosed *bool
	Limit         int64
	Offset        int64
}
