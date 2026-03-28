# Accounting Module API Documentation

Base URL: `/api/v1/accounting`

## 1. Chart of Accounts

### Create Account
**POST** `/accounts`

Create a new general ledger account.

**Request Body:**
```json
{
  "code": "1001",
  "name": "Cash on Hand",
  "type": "ASSET", // ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
  "description": "Main cash account",
  "parentId": "uuid-optional"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "code": "1001",
  "name": "Cash on Hand",
  "type": "ASSET"
  // ...
}
```

### List Accounts
**GET** `/accounts`

Retrieve all accounts for the tenant.

**Response:** `200 OK`
```json
[
  { "id": "uuid", "code": "1001", "name": "Cash", ... },
  { "id": "uuid", "code": "2001", "name": "Accounts Payable", ... }
]
```

---

## 2. Journal Entries

### Create Manual Journal
**POST** `/journals`

Create a draft journal entry. Must be balanced (Debits = Credits).

**Request Body:**
```json
{
  "fiscalYearId": "uuid",
  "date": "2026-02-15",
  "description": "Opening Balance",
  "lines": [
    { "accountId": "uuid-1", "debit": 1000, "credit": 0, "description": "Cash" },
    { "accountId": "uuid-2", "debit": 0, "credit": 1000, "description": "Capital" }
  ]
}
```

### Post Journal Entry
**POST** `/journals/{id}/post`

Finalize a journal entry. Status changes from `DRAFT` to `POSTED`. Cannot be modified after posting.

**Response:** `200 OK`
```json
{ "message": "Journal entry posted successfully" }
```

---

## 3. Reports

### General Ledger
**GET** `/reports/general-ledger`

Fetch flattened ledger entries for a specific account.

**Query Params:**
- `accountId`: UUID (Required)
- `startDate`: YYYY-MM-DD (Required)
- `endDate`: YYYY-MM-DD (Required)

**Response:** `200 OK`
```json
[
  {
    "transactionDate": "2026-02-15",
    "description": "Opening Balance",
    "debit": 1000,
    "credit": 0,
    "balance": 1000
  }
]
```
