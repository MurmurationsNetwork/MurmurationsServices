# Schema Parser

The Schema Parser is designed to process and transform schemas from the [MurmurationsLibrary](https://github.com/MurmurationsNetwork/MurmurationsLibrary). Initially, these schemas may use the `$ref` field to reference definitions from other locations.

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

Upon processing, the Schema Parser not only resolves these references but also accommodates any overwrite fields present alongside the `$ref`. If there's an additional field accompanying the reference, the Schema Parser will integrate the respective field into the detailed definition.

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

After the schemas are parsed and any customizations or overwrites are incorporated, they are stored in the library service for future usage.
