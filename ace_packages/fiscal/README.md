# Fiscal Year Management

Nepal-specific fiscal year management with Bikram Sambat calendar support and automatic invoice numbering.

## Features

- ✅ Nepal fiscal year support (Shrawan 1 to Ashad 32)
- ✅ Bikram Sambat (BS) calendar conversion
- ✅ Automatic invoice/purchase/voucher numbering
- ✅ Fiscal year open/close management
- ✅ Multi-tenant with RLS
- ✅ Atomic number generation (thread-safe)

## Nepal Fiscal Year

Nepal's fiscal year runs from **Shrawan 1** (mid-July) to **Ashad 32** (mid-July next year).

Example: Fiscal Year **2082/83**
- Start: 2082-04-01 BS (Shrawan 1, 2082) = 2025-07-17 AD
- End: 2083-03-32 BS (Ashad 32, 2082) = 2026-07-16 AD

## Usage

### Initialize Fiscal Module

```go
import "github.com/aceextension/fiscal"

func main() {
    fiscal.Init()
}
```

### Create Fiscal Year

```go
// Method 1: From Nepali fiscal year name
fy, err := fiscal.Service.CreateFromNepaliDate(ctx, tenantID, "2082/83")

// Method 2: From AD dates
fy, err := fiscal.Service.Create(ctx, tenantID, "2082/83", startDate, endDate)
```

### Set as Current

```go
err := fiscal.Service.SetAsCurrent(ctx, tenantID, fiscalYearID)
```

### Generate Document Numbers

```go
// Generate invoice number (e.g., "INV-8283-0001")
invoiceNum, err := fiscal.Service.GenerateInvoiceNumber(ctx, fiscalYearID)

// Generate purchase number (e.g., "PUR-8283-0001")
purchaseNum, err := fiscal.Service.GeneratePurchaseNumber(ctx, fiscalYearID)

// Generate voucher number (e.g., "JV-8283-0001")
voucherNum, err := fiscal.Service.GenerateVoucherNumber(ctx, fiscalYearID)
```

### Close/Reopen Fiscal Year

```go
// Close fiscal year (prevents new transactions)
err := fiscal.Service.Close(ctx, fiscalYearID, closedBy)

// Reopen fiscal year
err := fiscal.Service.Reopen(ctx, fiscalYearID)
```

### Nepali Date Utilities

```go
import "github.com/aceextension/fiscal/utils"

// Convert AD to BS
today := time.Now()
todayBS := utils.ADToBS(today)
fmt.Println(todayBS.String()) // "2082-04-15"

// Convert BS to AD
bsDate := utils.NepaliDate{Year: 2082, Month: 4, Day: 1}
adDate := utils.BSToAD(bsDate)

// Get current Nepali date
currentBS := utils.GetCurrentNepaliDate()

// Get fiscal year dates
startBS, endBS, startAD, endAD := utils.GetFiscalYearDates("2082/83")

// Format Nepali date
formatted := utils.FormatNepaliDate(bsDate, "DD MMMM YYYY")
// Output: "1 Shrawan 2082"
```

## Nepali Calendar

### Months (Bikram Sambat)

1. **Baishakh** (April-May) - 30/31 days
2. **Jestha** (May-June) - 31/32 days
3. **Ashad** (June-July) - 31/32 days
4. **Shrawan** (July-August) - 29/30/31/32 days ← **Fiscal Year Start**
5. **Bhadra** (August-September) - 29/30/31 days
6. **Ashwin** (September-October) - 29/30 days
7. **Kartik** (October-November) - 29/30 days
8. **Mangsir** (November-December) - 29/30 days
9. **Poush** (December-January) - 29/30 days
10. **Magh** (January-February) - 29/30 days
11. **Falgun** (February-March) - 29/30 days
12. **Chaitra** (March-April) - 30/31 days

### Supported Years

Calendar data included for BS years **2080-2090** (AD 2023-2033).

For production use beyond 2090 BS, extend the `nepaliMonthDays` map in `utils/nepali_date.go`.

## Document Numbering

Each fiscal year maintains separate counters for:
- **Invoices**: `INV-{year}-{number}` (e.g., `INV-8283-0001`)
- **Purchases**: `PUR-{year}-{number}` (e.g., `PUR-8283-0001`)
- **Vouchers**: `JV-{year}-{number}` (e.g., `JV-8283-0001`)

Numbers are:
- **Atomic**: Thread-safe increment using database
- **Sequential**: No gaps in numbering
- **Fiscal year scoped**: Resets each fiscal year

## Database Schema

```sql
CREATE TABLE fiscal_years (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(20) NOT NULL,              -- "2082/83"
    start_date DATE NOT NULL,               -- AD date
    end_date DATE NOT NULL,                 -- AD date
    start_date_bs VARCHAR(15) NOT NULL,     -- "2082-04-01"
    end_date_bs VARCHAR(15) NOT NULL,       -- "2083-03-32"
    is_current BOOLEAN DEFAULT FALSE,       -- Only one per tenant
    is_closed BOOLEAN DEFAULT FALSE,
    invoice_prefix VARCHAR(20),             -- "INV-8283-"
    purchase_prefix VARCHAR(20),            -- "PUR-8283-"
    voucher_prefix VARCHAR(20),             -- "JV-8283-"
    last_invoice_num INTEGER DEFAULT 0,
    last_purchase_num INTEGER DEFAULT 0,
    last_voucher_num INTEGER DEFAULT 0,
    ...
);
```

## Best Practices

1. **One Current Fiscal Year**: Only one fiscal year should be current per tenant
2. **Close Previous Years**: Close fiscal years after year-end to prevent modifications
3. **Use Atomic Numbering**: Always use the service methods for number generation
4. **Validate Dates**: Ensure fiscal year dates don't overlap
5. **Backup Before Close**: Closing is reversible, but backup first

## Example: Complete Flow

```go
// 1. Create fiscal year for 2082/83
fy, _ := fiscal.Service.CreateFromNepaliDate(ctx, tenantID, "2082/83")

// 2. Set as current
fiscal.Service.SetAsCurrent(ctx, tenantID, fy.ID)

// 3. Generate invoice numbers throughout the year
invoice1, _ := fiscal.Service.GenerateInvoiceNumber(ctx, fy.ID) // "INV-8283-0001"
invoice2, _ := fiscal.Service.GenerateInvoiceNumber(ctx, fy.ID) // "INV-8283-0002"

// 4. At year end, close the fiscal year
fiscal.Service.Close(ctx, fy.ID, adminUserID)

// 5. Create next fiscal year
nextFY, _ := fiscal.Service.CreateFromNepaliDate(ctx, tenantID, "2083/84")
fiscal.Service.SetAsCurrent(ctx, tenantID, nextFY.ID)

// 6. Numbers reset for new year
invoice1, _ := fiscal.Service.GenerateInvoiceNumber(ctx, nextFY.ID) // "INV-8384-0001"
```

## IRD Compliance

For Nepal IRD (Inland Revenue Department) compliance:
- Fiscal year aligns with Nepal government fiscal year
- Invoice numbers are sequential and auditable
- Closed fiscal years prevent backdating
- Bikram Sambat dates for official documents

## Testing

Run the example:
```bash
go run fiscal/examples/main.go
```

Expected output:
- ✅ Nepali date conversion works
- ✅ Fiscal year created with correct dates
- ✅ Invoice numbers generated sequentially
- ✅ Closed fiscal year prevents new numbers
- ✅ Reopened fiscal year allows new numbers
