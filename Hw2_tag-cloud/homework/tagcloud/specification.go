package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	tags []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	newTagClouds := TagCloud{make([]TagStat, 0)}
	return &newTagClouds
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (tc *TagCloud) AddTag(tag string) {
	tagIndex := -1
	for i := range tc.tags {
		if tc.tags[i].Tag == tag {
			tagIndex = i
			break
		}
	}
	if tagIndex == -1 {
		tc.tags = append(tc.tags, TagStat{tag, 1})
	} else {
		tc.tags[tagIndex].OccurrenceCount++
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (tc *TagCloud) TopN(n int) []TagStat {

	sort.Slice(tc.tags, func(i, j int) bool {
		return tc.tags[i].OccurrenceCount > tc.tags[j].OccurrenceCount
	})

	topNTags := make([]TagStat, 0, n)
	if n > len(tc.tags) {
		topNTags = tc.tags
	} else {
		for i := 0; i < n; i++ {
			topNTags = append(topNTags, tc.tags[i])
		}
	}
	return topNTags
}
