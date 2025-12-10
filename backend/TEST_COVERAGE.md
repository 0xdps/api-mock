# Backend Test Coverage Documentation

This document provides a comprehensive overview of all test cases in the API-Mockly backend.

## Test Statistics

| Component | Test File | Test Count | Coverage Focus |
|-----------|-----------|------------|----------------|
| Cache (Local) | `cache_test.go` | 35 | Cache operations, concurrency, edge cases |
| Cache (Redis Mock) | `cache_redis_test.go` | 12 | Mock Redis integration, Redis simulation |
| Filters | `filters_test.go` | 31 | Query parameter filtering |
| Handlers | `handlers_test.go` | 20+ | HTTP request handling |
| Schema Validation | `schema_validation_test.go` | 9 | JSON schema validation |
| Validation | `validation_test.go` | 6 | Data validation |
| CORS Middleware | `cors_test.go` | 10 | CORS configuration |
| **Total** | **7 files** | **~123 tests** | **Full stack coverage** |

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

### 4.1 GET Collection Tests
Tests retrieving collections of resources.

- ✅ `TestGetCollectionBasic` - Basic GET requests across 10 random resources
- ✅ `TestGetCollectionWithCount` - Count parameter handling (1, 5, 10, 25, 100)
- ✅ `TestGetCollectionWithFilters` - Filter parameter application
- ✅ `TestGetCollectionNoCache` - nocache=true/false behavior
- ✅ `TestGetMetadata` - GET /resource/meta endpoint

### 4.2 Cache Toggle Tests
Tests cache bypass functionality.

- ✅ `TestCRUDWithNoCacheToggle` - Toggling nocache parameter
- ✅ `TestNoCacheVariations` - Different nocache formats (nocache, fresh)

### 4.3 Error Handling Tests
Tests invalid inputs and error scenarios.

- ✅ `TestInvalidJSONBody` - Malformed JSON in POST requests
- ✅ `TestCountParameterEdgeCases` - Invalid count values (0, -1, abc, 1000)
- ✅ `TestCountCapping` - Count capped at 100

### 4.4 HTTP Handler Tests
Tests HTTP-level handler behavior.

- ✅ `TestHandler_InvalidHTTPMethods` - PATCH, OPTIONS, TRACE, CONNECT, HEAD
- ✅ `TestHandler_MalformedJSON_POST` - Various malformed JSON scenarios
- ✅ `TestHandler_MalformedJSON_PUT` - Malformed JSON in PUT requests
- ✅ `TestHandler_MissingContentType` - Missing or wrong Content-Type headers

### 4.5 POST (Create) Tests
Tests resource creation.

- ✅ `TestHandler_POST_ValidCreation` - Valid item creation
- ✅ `TestHandler_POST_MaxItemsConstraint` - Max items limit enforcement
- ✅ `TestHandler_POST_ConcurrentCreation` - Concurrent POST requests

### 4.6 PUT (Update) Tests
Tests resource updates.

- ✅ `TestHandler_PUT_ValidUpdate` - Valid item updates
- ✅ `TestHandler_PUT_InvalidID` - Updating non-existent IDs

### 4.7 DELETE Tests
Tests resource deletion.

- ✅ `TestHandler_DELETE_ValidDeletion` - Valid item deletion
- ✅ `TestHandler_DELETE_InvalidID` - Deleting non-existent IDs
- ✅ `TestHandler_DELETE_MinItemsConstraint` - Minimum items enforcement

### 4.8 Response Headers Tests
Tests HTTP response header correctness.

- ✅ `TestResponseHeaders` - Content-Type, X-Cache headers

### 4.9 Data Integrity Tests
Tests data consistency across operations.

- ✅ `TestDataIntegrity` - Data consistency between multiple GET requests

---

## 5. Schema Validation Tests (`internal/handlers/schema_validation_test.go`)

### 5.1 Type Validation Tests
Tests JSON schema type checking.

- ✅ `TestValidation_TypeChecking` - String, number, boolean type validation

### 5.2 Required Fields Tests
Tests required property enforcement.

- ✅ `TestValidation_RequiredFields` - Missing required fields rejection

### 5.3 Format Validation Tests
Tests string format validation.

- ✅ `TestValidation_EmailFormat` - Email format validation

### 5.4 Enum Validation Tests
Tests enumeration constraints.

- ✅ `TestValidation_EnumConstraint` - Enum value validation

### 5.5 Numeric Constraints Tests
Tests numeric boundary validation.

- ✅ `TestValidation_NumericConstraints` - Min, max value constraints

### 5.6 String Constraints Tests
Tests string length validation.

- ✅ `TestValidation_StringConstraints` - MinLength, maxLength constraints

### 5.7 Array Constraints Tests
Tests array validation rules.

- ✅ `TestValidation_ArrayConstraints` - MinItems, maxItems, uniqueItems

### 5.8 Additional Properties Tests
Tests schema strictness.

- ✅ `TestValidation_AdditionalProperties` - Rejection of extra fields

### 5.9 Update Validation Tests
Tests validation during updates.

- ✅ `TestValidation_UpdateOperation` - PUT request validation

### 5.10 Duplicate ID Tests
Tests ID field handling.

- ✅ `TestValidation_DuplicateID` - Rejection of user-provided IDs in POST, auto-generation

---

## 6. Validation Tests (`internal/handlers/validation_test.go`)

### 6.1 POST Validation Tests
Tests validation during resource creation.

- ✅ `TestPostCollectionValidatesRequiredFields` - Required fields enforcement
- ✅ `TestPostCollectionValidatesPropertyTypes` - Property type validation
- ✅ `TestPostCollectionAcceptsValidItem` - Valid item acceptance

### 6.2 PUT Validation Tests
Tests validation during resource updates.

- ✅ `TestPutSingleValidatesTypes` - Type validation in updates
- ✅ `TestPutSingleAcceptsValidUpdate` - Valid update acceptance

### 6.3 Multi-Resource Validation Tests
Tests validation across different resource types.

- ✅ `TestSchemaValidationForMultipleResources` - Validation for various resources (user, post, etc.)

---

## 7. CORS Middleware Tests (`internal/middleware/cors_test.go`)

### 7.1 CORS Configuration Tests
Tests CORS setup and configuration.

- ✅ `TestSetupCORS_ConfigurationExists` - CORS handler creation
- ✅ `TestCORS_AllowedOrigins` - Access-Control-Allow-Origin header
- ✅ `TestCORS_MaxAge` - Access-Control-Max-Age header

### 7.2 Preflight Request Tests
Tests OPTIONS request handling.

- ✅ `TestCORS_PreflightRequest` - OPTIONS request handling

### 7.3 Allowed Methods Tests
Tests HTTP method permissions.

- ✅ `TestCORS_AllowedMethods` - GET, POST, PUT, PATCH, DELETE

### 7.4 Allowed Headers Tests
Tests allowed request headers.

- ✅ `TestCORS_AllowedHeaders` - Content-Type, Authorization, X-No-Cache

### 7.5 Exposed Headers Tests
Tests response header exposure.

- ✅ `TestCORS_ExposedHeaders` - Link, X-Cache headers

### 7.6 Credentials Tests
Tests credential handling.

- ✅ `TestCORS_CredentialsNotAllowed` - Credentials not enabled

### 7.7 Multiple Origins Tests
Tests handling of various origins.

- ✅ `TestCORS_MultipleOrigins` - Various origin domains

---

## Test Coverage Summary

### What We Test

#### ✅ **Functionality**
- Cache operations (CRUD)
- Data filtering and querying
- HTTP request handling
- Schema validation
- CORS policies

#### ✅ **Performance**
- Concurrent access (50+ goroutines)
- Large dataset filtering (10,000 items)
- Cache warmup time
- Response time validation

#### ✅ **Reliability**
- Thread safety
- Race condition prevention
- Atomic operations
- Data consistency

#### ✅ **Error Handling**
- Invalid inputs
- Malformed JSON
- Missing required fields
- Non-existent resources
- Constraint violations

#### ✅ **Edge Cases**
- Empty values
- Nil values
- Special characters
- Unicode characters
- Boundary values
- Type coercion

#### ✅ **Configuration**
- Cache modes (off, local, remote, all)
- Fallback behavior
- Mode switching
- Redis unavailability

### What We Don't Test

#### ❌ **Infrastructure**
- Real Redis server connections (we use miniredis/mock Redis)
- Network failures and timeouts
- Real database connections
- File system operations

#### ❌ **External Dependencies**
- Third-party API integrations
- Authentication providers
- External services

#### ❌ **Deployment**
- Container orchestration
- Load balancing
- Service discovery
- Health checks

---

## Running Tests

### Run All Tests
```bash
cd backend
go test ./internal/...
```

### Run Specific Package
```bash
go test ./internal/cache
go test ./internal/filters
go test ./internal/handlers
go test ./internal/middleware
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

---

## Test Organization

### Test Naming Convention
```
Test<Component>_<Scenario>_<ExpectedBehavior>
```

Examples:
- `TestCache_Get_LocalMode_Hit`
- `TestApplyFilters_MultipleFilters`
- `TestValidation_RequiredFields`

### Test Structure
```go
func TestComponent_Scenario(t *testing.T) {
    // 1. Setup
    handler, cache := setupTestHandler(t)
    
    // 2. Execute
    result := handler.DoSomething()
    
    // 3. Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
    
    // 4. Cleanup (if needed)
}
```

---

## Contributing

When adding new tests:

1. ✅ Follow the naming convention
2. ✅ Add to appropriate test file
3. ✅ Update this documentation
4. ✅ Ensure test is isolated and repeatable
5. ✅ Test both success and failure paths
6. ✅ Include edge cases
7. ✅ Use table-driven tests for multiple scenarios

---

## Last Updated
December 10, 2025
