package parser

import (
	"encoding/xml"
	"sync"
)

// Event represents a single activity Event with its relationships
type Event struct {
	ID         string
	Conditions []*Event
	Responses  []*Event
	Includes   []*Event
	Excludes   []*Event
}

// Constraint represents the top-level XML structure
type Constraint struct {
	XMLName    xml.Name   `xml:"constraints"`
	Conditions []Relation `xml:"conditions>condition"`
	Responses  []Relation `xml:"responses>response"`
	Includes   []Relation `xml:"includes>include"`
	Excludes   []Relation `xml:"excludes>exclude"`
}

// Relation represents a generic relationship between two Events
type Relation struct {
	SourceID string `xml:"sourceId,attr"`
	TargetID string `xml:"targetId,attr"`
}

// XmlToEvents parses XML data into a map of Event pointers
func XmlToEvents(byteValue []byte) (*sync.Map, error) {
	var constraint Constraint
	if err := xml.Unmarshal(byteValue, &constraint); err != nil {
		return nil, err
	}

	events := sync.Map{}
	var wg sync.WaitGroup

	// Helper function to get or create an event
	getOrCreateEvent := func(id string) *Event {
		if event, exists := events.Load(id); exists {
			return event.(*Event)
		}
		event := &Event{ID: id}
		events.Store(id, event)
		return event
	}

	// Helper function to process relations
	processRelations := func(relations []Relation, relationType string) {
		defer wg.Done()
		for _, relation := range relations {
			source := getOrCreateEvent(relation.SourceID)
			target := getOrCreateEvent(relation.TargetID)

			switch relationType {
			case "Conditions":
				target.Conditions = append(target.Conditions, source)
			case "Responses":
				target.Responses = append(target.Responses, source)
			case "Includes":
				target.Includes = append(target.Includes, source)
			case "Excludes":
				target.Excludes = append(target.Excludes, source)
			}
		}
	}

	// Start goroutines for each type of relation
	wg.Add(4)
	go processRelations(constraint.Conditions, "Conditions")
	go processRelations(constraint.Responses, "Responses")
	go processRelations(constraint.Includes, "Includes")
	go processRelations(constraint.Excludes, "Excludes")

	// Wait for all goroutines to complete
	wg.Wait()

	return &events, nil
}

// EventsToAdjacencyMatrix converts events map to an adjacency matrix
func EventsToAdjacencyMatrix(events *sync.Map) ([][]int, []string) {
	eventIDs := getAllEventIDs(events)
	indexMap := createIndexMap(eventIDs)
	matrix := createAdjacencyMatrix(events, indexMap)
	return matrix, eventIDs
}

// getAllEventIDs returns a slice of all event IDs
func getAllEventIDs(events *sync.Map) []string {
	var eventIDs []string
	events.Range(func(id, event interface{}) bool {
		eventIDs = append(eventIDs, id.(string))
		return true
	})
	return eventIDs
}

// createIndexMap creates a mapping of event IDs to their indices
func createIndexMap(eventIDs []string) map[string]int {
	indexMap := make(map[string]int)
	for index, id := range eventIDs {
		indexMap[id] = index
	}
	return indexMap
}

// createAdjacencyMatrix creates an adjacency matrix from events
func createAdjacencyMatrix(events *sync.Map, indexMap map[string]int) [][]int {
	size := len(indexMap)
	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	events.Range(func(id, event interface{}) bool {
		n := event.(*Event)
		sourceIndex := indexMap[id.(string)]

		for _, relation := range [][]*Event{n.Conditions, n.Responses, n.Includes, n.Excludes} {
			for _, target := range relation {
				targetIndex := indexMap[target.ID]
				matrix[sourceIndex][targetIndex] = 1
			}
		}

		return true
	})

	return matrix
}

func (e *Event) GetProjection() []*Event {
	var relations []*Event
	relations = append(relations, e.Conditions...)
	relations = append(relations, e.Responses...)
	relations = append(relations, e.Includes...)
	relations = append(relations, e.Excludes...)
	for _, relation := range relations {
		relations = append(relations, relation.Includes...)
		relations = append(relations, relation.Excludes...)
	}
	return relations
}

func AllProjections(events *sync.Map) map[string][]*Event {
	relationsMap := make(map[string][]*Event)
	events.Range(func(id, event interface{}) bool {
		e := event.(*Event)
		relationsMap[id.(string)] = e.GetProjection()
		return true
	})
	return relationsMap
}
