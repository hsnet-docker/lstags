package tag

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Tag aggregates tag-related information: tag name, image digest etc
type Tag struct {
	name    string
	digest  string
	imageID string
	state   string
	created int64
}

// SortKey returns a sort key
func (tg *Tag) SortKey() string {
	return tg.name
}

// GetName gets tag name
func (tg *Tag) GetName() string {
	return tg.name
}

// GetDigest gets tagged image's digest
func (tg *Tag) GetDigest() string {
	return tg.digest
}

// GetShortDigest gets shorter form of tagged image's digest
func (tg *Tag) GetShortDigest() string {
	return tg.digest[0:40]
}

func calculateImageID(s string) string {
	fields := strings.Split(s, ":")

	if len(fields) < 2 {
		return s
	}

	if len(fields[1]) > 12 {
		return fields[1][0:12]
	}

	return fields[1]
}

// SetImageID sets local Docker image ID
func (tg *Tag) SetImageID(s string) {
	tg.imageID = calculateImageID(s)
}

// GetImageID gets local Docker image ID
func (tg *Tag) GetImageID() string {
	return tg.imageID
}

// SetState sets repo tag state
func (tg *Tag) SetState(state string) {
	tg.state = state
}

// GetState gets repo tag state
func (tg *Tag) GetState() string {
	return tg.state
}

// SetCreated sets image creation timestamp
func (tg *Tag) SetCreated(created int64) {
	tg.created = created
}

// GetCreated gets image creation timestamp
func (tg *Tag) GetCreated() int64 {
	return tg.created
}

// GetCreatedKey gets image creation timestamp in a string form (for a string sort e.g.)
func (tg *Tag) GetCreatedKey() string {
	return strconv.FormatInt(tg.created, 10)
}

// GetCreatedString gets image creation timestamp in a human-readable string form
func (tg *Tag) GetCreatedString() string {
	t := time.Unix(tg.created, 0)

	return t.Format(time.RFC3339)
}

// New creates a new instance of Tag
func New(name, digest string) (*Tag, error) {
	if name == "" {
		return nil, errors.New("Empty tag name not allowed")
	}

	if digest == "" {
		return nil, errors.New("Empty image digest not allowed")
	}

	return &Tag{
			name:   name,
			digest: digest,
		},
		nil
}

func calculateState(sortKey string, registryTags, localTags map[string]*Tag) string {
	r, definedInRegistry := registryTags[sortKey]
	l, definedLocally := localTags[sortKey]

	if definedInRegistry && !definedLocally {
		return "ABSENT"
	}

	if !definedInRegistry && definedLocally {
		return "LOCAL-ONLY"
	}

	if definedInRegistry && definedLocally {
		if r.GetDigest() == l.GetDigest() {
			return "PRESENT"
		}

		return "CHANGED"
	}

	return "UNKNOWN"
}

// Join joins local tags with ones from registry, performs state processing and returns:
// * sorted slice of sort keys
// * joined map of *tag.Tag
func Join(registryTags, localTags map[string]*Tag) ([]string, map[string]*Tag) {
	sortedKeys := make([]string, 0)
	joinedTags := make(map[string]*Tag)

	for sortKey := range registryTags {
		sortedKeys = append(sortedKeys, sortKey)
		joinedTags[sortKey] = registryTags[sortKey]

		ltg, defined := localTags[sortKey]
		if defined {
			joinedTags[sortKey].SetImageID(ltg.GetImageID())
		} else {
			joinedTags[sortKey].SetImageID("n/a")
		}
	}

	for sortKey := range localTags {
		_, defined := registryTags[sortKey]
		if !defined {
			sortedKeys = append(sortedKeys, sortKey)
			joinedTags[sortKey] = localTags[sortKey]
		}
	}

	for sortKey, jtg := range joinedTags {
		jtg.SetState(
			calculateState(
				sortKey,
				registryTags,
				localTags,
			),
		)
	}

	sort.Strings(sortedKeys)

	return sortedKeys, joinedTags
}