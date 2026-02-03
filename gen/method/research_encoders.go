package method

// Research RPC encoders for NotebookLM research operations.
//
// These encoders produce argument arrays for the research RPCs:
// - StartFastResearch (Ljjv0c) - Quick web search
// - StartDeepResearch (QA9ei) - Deep research with more sources
// - PollResearchResults (e3bVqc) - Poll for research completion
// - ImportResearchSources (LBwxtb) - Import sources from research results

// EncodeStartFastResearchArgs encodes arguments for StartFastResearch RPC.
// RPC ID: Ljjv0c
//
// Argument format: [["query", 1], null, 1, "project_id"]
// - [0]: Query tuple [query_string, source_type] where source_type 1 = Web
// - [1]: null (reserved)
// - [2]: 1 (source type flag)
// - [3]: Project ID
func EncodeStartFastResearchArgs(projectID, query string) []interface{} {
	return []interface{}{
		[]interface{}{query, 1}, // Query tuple with source type 1 (Web)
		nil,                     // Reserved
		1,                       // Source type flag
		projectID,               // Project ID
	}
}

// EncodeStartDeepResearchArgs encodes arguments for StartDeepResearch RPC.
// RPC ID: QA9ei
//
// Argument format: [null, [1], ["query", 1], 5, "project_id"]
// - [0]: null (reserved)
// - [1]: [1] - Source types array (1 = Web)
// - [2]: Query tuple [query_string, source_type]
// - [3]: 5 - Research depth/thoroughness parameter
// - [4]: Project ID
func EncodeStartDeepResearchArgs(projectID, query string) []interface{} {
	return []interface{}{
		nil,                     // Reserved
		[]interface{}{1},        // Source types [1] = Web
		[]interface{}{query, 1}, // Query tuple
		5,                       // Research depth parameter
		projectID,               // Project ID
	}
}

// EncodePollResearchResultsArgs encodes arguments for PollResearchResults RPC.
// RPC ID: e3bVqc
//
// Argument format: [null, null, "project_id"]
// - [0]: null (reserved)
// - [1]: null (reserved)
// - [2]: Project ID
func EncodePollResearchResultsArgs(projectID string) []interface{} {
	return []interface{}{
		nil,       // Reserved
		nil,       // Reserved
		projectID, // Project ID
	}
}

// EncodeImportResearchSourcesArgs encodes arguments for ImportResearchSources RPC.
// RPC ID: LBwxtb
//
// Argument format: [null, [1], "task_id", "project_id", [sources]]
// - [0]: null (reserved)
// - [1]: [1] - Source types array (1 = Web)
// - [2]: Task ID from research operation
// - [3]: Project ID
// - [4]: Array of source URLs to import
func EncodeImportResearchSourcesArgs(projectID, taskID string, sources []string) []interface{} {
	// Convert []string to []interface{} for JSON marshaling
	sourceArray := make([]interface{}, len(sources))
	for i, s := range sources {
		sourceArray[i] = s
	}

	return []interface{}{
		nil,              // Reserved
		[]interface{}{1}, // Source types [1] = Web
		taskID,           // Task ID
		projectID,        // Project ID
		sourceArray,      // Sources array
	}
}
