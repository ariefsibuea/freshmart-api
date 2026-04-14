# API Documentation — Freshmart API

Base URL: `http://localhost:8080/api/v1`

---

## Response Envelope

All responses follow a consistent JSON envelope.

**Success (single item):**

```json
{
  "status": "success",
  "data": { ... }
}
```

**Success (list):**

```json
{
  "status": "success",
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 22,
    "total_pages": 3
  }
}
```

**Error:**

```json
{
  "status": "error",
  "error": {
    "message": "..."
  }
}
```

---

## Product Object

| Field          | Type      | Nullable | Description                                         |
| -------------- | --------- | -------- | --------------------------------------------------- |
| `id`           | `integer` | No       | Auto-generated product ID                           |
| `name`         | `string`  | No       | Product name                                        |
| `price`        | `float`   | No       | Price per product                                   |
| `product_type` | `string`  | No       | One of: `Sayuran`, `Protein`, `Buah`, `Snack`       |
| `description`  | `string`  | Yes      | Product description (omitted from response if null) |
| `quantity`     | `integer` | No       | Available stock                                     |
| `created_at`   | `string`  | No       | ISO 8601 timestamp                                  |
| `updated_at`   | `string`  | No       | ISO 8601 timestamp                                  |

---

## Endpoints

### POST /products

Add a new product.

**Request Body:**

```json
{
  "name": "Ayam Broiler",
  "price": 42900,
  "product_type": "Protein",
  "description": "Ayam broiler utuh segar",
  "quantity": 30
}
```

| Field          | Type      | Required | Notes                                                   |
| -------------- | --------- | -------- | ------------------------------------------------------- |
| `name`         | `string`  | Yes      | Cannot be empty or whitespace                           |
| `price`        | `float`   | Yes      | Must be greater than 0                                  |
| `product_type` | `string`  | Yes      | Exact PascalCase: `Sayuran`, `Protein`, `Buah`, `Snack` |
| `description`  | `string`  | No       | Omit or set to `null` to store as NULL in database      |
| `quantity`     | `integer` | Yes      | Must be >= 0                                            |

**Response `201 Created`:**

```json
{
  "status": "success",
  "data": {
    "id": 23,
    "name": "Ayam Broiler",
    "price": 42900,
    "product_type": "Protein",
    "description": "Ayam broiler utuh segar",
    "quantity": 30,
    "created_at": "2026-04-14T10:00:00Z",
    "updated_at": "2026-04-14T10:00:00Z"
  }
}
```

**Error responses:**

| Status | Message                                       | Cause                        |
| ------ | --------------------------------------------- | ---------------------------- |
| `400`  | `invalid request body`                        | Malformed JSON               |
| `422`  | `name is required`                            | Empty or missing `name`      |
| `422`  | `price must be greater than 0`                | `price` is 0 or negative     |
| `422`  | `invalid product_type: must be one of [...]`  | Invalid `product_type` value |
| `422`  | `quantity must be greater than or equal to 0` | Negative `quantity`          |

---

### GET /products

List products with optional search, filter, sort, and pagination.

**Query Parameters:**

| Parameter      | Type      | Default      | Description                                           |
| -------------- | --------- | ------------ | ----------------------------------------------------- |
| `name`         | `string`  | —            | Case-insensitive partial match on product name        |
| `product_type` | `string`  | —            | Exact match: `Sayuran`, `Protein`, `Buah`, or `Snack` |
| `sort_by`      | `string`  | `created_at` | Sort field: `created_at`, `price`, or `name`          |
| `order`        | `string`  | `desc`       | Sort direction: `asc` or `desc`                       |
| `page`         | `integer` | `1`          | Page number (>= 1)                                    |
| `page_size`    | `integer` | `10`         | Items per page (>= 1, max `100`)                      |

**Example requests:**

```bash
# Default list
GET /api/v1/products

# Search by name
GET /api/v1/products?name=ayam

# Filter by type
GET /api/v1/products?product_type=Protein

# Sort by price ascending
GET /api/v1/products?sort_by=price&order=asc

# Combined: search + filter + sort + paginate
GET /api/v1/products?name=ayam&product_type=Protein&sort_by=price&order=asc&page=1&page_size=5
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "id": 12,
      "name": "Ayam Broiler",
      "price": 42900,
      "product_type": "Protein",
      "description": "Ayam broiler utuh segar",
      "quantity": 30,
      "created_at": "2026-04-14T10:00:00Z",
      "updated_at": "2026-04-14T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 5,
    "total_items": 6,
    "total_pages": 2
  }
}
```

**Error responses:**

| Status | Message                                                   | Cause                               |
| ------ | --------------------------------------------------------- | ----------------------------------- |
| `400`  | `invalid query parameter 'page'`                          | Non-integer `page` value            |
| `400`  | `invalid query parameter 'page_size'`                     | Non-integer `page_size` value       |
| `400`  | `invalid product_type: must be one of [...]`              | Invalid `product_type` filter value |
| `400`  | `invalid sort_by: must be one of 'price', 'name', 'date'` | Invalid `sort_by` value             |
| `400`  | `invalid order: must be one of 'asc', 'desc'`             | Invalid `order` value               |

---

### GET /products/:id

Get a single product by its ID. Responses are cached in Redis (TTL: 5 minutes by default).

**Path Parameter:**

| Parameter | Type      | Description    |
| --------- | --------- | -------------- |
| `id`      | `integer` | The product ID |

**Example request:**

```bash
GET /api/v1/products/12
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "id": 12,
    "name": "Ayam Broiler",
    "price": 42900,
    "product_type": "Protein",
    "description": "Ayam broiler utuh segar",
    "quantity": 30,
    "created_at": "2026-04-14T10:00:00Z",
    "updated_at": "2026-04-14T10:00:00Z"
  }
}
```

**Error responses:**

| Status | Message                          | Cause                  |
| ------ | -------------------------------- | ---------------------- |
| `400`  | `invalid product id`             | Non-integer `:id`      |
| `404`  | `product with id {id} not found` | Product does not exist |

---

## HTTP Status Codes

| Code  | Meaning                                                        |
| ----- | -------------------------------------------------------------- |
| `200` | OK — successful GET                                            |
| `201` | Created — successful POST                                      |
| `400` | Bad Request — invalid path/query parameter                     |
| `404` | Not Found — product does not exist                             |
| `408` | Request Timeout — request exceeded the configured timeout (3s) |
| `422` | Unprocessable Entity — request body validation failed          |
| `500` | Internal Server Error — unexpected server-side failure         |
