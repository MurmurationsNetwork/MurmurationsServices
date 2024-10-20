# Front-end MongoDB

Our front-end MongoDB supports the [Murmurations Tools](https://github.com/MurmurationsNetwork/MurmurationsTools) and [Murmurations Map](https://github.com/MurmurationsNetwork/MurmurationsMap) applications with distinct databases named `mpgData` and `mapData`.

## Database: `mpgData`

This database is comprised of two primary collections: `users` and `profiles`.

### 1. Collection: `users`

Stores user login information.

- _id: Unique identifier for the user.
- cuid: Client Unique Identifier, specific to the user.
- email_hash: Hashed email address for user privacy.
- last_login: Timestamp of the user's last login.
- password: Hashed password for user authentication.
- profiles: Array of profile IDs associated with the user.

Deprecated fields:

- ipfs: InterPlanetary File System hash, associated with user data.
- ipns: InterPlanetary Name System hash, providing a stable address for the user's IPFS data.

### 2. Collection: `profiles`

Contains saved user profile information.

- _id: Unique identifier for the profile.
- cuid: Client Unique Identifier, specific to the profile.
- last_updated: Timestamp of the last update to the profile.
- linked_schemas: Array of schema identifiers that the profile is linked to.
- node_id: Index Service id for the profile.
- profile: JSON string containing detailed profile information.
- title: The title or name associated with the profile.

Deprecated fields:

- ipfs: Array of IPFS hashes related to the profile data.

## Database: `mapData`

This database contains two collections crucial for managing profile data and settings for data retrieval and sorting.

### 1. Collection: `profiles`

Stores profile information in a structured JSON format.

### 2. Collection: `settings`

Contains configuration settings for data sorting and retrieval, particularly for integration with Elasticsearch.

`sort` field: Specifies the criteria used to sort data within the Elasticsearch engine. The sort value stores the most recently processed profile.
