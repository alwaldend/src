package experiments

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// digest returns an exact SHA256 identifier for artifact provenance.
func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// workflowIdentity canonicalizes JSON after removing explicitly declared variable fields.
func workflowIdentity(data []byte, exclusions []string) (string, error) {
	var value any
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return "", fmt.Errorf("decode workflow identity: %w", err)
	}
	if _, ok := value.(map[string]any); !ok || !json.Valid(data) {
		return "", fmt.Errorf("workflow must be one JSON object")
	}
	for _, pointer := range exclusions {
		if !strings.HasPrefix(pointer, "/") || pointer == "/" {
			return "", fmt.Errorf("identity exclusion must name an existing non-root field")
		}
		parts := strings.Split(pointer[1:], "/")
		for i, part := range parts {
			for offset := 0; offset < len(part); offset++ {
				if part[offset] == '~' {
					if offset+1 >= len(part) || (part[offset+1] != '0' && part[offset+1] != '1') {
						return "", fmt.Errorf("invalid JSON pointer escape")
					}
					offset++
				}
			}
			parts[i] = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		}
		current := value
		for i, part := range parts {
			last := i == len(parts)-1
			switch node := current.(type) {
			case map[string]any:
				child, exists := node[part]
				if !exists {
					return "", fmt.Errorf("identity exclusion field does not exist")
				}
				if last {
					delete(node, part)
				} else {
					current = child
				}
			case []any:
				index, err := strconv.Atoi(part)
				if err != nil || index < 0 || index >= len(node) || strconv.Itoa(index) != part {
					return "", fmt.Errorf("identity exclusion has an invalid array index")
				}
				if last {
					node[index] = nil
				} else {
					current = node[index]
				}
			default:
				return "", fmt.Errorf("identity exclusion traverses a scalar")
			}
		}
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("encode canonical workflow: %w", err)
	}
	return digest(canonical), nil
}
