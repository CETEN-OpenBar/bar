# RefillsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getAccountRefills**](#getaccountrefills) | **GET** /accounts/{account_id}/refills | |
|[**getRefills**](#getrefills) | **GET** /refills | |
|[**getSelfRefills**](#getselfrefills) | **GET** /account/refills | |
|[**markDeleteRefill**](#markdeleterefill) | **DELETE** /accounts/{account_id}/refills/{refill_id} | |
|[**patchRefillId**](#patchrefillid) | **PATCH** /accounts/{account_id}/refills/{refill_id} | |
|[**postRefill**](#postrefill) | **POST** /accounts/{account_id}/refills | |

# **getAccountRefills**
> GetRefills200Response getAccountRefills()

Get all refills of an account

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let accountId: string; //ID or CardID of the account (default to undefined)
let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of transactions per page (optional) (default to undefined)
let startDate: string; //Start date of the refill (optional) (default to undefined)
let endDate: string; //End date of the refill (optional) (default to undefined)

const { status, data } = await apiInstance.getAccountRefills(
    accountId,
    page,
    limit,
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID or CardID of the account | defaults to undefined|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the refill | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the refill | (optional) defaults to undefined|


### Return type

**GetRefills200Response**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRefills**
> GetRefills200Response getRefills()

Get all refills

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of transactions per page (optional) (default to undefined)
let startDate: string; //Start date of the refill (optional) (default to undefined)
let endDate: string; //End date of the refill (optional) (default to undefined)

const { status, data } = await apiInstance.getRefills(
    page,
    limit,
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the refill | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the refill | (optional) defaults to undefined|


### Return type

**GetRefills200Response**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getSelfRefills**
> GetRefills200Response getSelfRefills()

Get all refills

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of transactions per page (optional) (default to undefined)
let startDate: string; //Start date of the refill (optional) (default to undefined)
let endDate: string; //End date of the refill (optional) (default to undefined)

const { status, data } = await apiInstance.getSelfRefills(
    page,
    limit,
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the refill | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the refill | (optional) defaults to undefined|


### Return type

**GetRefills200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**401** | Not authorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markDeleteRefill**
> markDeleteRefill()

Cancels a refill

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let refillId: string; //ID of the refill (default to undefined)

const { status, data } = await apiInstance.markDeleteRefill(
    accountId,
    refillId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **refillId** | [**string**] | ID of the refill | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** |  |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account or refill not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchRefillId**
> Refill patchRefillId()

Update refill\'s state

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let refillId: string; //ID of the refill (default to undefined)
let state: RefillState; //New state of the refill (optional) (default to undefined)
let type: RefillType; //New type of the refill (optional) (default to undefined)

const { status, data } = await apiInstance.patchRefillId(
    accountId,
    refillId,
    state,
    type
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **refillId** | [**string**] | ID of the refill | defaults to undefined|
| **state** | [**RefillState**] | New state of the refill | (optional) defaults to undefined|
| **type** | [**RefillType**] | New type of the refill | (optional) defaults to undefined|


### Return type

**Refill**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**400** | Bad request |  -  |
|**401** | Not connected |  -  |
|**403** | Forbidden |  -  |
|**404** | Account or refill not found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **postRefill**
> Refill postRefill()

Create a new refill

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let accountId: string; //ID or CardID of the account (default to undefined)
let amount: number; //Amount of the refill (default to undefined)
let type: RefillType; //Type of the refill (default to undefined)

const { status, data } = await apiInstance.postRefill(
    accountId,
    amount,
    type
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID or CardID of the account | defaults to undefined|
| **amount** | [**number**] | Amount of the refill | defaults to undefined|
| **type** | [**RefillType**] | Type of the refill | defaults to undefined|


### Return type

**Refill**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** |  |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

