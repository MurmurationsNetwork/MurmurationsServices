package constant

var MongoIndex = struct {
	Node    string
	Schema  string
	Mapping string
	Profile string
	Update  string
	Batch   string
}{
	Node:    "nodes",
	Schema:  "schemas",
	Mapping: "mappings",
	Profile: "profiles",
	Update:  "updates",
	Batch:   "batches",
}
