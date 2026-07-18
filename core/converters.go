package core

import (
	"time"

	"github.com/m87/nod"
)

func ConvertToNodTags(tags []string) []*nod.Tag {
	nodTags := make([]*nod.Tag, len(tags))
	for i, tagName := range tags {
		nodTags[i] = &nod.Tag{Name: tagName}
	}
	return nodTags
}

func ConvertFromNodTags(nodTags []*nod.Tag) []string {
	tags := make([]string, len(nodTags))
	for i, nodTag := range nodTags {
		tags[i] = nodTag.Name
	}
	return tags
}

func ConvertToNodKV(properties map[string]string) map[string]*nod.NodeKV {
	kvMap := make(map[string]*nod.NodeKV)
	for key, value := range properties {
		kvMap[key] = &nod.NodeKV{Key: key, ValueText: &value}
	}
	return kvMap
}

func ConvertFromNodKV(nodKV map[string]*nod.NodeKV) map[string]string {
	properties := make(map[string]string)
	for key, kv := range nodKV {
		if kv != nil && kv.ValueText != nil {
			properties[key] = *kv.ValueText
		}
	}
	return properties
}

func ConvertToNodContent(content map[string]string) map[string]*nod.NodeContent {
	nodContent := make(map[string]*nod.NodeContent)
	for key, value := range content {
		nodContent[key] = &nod.NodeContent{Key: key, Value: &value}
	}
	return nodContent
}

func ConvertFromNodContent(nodContent map[string]*nod.NodeContent) map[string]string {
	content := make(map[string]string)
	for key, c := range nodContent {
		if c != nil && c.Value != nil {
			content[key] = *c.Value
		}
	}
	return content
}

func nodString(kv map[string]*nod.NodeKV, key string) string {
	if value := kv[key]; value != nil && value.ValueText != nil {
		return *value.ValueText
	}
	return ""
}

func nodBool(kv map[string]*nod.NodeKV, key string) bool {
	if value := kv[key]; value != nil && value.ValueBool != nil {
		return *value.ValueBool
	}
	return false
}

func nodInt64(kv map[string]*nod.NodeKV, key string) int64 {
	if value := kv[key]; value != nil && value.ValueInt64 != nil {
		return *value.ValueInt64
	}
	return 0
}

func nodTime(kv map[string]*nod.NodeKV, key string) time.Time {
	if value := kv[key]; value != nil && value.ValueTime != nil {
		return *value.ValueTime
	}
	return time.Time{}
}
