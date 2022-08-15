package constant

var MongoIndex = struct {
	Node    string
	Schema  string
	Mapping string
	Profile string
	Update  string
}{
	Node:    "nodes",
	Schema:  "schemas",
	Mapping: "mappings",
	Profile: "profiles",
	Update:  "updates",
}
