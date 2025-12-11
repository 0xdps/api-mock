# Backend Test Coverage Documentation

This document provides a comprehensive overview of all test cases in the API-Mockly backend.

**Last Updated:** December 11, 2025  
**Test Framework:** Go's built-in `testing` package  
**Total Coverage:** ~86%

## Test Statistics

| Component | Test File | Test Count | Coverage | Focus Areas |
|-----------|-----------|------------|----------|-------------|
| Cache (Local) | `cache_test.go` | 35 | ~87% | Cache operations, concurrency, edge cases |
| Cache (Redis Mock) | `cache_redis_test.go` | 12 | ~85% | Mock Redis integration, Redis simulation |
| Filters | `filters_test.go` | 31 | ~92% | Query parameter filtering, operators |
| Handlers | `handlers_test.go` | 30+ | ~75% | HTTP request handling, CRUD operations |
| Utility Handlers | `utility_test.go` | 9 | ~90% | Testing utilities (echo, delay, status, flaky, chaos) |
| Schema Validation | `schema_validation_test.go` | 10 | ~85% | JSON schema validation, constraints |
| Validation | `validation_test.go` | 6 | ~80% | Input validation, type checking |
| CORS Middleware | `cors_test.go` | 10 | ~95% | CORS configuration, preflight |
| **Total** | **8 files** | **~143 tests** | **~86%** | **Full stack coverage** |

---

## Quick Start

### Run All Tests
```bash
cd backend
go test ./internal/...
```

### Run with Coverage
```bash
go test ./internal/... -cover
```

### Run with Verbose Output
```bash
go test ./internal/... -v
```

### Run Specific Test
```bash
go test ./internal/cache -run TestCache_ConcurrentWrites
```

### Run with Race Detection
```bash
go test ./internal/... -race
```

---

## 1. Cache Tests (`internal/cache/cache_test.go`)

### 1.1 Cache Mode Tests
Tests different caching strategies and their behavior.

- ✅ `TestNewCache_LocalMode` - Verifies local in-memory cache initialization
- ✅ `TestNewCache_OffMode` - Tests cache creation in off mode
- ✅ `TestNewCache_ModeFallback` - Tests fallback to local mode when Redis unavailable

### 1.2 Warmup Tests
Tests cache pre-population during startup.

- ✅ `TestCache_Warmup_LocalMode` - Verifies data generation for all resources
- ✅ `TestCache_Warmup_OffMode` - Ensures warmup is skipped in off mode

### 1.3 GET Operations Tests
Tests data retrieval from cache.

- ✅ `TestCache_Get_LocalMode_Hit` - Cache hit scenario
- ✅ `TestCache_Get_LocalMode_Miss` - Cache miss scenario
- ✅ `TestCache_Get_OffMode_GeneratesOnDemand` - On-demand generation in off mode
- ✅ `TestCache_GetByID_LocalMode_Found` - Retrieval by ID (found)
- ✅ `TestCache_GetByID_LocalMode_NotFound` - Retrieval by ID (not found)

### 1.4 CRUD Operations Tests
Tests create, update, and delete operations.

- ✅ `TestCache_AddItem_LocalMode` - Adding new items to cache
- ✅ `TestCache_AddItem_MaxConstraint` - Max items per resource enforcement
- ✅ `TestCache_UpdateItemByID_LocalMode` - Updating existing items
- ✅ `TestCache_DeleteItemByID_LocalMode` - Deleting items
- ✅ `TestCache_DeleteItemByID_MinConstraint` - Min items constraint enforcement

### 1.5 Concurrent Access Tests
Tests thread-safety and race conditions.

- ✅ `TestCache_ConcurrentWrites` - Multiple goroutines writing simultaneously
- ✅ `TestCache_ConcurrentUpdates` - Concurrent updates to same item
- ✅ `TestCache_ConcurrentDeletes` - Concurrent deletion operations
- ✅ `TestCache_MixedConcurrentOperations` - Mixed read/write/update/delete operations
- ✅ `TestCache_ConcurrentReads` - Multiple goroutines reading simultaneously
- ✅ `TestCache_AtomicStatsUnderConcurrency` - Atomic counter correctness

### 1.6 Edge Cases & Error Handling
Tests boundary conditions and error scenarios.

- ✅ `TestCache_EmptyResourceName` - Empty string resource name
- ✅ `TestCache_InvalidResourceName` - Invalid resource names (special chars, paths, etc.)
- ✅ `TestCache_AddNilItem` - Adding nil items
- ✅ `TestCache_AddItemWithEmptyMap` - Adding empty map items
- ✅ `TestCache_GetByID_WithInvalidIDTypes` - Invalid ID types (nil, objects, arrays, etc.)
- ✅ `TestCache_UpdateItemByID_NonExistent` - Updating non-existent items
- ✅ `TestCache_DeleteItemByID_NonExistent` - Deleting non-existent items

### 1.7 Data Consistency Tests
Tests reproducibility and data integrity.

- ✅ `TestCache_SeedConsistency` - Same seed produces same data
- ✅ `TestCache_DifferentSeedsProduceDifferentData` - Different seeds produce different data

### 1.8 Statistics & Metadata Tests
Tests cache metrics and statistics.

- ✅ `TestCache_GetStats` - Retrieval of cache statistics (hits, misses, etc.)

---

## 2. Cache Redis Tests (`internal/cache/cache_redis_test.go`)

### 2.1 Cache Mode Behavior Tests
Tests how cache behaves with different Redis configurations.

- ✅ `TestCache_RemoteMode_RequiresRedis` - Fallback to local when Redis unavailable in remote mode
- ✅ `TestCache_AllMode_RequiresRedis` - Fallback to local when Redis unavailable in all mode
- ✅ `TestCache_ModeOff_DoesNotUseRedis` - Off mode generates data without caching
- ✅ `TestCache_RedisConnectionFails_Graceful` - Graceful handling of Redis connection failures

### 2.2 Miniredis Integration Tests
Tests using miniredis (mock Redis) to verify Redis operations.

- ✅ `TestMiniredis_BasicOperations` - SET, GET, RPUSH, LRANGE operations with mock Redis
- ✅ `TestMiniredis_KeyExpiry` - Key expiration and TTL with mock Redis
- ✅ `TestMiniredis_PatternMatching` - Key pattern matching and filtering with mock Redis

### 2.3 Redis Simulation Tests
Tests simulating redisstore operations using miniredis.

- ✅ `TestCache_RedisSimulation_SaveAndRetrieve` - Simulating SaveItems/GetItems operations
- ✅ `TestCache_RedisSimulation_UpdateItem` - Simulating UpdateItem operation with LSET
- ✅ `TestCache_RedisSimulation_DeleteItem` - Simulating DeleteItem operation with LREM
- ✅ `TestCache_RedisSimulation_GetByID` - Simulating GetItemByID with pattern matching
- ✅ `TestCache_RedisSimulation_Meta` - Simulating SaveMeta/GetMeta operations

---

## 3. Filters Tests (`internal/filters/filters_test.go`)

### 3.1 Filter Parsing Tests
Tests extraction of filters from query parameters.

- ✅ `TestParseFilters_Equal` - Default equality operator
- ✅ `TestParseFilters_NotEqual` - Not equal operator (!=)
- ✅ `TestParseFilters_NumericOperators` - Greater than, less than, >=, <=
- ✅ `TestParseFilters_StringOperators` - Contains, startsWith, endsWith
- ✅ `TestParseFilters_SkipsReservedParams` - Ignores count, nocache, fresh params

### 3.2 Filter Application Tests
Tests applying filters to data collections.

- ✅ `TestApplyFilters_NoFilters` - No filtering when no filters provided
- ✅ `TestApplyFilters_Equal` - Equality filtering
- ✅ `TestApplyFilters_NotEqual` - Not equal filtering
- ✅ `TestApplyFilters_NumericGreater` - Greater than comparison
- ✅ `TestApplyFilters_NumericLess` - Less than comparison
- ✅ `TestApplyFilters_Contains` - Substring matching
- ✅ `TestApplyFilters_StartsWith` - Prefix matching
- ✅ `TestApplyFilters_EndsWith` - Suffix matching
- ✅ `TestApplyFilters_MultipleFilters` - Multiple filters combined (AND logic)
- ✅ `TestApplyFilters_NoMatches` - No results when no items match
- ✅ `TestApplyFilters_MissingField` - Filtering on non-existent fields

### 3.3 Comparison Helper Tests
Tests type conversion and comparison logic.

- ✅ `TestCompareEqual_DifferentTypes` - Type coercion in equality checks
- ✅ `TestToFloat_Conversions` - Converting various types to float64

### 3.4 Case Sensitivity Tests
Tests case-insensitive string operations.

- ✅ `TestMatchesFilter_CaseInsensitive` - Case-insensitive contains and startsWith

### 3.5 Edge Cases Tests
Tests boundary conditions and special inputs.

- ✅ `TestParseFilters_EmptyValues` - Empty filter values
- ✅ `TestParseFilters_SpecialCharacters` - Field names with special characters
- ✅ `TestParseFilters_MalformedOperators` - Invalid operator syntax
- ✅ `TestParseFilters_UnicodeCharacters` - Unicode in filter values
- ✅ `TestApplyFilters_NilValues` - Handling nil values in data
- ✅ `TestApplyFilters_NestedObjects` - Nested object filtering
- ✅ `TestApplyFilters_EmptyFilterValue` - Empty string filter values
- ✅ `TestApplyFilters_NumericStringMixedComparison` - Mixed numeric/string comparisons

### 3.6 Performance Tests
Tests filter performance with large datasets.

- ✅ `TestApplyFilters_LargeDataset` - Filtering 10,000 items
- ✅ `TestApplyFilters_MultipleComplexFilters` - Complex multi-filter scenarios

### 3.7 Type Conversion Edge Cases
Tests edge cases in type conversion.

- ✅ `TestCompareEqual_EdgeCases` - Nil, boolean, array comparisons
- ✅ `TestToFloat_EdgeCases` - Nil, empty string, negative, scientific notation

---

## 4. Handlers Tests (`internal/handlers/handlers_test.go`)

**Total Tests:** 30+ | **Coverage:** ~75%

### 4.1 GET Collection Tests
Tests retrieving collections of resources.

- ✅ `TestGetCollectionBasic` - Basic GET requests across 10 random resources
- ✅ `TestGetCollectionWithCount` - Count parameter handling (1, 5, 10, 25, 100)
- ✅ `TestGetCollectionWithFilters` - Filter parameter application
- ✅ `TestGetCollectionNoCache` - nocache=true/false behavior
- ✅ `TestGetMetadata` - GET /resource/meta endpoint

### 4.2 Cache Control Tests
Tests cache bypass functionality.

- ✅ `TestCRUDWithNoCacheToggle` - Toggling nocache parameter
- ✅ `TestNoCacheVariations` - Different nocache formats (nocache, fresh)

### 4.3 Error Handling Tests
Tests invalid inputs and error scenarios.

- ✅ `TestInvalidJSONBody` - Malformed JSON in POST requests
- ✅ `TestCountParameterEdgeCases` - Invalid count values (0, -1, abc, 1000)
- ✅ `TestCountCapping` - Count capped at 100

### 4.4 HTTP Method Validation Tests
Tests HTTP-level handler behavior.

- ✅ `TestHandler_InvalidHTTPMethods` - PATCH, OPTIONS, TRACE, CONNECT, HEAD
- ✅ `TestHandler_MalformedJSON_POST` - 8 malformed JSON scenarios (empty, invalid, null, etc.)
- ✅ `TestHandler_MalformedJSON_PUT` - Malformed JSON in PUT requests
- ✅ `TestHandler_MissingContentType` - Missing or wrong Content-Type headers

### 4.5 POST (Create) Tests
Tests resource creation with validation.

- ✅ `TestHandler_POST_ValidCreation` - Valid item creation with auto-generated ID
- ✅ `TestHandler_POST_MaxItemsConstraint` - Max 100 items limit enforcement
- ✅ `TestHandler_POST_ConcurrentCreation` - 10 concurrent POST requests

### 4.6 PUT (Update) Tests
Tests resource updates with validation.

- ✅ `TestHandler_PUT_ValidUpdate` - Valid partial item updates
- ✅ `TestHandler_PUT_InvalidID` - Updating non-existent IDs (multiple invalid formats)

### 4.7 DELETE Tests
Tests resource deletion with constraints.

- ✅ `TestHandler_DELETE_ValidDeletion` - Valid item deletion
- ✅ `TestHandler_DELETE_InvalidID` - Deleting non-existent IDs (multiple formats)
- ✅ `TestHandler_DELETE_MinItemsConstraint` - Minimum 1 item enforcement

### 4.8 Response Headers Tests
Tests HTTP response header correctness.

- ✅ `TestResponseHeaders` - Content-Type, X-Cache headers verification

### 4.9 Data Integrity Tests
Tests data consistency across operations.

- ✅ `TestDataIntegrity` - Data consistency between multiple GET requests (cache vs response)

---

## 5. Utility Handler Tests (`internal/handlers/utility_test.go`)

**Total Tests:** 9 | **Coverage:** ~90%

### 5.1 Echo Endpoint Tests
Tests request reflection and debugging capabilities.

- ✅ `TestUtility_Echo/GET_request` - Reflects GET request details
  - Method, headers, query parameters captured correctly
  - User-Agent and other headers preserved
  - Query parameters parsed and returned

- ✅ `TestUtility_Echo/POST_request_with_JSON_body` - Reflects POST with body
  - POST method captured
  - JSON body parsed and included in response
  - Content-Type headers respected

### 5.2 Delay Endpoint Tests
Tests artificial latency simulation.

- ✅ `TestUtility_Delay/valid_delay` - Delays for specified milliseconds
  - 100ms delay verified with time measurement
  - Response includes delay confirmation
  - Status 200 returned after delay

- ✅ `TestUtility_Delay/invalid_delay` - Rejects invalid delay values
  - Non-numeric values rejected with 400 Bad Request
  - Clear error message provided

- ✅ `TestUtility_Delay/delay_capped_at_30_seconds` - Maximum delay enforcement
  - Delays > 30 seconds capped at 30s
  - Response indicates capped value (30000ms)
  - Prevents abuse with excessive delays

### 5.3 Random Delay Tests
Tests random latency between min/max bounds.

- ✅ `TestUtility_DelayRandom/random_delay_with_defaults` - Default range behavior
  - Default: 100ms-2000ms range
  - Actual delay falls within expected bounds
  - Randomness verified over multiple executions

- ✅ `TestUtility_DelayRandom/random_delay_with_custom_range` - Custom min/max
  - Custom ranges respected (e.g., 50-150ms)
  - Min and max parameters parsed correctly
  - Delay always between min and max

### 5.4 Status Code Tests
Tests returning arbitrary HTTP status codes.

- ✅ `TestUtility_Status/status_*` - Multiple status codes (200, 201, 400, 404, 500, 503)
  - Each status code returns correctly
  - Response includes status categorization (2xx, 4xx, 5xx)
  - Clear message describing the status category

- ✅ `TestUtility_Status/invalid_status_code` - Out of range codes rejected
  - Codes outside 100-599 range return 400
  - Error message explains valid range

### 5.5 Validation Error Tests
Tests sample validation error responses.

- ✅ `TestUtility_ErrorValidation` - Returns 422 validation error
  - Status 422 Unprocessable Entity
  - Multiple validation errors in array format
  - Each error includes: field, message, code
  - Sample errors: email required, age range, username length

### 5.6 Flaky Endpoint Tests
Tests probabilistic success/failure.

- ✅ `TestUtility_Flaky/flaky_with_100%_success_rate` - Always succeeds
  - successRate=1.0 → 200 OK every time
  - Response includes success=true

- ✅ `TestUtility_Flaky/flaky_with_0%_success_rate` - Always fails
  - successRate=0.0 → 503 every time
  - Response includes success=false
  - Clear error message with retry hint

- ✅ `TestUtility_Flaky/flaky_with_default_rate` - Probabilistic behavior
  - Default 50% success rate
  - 20 requests produce mix of successes and failures
  - Randomness validated (not all success or all failure)

### 5.7 Chaos Engineering Tests
Tests random status codes and response shapes.

- ✅ `TestUtility_Chaos` - Random status codes
  - Multiple different status codes generated (200, 400, 401, 403, 404, 500, 502, 503)
  - 50 requests produce variety (not stuck on one code)
  - All responses are valid JSON

- ✅ `TestUtility_ChaosResponseShapes` - Random response structures
  - 5 different response shapes generated
  - Responses vary between: objects with data/result/nested, arrays, simple key-value
  - 30 requests produce multiple different shapes
  - All shapes are valid JSON

### Testing Utilities Purpose
These endpoints enable developers to:
- **Echo:** Debug request details and inspect what server receives
- **Delay:** Test timeout handling and loading states (1-30000ms)
- **DelayRandom:** Simulate variable network latency
- **Status:** Test error handling for any HTTP status code (100-599)
- **ErrorValidation:** Test form validation error handling (422 responses)
- **Flaky:** Test retry logic and resilience (configurable success rate)
- **Chaos:** Test error handling with unpredictable responses (chaos engineering)

---

## 6. Schema Validation Tests (`internal/handlers/schema_validation_test.go`)

**Total Tests:** 10 | **Coverage:** ~85%

### 5.1 Type Validation Tests
Tests JSON schema type checking with multiple scenarios.

- ✅ `TestValidation_TypeChecking` - String, number, boolean type validation
  - Valid types accepted
  - Invalid types rejected (e.g., age as string instead of number)
  - Type mismatches produce clear error messages

### 5.2 Required Fields Tests
Tests required property enforcement.

- ✅ `TestValidation_RequiredFields` - Missing required fields rejection
  - All required fields present → Success
  - Missing username → Rejection
  - Missing email → Rejection

### 5.3 Format Validation Tests
Tests string format validation (email, URL, etc.).

- ✅ `TestValidation_EmailFormat` - Email format validation
  - Valid emails: `john@example.com`, `john.doe@mail.example.com`
  - Invalid emails: `johnexample.com`, `notanemail`

### 5.4 Enum Validation Tests
Tests enumeration constraints.

- ✅ `TestValidation_EnumConstraint` - Enum value validation
  - Valid values: `male`, `female`
  - Invalid values: `other`, `unknown`

### 5.5 Numeric Constraints Tests
Tests numeric boundary validation.

- ✅ `TestValidation_NumericConstraints` - Min/max value constraints (age 18-100)
  - Boundary values: 18 (min), 100 (max) → Accepted
  - Below minimum: 17 → Rejected
  - Above maximum: 101 → Rejected

### 5.6 String Constraints Tests
Tests string length validation.

- ✅ `TestValidation_StringConstraints` - MinLength/maxLength (username 3-20 chars)
  - Boundary values: "abc" (3), "abcdefghijklmnopqrst" (20) → Accepted
  - Too short: "ab" → Rejected
  - Too long: "abcdefghijklmnopqrstu" (21) → Rejected

### 5.7 Array Constraints Tests
Tests array validation rules.

- ✅ `TestValidation_ArrayConstraints` - MinItems, maxItems, uniqueItems
  - Within bounds → Accepted
  - Below minimum → Rejected
  - Above maximum → Rejected
  - Duplicate items (when uniqueItems=true) → Rejected

### 5.8 Additional Properties Tests
Tests schema strictness (additionalProperties: false).

- ✅ `TestValidation_AdditionalProperties` - Rejection of extra fields
  - Only schema-defined properties → Accepted
  - Extra/unknown properties → Rejected

### 5.9 Update Validation Tests
Tests validation during PUT operations.

- ✅ `TestValidation_UpdateOperation` - PUT request validation
  - Valid updates: email, age changes → Accepted
  - Invalid email format → Rejected
  - Age below minimum → Rejected
  - Invalid gender enum → Rejected
  - Username too short → Rejected

### 5.10 ID Field Validation Tests
Tests ID field handling in POST requests.

- ✅ `TestValidation_DuplicateID` - ID field rules
  - POST with ID field → Rejected (IDs auto-generated)
  - POST without ID → Accepted with auto-generated ID

---

## 6. Validation Tests (`internal/handlers/validation_test.go`)

**Total Tests:** 6 | **Coverage:** ~80%

### 6.1 POST Validation Tests
Tests validation during resource creation.

- ✅ `TestPostCollectionValidatesRequiredFields` - Required fields enforcement
  - Missing required fields (username, email) → 400 Bad Request
  - Clear error messages indicating which field is missing
  
- ✅ `TestPostCollectionValidatesPropertyTypes` - Property type validation
  - Wrong type for numeric field (price as string) → 400 Bad Request
  - Type mismatch errors with field name
  
- ✅ `TestPostCollectionAcceptsValidItem` - Valid item acceptance
  - Correctly formatted item → 201 Created
  - Item added to cache
  - Response contains created data

### 6.2 PUT Validation Tests
Tests validation during resource updates.

- ✅ `TestPutSingleValidatesTypes` - Type validation in updates
  - Invalid type for update (price as string) → 400 Bad Request
  - Validation applied to partial updates
  
- ✅ `TestPutSingleAcceptsValidUpdate` - Valid update acceptance
  - Valid partial update → 200 OK
  - Updated fields reflected in response
  - ID preserved during update

### 6.3 Multi-Resource Validation Tests
Tests validation across different resource types.

- ✅ `TestSchemaValidationForMultipleResources` - Cross-resource validation
  - User schema validation (ID should be numeric)
  - Post schema validation (userId should be numeric)
  - Each resource uses its own schema rules

---

## 7. CORS Middleware Tests (`internal/middleware/cors_test.go`)

**Total Tests:** 10 | **Coverage:** ~95%

### 7.1 CORS Configuration Tests
Tests CORS setup and configuration.

- ✅ `TestSetupCORS_ConfigurationExists` - CORS handler creation
  - Handler is non-nil and properly initialized
  
- ✅ `TestCORS_AllowedOrigins` - Access-Control-Allow-Origin header
  - Wildcard (*) for all origins
  - Works with any Origin header
  
- ✅ `TestCORS_MaxAge` - Access-Control-Max-Age header
  - Max age header present on preflight requests
  - Proper cache duration for CORS preflight

### 7.2 Preflight Request Tests
Tests OPTIONS request handling.

- ✅ `TestCORS_PreflightRequest` - OPTIONS request handling
  - OPTIONS requests return 200/204
  - Access-Control-Allow-Methods header set
  - Proper CORS headers in response

### 7.3 Allowed Methods Tests
Tests HTTP method permissions.

- ✅ `TestCORS_AllowedMethods` - All standard REST methods
  - GET, POST, PUT, PATCH, DELETE
  - Each method tested individually
  - Access-Control-Allow-Methods header verification

### 7.4 Allowed Headers Tests
Tests allowed request headers.

- ✅ `TestCORS_AllowedHeaders` - Standard and custom headers
  - Content-Type (for JSON requests)
  - Authorization (for auth tokens)
  - X-No-Cache (custom cache control)

### 7.5 Exposed Headers Tests
Tests response header exposure.

- ✅ `TestCORS_ExposedHeaders` - Custom response headers
  - Link header preserved
  - X-Cache header accessible to clients

### 7.6 Credentials Tests
Tests credential handling.

- ✅ `TestCORS_CredentialsNotAllowed` - Security verification
  - Access-Control-Allow-Credentials not set to true
  - No credential support (public API)

### 7.7 Multiple Origins Tests
Tests handling of various origins.

- ✅ `TestCORS_MultipleOrigins` - Origin diversity
  - localhost:3000 (development)
  - https://example.com (production)
  - https://test.org (staging)
  - All origins receive wildcard (*) response

---

## Test Coverage Summary

### Coverage by Package

| Package | Statement Coverage | Critical Paths | Concurrent Tests |
|---------|-------------------|----------------|------------------|
| `cache` | ~87% | ✅ CRUD, modes, warmup | ✅ 50+ goroutines |
| `filters` | ~92% | ✅ All operators | ✅ 10K items |
| `handlers` | ~75% | ✅ HTTP CRUD | ✅ Concurrent POST |
| `handlers (utility)` | ~90% | ✅ Testing utils | ✅ Probabilistic |
| `middleware` | ~95% | ✅ CORS | N/A |
| **Overall** | **~86%** | **Full stack** | **Thread-safe** |

### What We Test

#### ✅ **Functionality** (100% core features)
- Cache operations (CRUD) with all modes
- Data filtering with 8+ operators
- HTTP request handling (GET, POST, PUT, DELETE)
- JSON schema validation (types, formats, constraints)
- CORS policies and preflight requests

#### ✅ **Performance** (Validated)
- Concurrent access with 50+ goroutines
- Large dataset filtering: 10,000 items in <100ms
- Cache warmup time tracking
- Response time validation

#### ✅ **Reliability** (Thread-safe)
- Thread safety verified with race detector
- Race condition prevention
- Atomic operations for counters
- Data consistency across concurrent operations

#### ✅ **Error Handling** (Comprehensive)
- Invalid inputs and malformed JSON (8+ scenarios)
- Missing required fields
- Non-existent resources (404 handling)
- Constraint violations (min/max items, string length, numeric bounds)
- Type mismatches with clear error messages

#### ✅ **Edge Cases** (Robust)
- Empty and nil values
- Special characters and Unicode
- Boundary values (min/max for age, string lengths)
- Type coercion (numeric strings vs numbers)
- Invalid operators and malformed filters

#### ✅ **Configuration** (All modes)
- Cache modes: off, local, remote, all
- Fallback behavior when Redis unavailable
- Mode switching at runtime
- Graceful degradation

### What We Don't Test (Out of Scope)

#### ❌ **Infrastructure** (Production environment)
- Real Redis server connections (we use miniredis)
- Network failures and timeouts
- Actual HTTP server on port 8080
- File system operations (schemas embedded)
- Docker container behavior

#### ❌ **External Dependencies** (None exist)
- Third-party API integrations
- Authentication providers (public API)
- External databases
- Payment gateways

#### ❌ **Deployment** (Platform-specific)
- Container orchestration (Kubernetes, Docker Swarm)
- Load balancing and service discovery
- Fly.io deployment specifics
- CDN caching behavior
- SSL/TLS certificate validation

#### ❌ **Observability** (Not in test scope)
- Log output format verification
- Metrics collection and exporters
- Distributed tracing
- Error tracking integration (Sentry, etc.)

---

## Testing Strategy

### 1. Test Isolation
- Each test is independent (no shared state)
- Setup/teardown handled automatically
- No external service dependencies
- Deterministic with seeded random data (seed=42)

### 2. Mock Strategy
- **Redis:** Miniredis (in-process mock server)
- **HTTP:** httptest.ResponseRecorder
- **Time:** Fixed timestamps where needed
- **Random Data:** Seeded with gofakeit

### 3. Performance Targets
- All tests complete in <5 seconds
- Individual test <500ms
- Concurrent tests handle 50+ goroutines
- Filter 10K items in <100ms

### 4. Race Detection
```bash
go test -race ./internal/...
```
All concurrent tests pass race detector.

---

## Running Tests

### Basic Commands
```bash
cd backend

# All tests
go test ./internal/...

# Specific package
go test ./internal/cache
go test ./internal/filters
go test ./internal/handlers
go test ./internal/middleware

# With coverage report
go test ./internal/... -cover

# Verbose output
go test ./internal/... -v

# Specific test
go test ./internal/cache -run TestCache_ConcurrentWrites

# With race detection
go test ./internal/... -race

# Parallel execution
go test ./internal/... -parallel 4
```

### Coverage Reports
```bash
# Generate coverage profile
go test ./internal/... -coverprofile=coverage.out

# View in browser
go tool cover -html=coverage.out

# View in terminal
go tool cover -func=coverage.out
```

### Makefile Commands
```bash
# From project root
make api-test

# From backend directory
cd backend
make test
```

---

## Test Organization

### File Structure
```
backend/internal/
├── cache/
│   ├── cache.go                    # Implementation (~500 lines)
│   ├── cache_test.go               # 35 tests (~800 lines)
│   └── cache_redis_test.go         # 12 tests (~400 lines)
├── filters/
│   ├── filters.go                  # Implementation (~200 lines)
│   └── filters_test.go             # 31 tests (~600 lines)
├── handlers/
│   ├── dynamic.go                  # Implementation (~400 lines)
│   ├── handlers_test.go            # 30+ tests (~1000 lines)
│   ├── schema_validation_test.go   # 10 tests (~450 lines)
│   └── validation_test.go          # 6 tests (~250 lines)
└── middleware/
    ├── cors.go                     # Implementation (~50 lines)
    └── cors_test.go                # 10 tests (~200 lines)
```

### Test Naming Convention
```
Test<Component>_<Scenario>_<ExpectedBehavior>
```

**Examples:**
- `TestCache_Get_LocalMode_Hit` - Cache hit in local mode
- `TestApplyFilters_MultipleFilters` - Multiple filters combined
- `TestValidation_RequiredFields` - Required field enforcement
- `TestHandler_POST_MaxItemsConstraint` - Max items limit

### Test Structure Pattern
```go
func TestComponent_Scenario(t *testing.T) {
    // 1. SETUP - Create test environment
    handler, cache := setupTestHandler(t)
    
    // 2. EXECUTE - Perform operation
    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()
    handler.GetCollection("users")(w, req)
    
    // 3. ASSERT - Verify results
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
    
    // Decode and verify response
    var data []map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &data)
    
    if len(data) == 0 {
        t.Error("Expected non-empty data")
    }
    
    // 4. CLEANUP - Automatic via defer or t.Cleanup()
}
```

### Test Helpers
```go
// Setup helpers (in *_test.go files)
func setupTestHandler(t *testing.T) (*DynamicHandler, *cache.Cache)
func setupTestCache(t *testing.T, mode CacheMode) (*Cache, *schema.Registry)
func setupRedisTest(t *testing.T, mode CacheMode) (*Cache, *miniredis.Miniredis)

// Utility helpers
func createCtxWithParams(params map[string]string) context.Context
func randomResources(t *testing.T, handler *DynamicHandler, n int) []string
func findFilter(filters []Filter, field string) *Filter
```

---

## Test Metrics

### Execution Time (Approximate)
```
Cache tests:           ~1.2s (35 tests)
Cache Redis tests:     ~0.4s (12 tests)
Filter tests:          ~0.5s (31 tests)
Handler tests:         ~2.9s (30+ tests)
Utility handler tests: ~31.5s (9 tests - includes actual 30s delay test)
Schema validation:     ~0.3s (10 tests)
Validation tests:      ~0.2s (6 tests)
CORS tests:            ~0.2s (10 tests)
─────────────────────────────────────
Total:                 ~37.2s (143 tests)
```

### Coverage Gaps (Areas for Improvement)

**Handlers (~75% coverage):**
- GET by ID error scenarios
- More concurrent update scenarios
- Stress testing with high request volumes

**Validation (~80% coverage):**
- More complex nested object validation
- Array of objects validation
- Pattern/regex validation

**Future Tests Planned:**
- Benchmark tests for performance tracking
- Fuzz testing for edge case discovery
- Integration tests for end-to-end workflows
- Load testing for production readiness

---

## Contributing

### Adding New Tests

When adding new tests:

1. ✅ **Follow naming convention:** `Test<Component>_<Scenario>`
2. ✅ **Add to appropriate file:** Group related tests together
3. ✅ **Update this document:** Add test description to relevant section
4. ✅ **Ensure isolation:** No shared state between tests
5. ✅ **Test both paths:** Success and failure scenarios
6. ✅ **Include edge cases:** Nil, empty, boundary values
7. ✅ **Use table-driven tests:** For multiple similar scenarios
8. ✅ **Add helpful logs:** Use `t.Logf()` for debugging
9. ✅ **Run race detector:** `go test -race`
10. ✅ **Verify coverage:** Aim for >80% for new code

### Example: Adding a New Test
```go
func TestCache_NewFeature(t *testing.T) {
    // Setup
    cache, _ := setupTestCache(t, CacheModeLocal)
    cache.Warmup()
    
    // Test scenarios
    testCases := []struct {
        name     string
        input    interface{}
        expected bool
    }{
        {"valid input", "test", true},
        {"empty input", "", false},
        {"nil input", nil, false},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := cache.NewFeature(tc.input)
            if result != tc.expected {
                t.Errorf("Expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```

---

## Last Updated
December 11, 2025
