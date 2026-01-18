package core

import "github.com/m87/nod"

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

func ConvertToNodKV(properties map[string]string) map[string]*nod.KV {
	kvMap := make(map[string]*nod.KV)
	for key, value := range properties {
		kvMap[key] = &nod.KV{Key: key, ValueText: &value}
	}
	return kvMap
}

func ConvertFromNodKV(nodKV map[string]*nod.KV) map[string]string {
	properties := make(map[string]string)
	for key, kv := range nodKV {
		properties[key] = *kv.ValueText
	}
	return properties
}

func ConvertToNodContent(content map[string]string) map[string]*nod.Content {
	nodContent := make(map[string]*nod.Content)
	for key, value := range content {
		nodContent[key] = &nod.Content{Key: key, Value: &value}
	}
	return nodContent
}

func ConvertFromNodContent(nodContent map[string]*nod.Content) map[string]string {
	content := make(map[string]string)
	for key, c := range nodContent {
		content[key] = *c.Value
	}
	return content
}
