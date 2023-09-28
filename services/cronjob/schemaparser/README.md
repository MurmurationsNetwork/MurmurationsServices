# Schema Parser

The Schema Parser is designed to process and transform schemas from the [Murmurations Library](https://github.com/MurmurationsNetwork/MurmurationsLibrary). Initially, these schemas may use the `$ref` field to reference definitions from other locations.

For example:

```json
{
  ...
  "properties": {
    "lat": {
      "$ref": "../fields/latitude.json"
    },
    ...
  },
  ...
}
```

Upon processing, the Schema Parser not only resolves these references but also accommodates any additional or overidden fields present alongside the `$ref`. After the schemas are parsed and any additions or overrides are incorporated, they are stored in the library service for future use.

## Additional Fields

If there's an additional field accompanying the reference, the Schema Parser will integrate the respective field into the detailed definition.

As an illustration:

```json
{
  ...
  "properties": {
    "lat": {
      "$ref": "../fields/latitude.json",
      "otherfield": "value"
    },
    ...
  },
  ...
}
```

When processed, it transforms to:

```json
{
  ...
  "properties": {
    "lat": {
      // Original detailed fields from latitude.json
      "otherfield": "value"
    },
    ...
  },
  ...
}
```

## Overidden Fields

Sometimes it may be useful to override the value of certain parameters in a field, in particular the `title` and `description` values so that the context of the field they describe are clearer. For example, the standard title/description for the `tags` field is:

```json
{
  "title": "Tags",
  "description": "Keywords relevant to this entity and its activities or attributes, searchable in the Murmurations index"
}
```

But that title/description may not give enough context when used within a specific schema, so the values can be overidden as in [this example](https://github.com/MurmurationsNetwork/MurmurationsLibrary/blob/ae5d972941b053d52016aac5303fb1af4509fe08/schemas/organizations_schema-v1.0.0.json#L26-L30):

```json
"tags": {
  "$ref": "../fields/tags.json",
  "title": "Tags/Type",
  "description": "Keywords that describe the group such as its type, searchable in the Murmurations index"
}
```

After the schema is parsed by Schema Parser, these new values for the `title` and `description` keys will override the default values in the `tags` field.
