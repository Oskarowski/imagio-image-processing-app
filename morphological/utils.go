package morphological

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	cachedStructureElements    map[string]StructuringElement
	cachedStructureElementsErr error
	once                       sync.Once
)

type StructuringElement struct {
	Data    [][]int `json:"data"`
	OriginX int     `json:"originX"`
	OriginY int     `json:"originY"`
}

func getStructureElements() (map[string]StructuringElement, error) {
	once.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			cachedStructureElementsErr = err
			return
		}

		structureElementsFilePath := filepath.Join(wd, "morphological", "structure_elements.json")

		cachedStructureElements, cachedStructureElementsErr = LoadStructureElementsFromJSON(structureElementsFilePath)
	})

	return cachedStructureElements, cachedStructureElementsErr
}

func GetStructureElement(structureElementName string) (StructuringElement, error) {
	structureElements, err := getStructureElements()
	if err != nil {
		return StructuringElement{}, err
	}

	if structureElement, exists := structureElements[structureElementName]; exists {
		return structureElement, nil
	}

	return StructuringElement{}, fmt.Errorf("structure element %s not found", structureElementName)
}

func LoadStructureElementsFromJSON(filenamePath string) (map[string]StructuringElement, error) {
	bytes, err := os.ReadFile(filenamePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file to load structure elements: %w", err)
	}

	var structureElements map[string]StructuringElement
	if err := json.Unmarshal(bytes, &structureElements); err != nil {
		return nil, fmt.Errorf("could not parse JSON structure elements content: %w", err)
	}

	return structureElements, nil
}

func ReloadStructureElements() error {
	once = sync.Once{}
	_, err := getStructureElements()
	return err
}

func GetAvailableStructureElementsNames() ([]string, error) {
	structuralElements, err := getStructureElements()
	if err != nil {
		return nil, err
	}

	var names []string
	for name := range structuralElements {
		names = append(names, name)
	}

	return names, nil
}
