package docs

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

func MergeOpenAPISpecs(baseBytes []byte, moduleFiles [][]byte) ([]byte, error) {
	var baseDoc yaml.Node
	if err := yaml.Unmarshal(baseBytes, &baseDoc); err != nil {
		return nil, fmt.Errorf("gagal parse base openapi spec: %w", err)
	}
	if len(baseDoc.Content) == 0 {
		return nil, fmt.Errorf("base openapi spec kosong")
	}

	rootMapping := baseDoc.Content[0]
	if rootMapping.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("root node base openapi harus berupa mapping")
	}

	pathsNode := findOrCreateMappingChild(rootMapping, "paths")

	componentsNode := findOrCreateMappingChild(rootMapping, "components")
	schemasNode := findOrCreateMappingChild(componentsNode, "schemas")

	for idx, modBytes := range moduleFiles {
		if len(bytes.TrimSpace(modBytes)) == 0 {
			continue
		}

		var modDoc yaml.Node
		if err := yaml.Unmarshal(modBytes, &modDoc); err != nil {
			return nil, fmt.Errorf("gagal parse module spec index %d: %w", idx, err)
		}
		if len(modDoc.Content) == 0 {
			continue
		}

		modRoot := modDoc.Content[0]
		if modRoot.Kind != yaml.MappingNode {
			continue
		}

		for i := 0; i < len(modRoot.Content)-1; i += 2 {
			sectionKey := modRoot.Content[i].Value
			sectionVal := modRoot.Content[i+1]

			switch sectionKey {
			case "paths":
				mergeMappingNodes(pathsNode, sectionVal)
			case "components":
				if sectionVal.Kind == yaml.MappingNode {
					for j := 0; j < len(sectionVal.Content)-1; j += 2 {
						compKey := sectionVal.Content[j].Value
						compVal := sectionVal.Content[j+1]
						if compKey == "schemas" {
							mergeMappingNodes(schemasNode, compVal)
						} else {
							targetComp := findOrCreateMappingChild(componentsNode, compKey)
							mergeMappingNodes(targetComp, compVal)
						}
					}
				}
			case "tags":
				tagsNode := findChildNode(rootMapping, "tags")
				if tagsNode != nil && tagsNode.Kind == yaml.SequenceNode {
					mergeTags(tagsNode, sectionVal)
				}
			}
		}
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(4)
	if err := encoder.Encode(&baseDoc); err != nil {
		return nil, fmt.Errorf("gagal encode merged openapi spec: %w", err)
	}
	_ = encoder.Close()

	return buf.Bytes(), nil
}

func findChildNode(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

func findOrCreateMappingChild(mapping *yaml.Node, key string) *yaml.Node {
	if child := findChildNode(mapping, key); child != nil {
		if child.Kind == yaml.MappingNode {
			return child
		}
	}
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}
	valNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
	mapping.Content = append(mapping.Content, keyNode, valNode)
	return valNode
}

func mergeMappingNodes(targetMapping, sourceMapping *yaml.Node) {
	if targetMapping == nil || sourceMapping == nil {
		return
	}
	if targetMapping.Kind != yaml.MappingNode || sourceMapping.Kind != yaml.MappingNode {
		return
	}

	for i := 0; i < len(sourceMapping.Content)-1; i += 2 {
		keyNode := sourceMapping.Content[i]
		valNode := sourceMapping.Content[i+1]

		existingIdx := -1
		for j := 0; j < len(targetMapping.Content)-1; j += 2 {
			if targetMapping.Content[j].Value == keyNode.Value {
				existingIdx = j
				break
			}
		}

		if existingIdx >= 0 {
			if targetMapping.Content[existingIdx+1].Kind == yaml.MappingNode && valNode.Kind == yaml.MappingNode {
				mergeMappingNodes(targetMapping.Content[existingIdx+1], valNode)
			} else {
				targetMapping.Content[existingIdx+1] = valNode
			}
		} else {
			targetMapping.Content = append(targetMapping.Content, keyNode, valNode)
		}
	}
}

func mergeTags(targetSeq, sourceSeq *yaml.Node) {
	if targetSeq == nil || sourceSeq == nil {
		return
	}
	if targetSeq.Kind != yaml.SequenceNode || sourceSeq.Kind != yaml.SequenceNode {
		return
	}

	existingTags := make(map[string]bool)
	for _, item := range targetSeq.Content {
		if item.Kind == yaml.MappingNode {
			for i := 0; i < len(item.Content)-1; i += 2 {
				if item.Content[i].Value == "name" {
					existingTags[item.Content[i+1].Value] = true
				}
			}
		}
	}

	for _, item := range sourceSeq.Content {
		tagName := ""
		if item.Kind == yaml.MappingNode {
			for i := 0; i < len(item.Content)-1; i += 2 {
				if item.Content[i].Value == "name" {
					tagName = item.Content[i+1].Value
					break
				}
			}
		}
		if tagName != "" && !existingTags[tagName] {
			targetSeq.Content = append(targetSeq.Content, item)
			existingTags[tagName] = true
		}
	}
}
