# RefillsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getAccountRefills**](#getaccountrefills) | **GET** /accounts/{account_id}/refills | |
|[**getPendingRemoteRefills**](#getpendingremoterefills) | **GET** /account/remote-refills/pending | |
|[**getRefills**](#getrefills) | **GET** /refills | |
|[**getRemoteRefillStatus**](#getremoterefillstatus) | **GET** /remote-refills/status | |
|[**getRemoteRefills**](#getremoterefills) | **GET** /remote-refills | |
|[**getSelfRefills**](#getselfrefills) | **GET** /account/refills | |
|[**markDeleteRefill**](#markdeleterefill) | **DELETE** /accounts/{account_id}/refills/{refill_id} | |
|[**patchRefillId**](#patchrefillid) | **PATCH** /accounts/{account_id}/refills/{refill_id} | |
|[**postRefill**](#postrefill) | **POST** /accounts/{account_id}/refills | |
|[**selfValidateRemoteRefill**](#selfvalidateremoterefill) | **POST** /account/remote-refills/validate | |
|[**startRemoteRefill**](#startremoterefill) | **POST** /account/remote-refills/start | |
|[**verifyRemoteRefill**](#verifyremoterefill) | **POST** /remote-refills/verify | |

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

# **getPendingRemoteRefills**
> GetPendingRemoteRefills200Response getPendingRemoteRefills()

Get all pending remote refills for your account

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

const { status, data } = await apiInstance.getPendingRemoteRefills();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**GetPendingRemoteRefills200Response**

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

# **getRemoteRefillStatus**
> getRemoteRefillStatus()

Get the status of the remote refill system

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

const { status, data } = await apiInstance.getRemoteRefillStatus();
```

### Parameters
This endpoint does not have any parameters.


### Return type

void (empty response body)

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | The remote refill subsystem is working |  -  |
|**401** | Not authorized |  -  |
|**503** | Remote refills are not available |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRemoteRefills**
> GetRemoteRefills200Response getRemoteRefills()

Get all remote refills

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
let state: RemoteRefillState; //State of the refill (optional) (default to undefined)
let accountName: string; //Filter by account name (optional) (default to undefined)

const { status, data } = await apiInstance.getRemoteRefills(
    page,
    limit,
    startDate,
    endDate,
    state,
    accountName
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the refill | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the refill | (optional) defaults to undefined|
| **state** | **RemoteRefillState** | State of the refill | (optional) defaults to undefined|
| **accountName** | [**string**] | Filter by account name | (optional) defaults to undefined|


### Return type

**GetRemoteRefills200Response**

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

[admin_auth](../README.md#admin_auth)

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

# **selfValidateRemoteRefill**
> Refill selfValidateRemoteRefill()

Validate a remote refill

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let checkoutIntentId: number; //HelloAsso checkout intent id to validate (default to undefined)

const { status, data } = await apiInstance.selfValidateRemoteRefill(
    checkoutIntentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **checkoutIntentId** | [**number**] | HelloAsso checkout intent id to validate | defaults to undefined|


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
|**200** | Payment successfully verified, money has been added to the account |  -  |
|**409** | Remote refill already validated |  -  |
|**404** | Checkout id unknown or not associated with this account |  -  |
|**402** | HelloAsso did not validate the payment yet |  -  |
|**401** | Not authorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **startRemoteRefill**
> StartRemoteRefill200Response startRemoteRefill()

Start a remote refill

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let amount: number; //Amount of the refill in cents (default to undefined)

const { status, data } = await apiInstance.startRemoteRefill(
    amount
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **amount** | [**number**] | Amount of the refill in cents | defaults to undefined|


### Return type

**StartRemoteRefill200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**400** | Invalid amount |  -  |
|**401** | Not authorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **verifyRemoteRefill**
> Refill verifyRemoteRefill()

Verify a remote refill

### Example

```typescript
import {
    RefillsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RefillsApi(configuration);

let id: string; //Remote Refill id (default to undefined)

const { status, data } = await apiInstance.verifyRemoteRefill(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Remote Refill id | defaults to undefined|


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
|**200** | Payment successfully verified, money has been added to the account |  -  |
|**409** | Remote refill already validated |  -  |
|**404** | Refill id unknown |  -  |
|**402** | HelloAsso did not validate the payment yet |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

