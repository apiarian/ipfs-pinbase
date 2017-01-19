IPFS Pin Database
=================

# API

## GET `/nodes`

### Response

#### 200 OK:

```json
{
	"nodes": [
		{
			"hash": "node-hash-a",
			"description": "a helpful description"
		},
		{
			"hash": "node-hash-b",
			"description": "another helpful description"
		},
	],
}
```


## POST `/nodes`

### Data

```json
{
	"api-address": "1.2.3.4:5001",
	"description": "a helpful description"
}
```

### Responses

#### 201 Created:

```json
{
	"hash": "node-hash",
	"description": "a helpful description"
}
```

#### 409 Conflict:

```json
{
	"error-message": "node already exists",
	"error-details": "node 'node-hash' already exists"
}
```

#### 400 Bad Request:

```json
{
	"error-message": "request format invalid",
	"error-details": "api-address and description are required fields"
}
```

#### 403 Forbidden:

```json
{
	"error-message": "cannot add node",
	"error-details": "only administrators may add nodes"
}
```

#### 500 Internal Server Error:

```json
{
	"error-message": "failed to add node",
	"error-details": "could not reach node API: details"
}
```


## GET `/nodes/:node-hash`

### Responses

#### 200 OK:

```json
{
	"hash": "node-hash",
	"description": "a helpful description"
}
```


## POST `/nodes/:node-hash`

### Data

```json
{
	"description": "a new helpful description"
}
```

### Responses

#### 200 OK:

```json
{
	"


## GET `/parties`
