package manipulations

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed masks.json
var masksJSON []byte

var (
	cachedMasks    map[string][][]int
	cachedMasksErr error
	once           sync.Once
)

func LoadMasksFromEmbeddedJSON() (map[string][][]int, error) {
	var masks map[string][][]int
	if err := json.Unmarshal(masksJSON, &masks); err != nil {
		return nil, fmt.Errorf("could not parse embedded JSON content: %w", err)
	}

	// Basic validation: check that each mask is non-empty and rectangular.
	for maskName, mask := range masks {
		if len(mask) == 0 {
			return nil, fmt.Errorf("mask %s is empty", maskName)
		}

		rowLen := len(mask[0])
		for i, row := range mask {
			if len(row) != rowLen {
				return nil, fmt.Errorf("mask %s is not rectangular: row %d has length %d, expected %d", maskName, i, len(row), rowLen)
			}
		}
	}

	return masks, nil
}

func getMasks() (map[string][][]int, error) {
	once.Do(func() {
		cachedMasks, cachedMasksErr = LoadMasksFromEmbeddedJSON()
	})

	return cachedMasks, cachedMasksErr
}

func ReloadMasks() error {
	once = sync.Once{}
	_, err := getMasks()
	return err
}

// GetMask gets a mask matrix representation by its name from a map of masks.
//
// Args:
// masks: a map of masks, where the key is the mask name and the value is the mask itself
// maskName: the name of the mask to get
//
// Returns:
// the mask matrix if it exists in the map, otherwise an error
func GetMask(maskName string) ([][]int, error) {
	masks, err := getMasks()
	if err != nil {
		return nil, err
	}

	if mask, exists := masks[maskName]; exists {
		return mask, nil
	}

	return nil, fmt.Errorf("mask %s not found", maskName)
}

// GetAvailableEdgeSharpeningMasksNames returns a list of all available edge sharpening mask names as strings.
//
// Returns:
// a list of strings representing the names of all available edge sharpening masks
// an error if there was an issue reading the masks
func GetAvailableEdgeSharpeningMasksNames() ([]string, error) {
	masks, err := getMasks()
	if err != nil {
		return nil, err
	}

	var maskNames []string
	for maskName := range masks {
		maskNames = append(maskNames, maskName)
	}

	return maskNames, nil
}
