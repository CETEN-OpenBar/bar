# StarsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getAccountStarring**](#getaccountstarring) | **GET** /accounts/{account_id}/stars | |
|[**getSelfStarring**](#getselfstarring) | **GET** /account/stars | |
|[**getStarrings**](#getstarrings) | **GET** /stars | |
|[**markDeleteStarring**](#markdeletestarring) | **DELETE** /accounts/{account_id}/stars/{starring_id} | |
|[**patchStarringId**](#patchstarringid) | **PATCH** /accounts/{account_id}/stars/{starring_id} | |
|[**postStarring**](#poststarring) | **POST** /accounts/{account_id}/stars | |

# **getAccountStarring**
> GetStarrings200Response getAccountStarring()

Get all stars donations of an account

### Example

```typescript
import {
    StarsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new StarsApi(configuration);

let accountId: string; //ID or CardID of the account (default to undefined)
let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of donations per page (optional) (default to undefined)
let startDate: string; //Start date of the donation (optional) (default to undefined)
let endDate: string; //End date of the donation (optional) (default to undefined)

const { status, data } = await apiInstance.getAccountStarring(
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
| **limit** | [**number**] | Number of donations per page | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the donation | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the donation | (optional) defaults to undefined|


### Return type

**GetStarrings200Response**

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

# **getSelfStarring**
> GetStarrings200Response getSelfStarring()

Get all stars donations

### Example

```typescript
import {
    StarsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new StarsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of donations per page (optional) (default to undefined)
let startDate: string; //Start date of the donation (optional) (default to undefined)
let endDate: string; //End date of the donation (optional) (default to undefined)

const { status, data } = await apiInstance.getSelfStarring(
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
| **limit** | [**number**] | Number of donations per page | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the donation | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the donation | (optional) defaults to undefined|


### Return type

**GetStarrings200Response**

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

# **getStarrings**
> GetStarrings200Response getStarrings()

Get all stars donations

### Example

```typescript
import {
    StarsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new StarsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of donations per page (optional) (default to undefined)
let name: string; //Filter by account name (optional) (default to undefined)
let startDate: string; //Start date of the donations (optional) (default to undefined)
let endDate: string; //End date of the donations (optional) (default to undefined)

const { status, data } = await apiInstance.getStarrings(
    page,
    limit,
    name,
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of donations per page | (optional) defaults to undefined|
| **name** | [**string**] | Filter by account name | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date of the donations | (optional) defaults to undefined|
| **endDate** | [**string**] | End date of the donations | (optional) defaults to undefined|


### Return type

**GetStarrings200Response**

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

# **markDeleteStarring**
> markDeleteStarring()

Cancels a donation

### Example

```typescript
import {
    StarsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new StarsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let starringId: string; //ID of the donation (default to undefined)

const { status, data } = await apiInstance.markDeleteStarring(
    accountId,
    starringId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **starringId** | [**string**] | ID of the donation | defaults to undefined|


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

# **patchStarringId**
> Starring patchStarringId()

Update donation\'s state

### Example

```typescript
import {
    StarsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new StarsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let starringId: string; //ID of the donation (default to undefined)
let state: StarringState; //New state of the donation (optional) (default to undefined)

const { status, data } = await apiInstance.patchStarringId(
    accountId,
    starringId,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **starringId** | [**string**] | ID of the donation | defaults to undefined|
| **state** | [**StarringState**] | New state of the donation | (optional) defaults to undefined|


### Return type

**Starring**

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

# **postStarring**
> Starring postStarring()

Create a new stars donations

### Example

```typescript
import {
    StarsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new StarsApi(configuration);

let accountId: string; //ID or CardID of the account (default to undefined)
let amount: number; //Amount of the starring (default to undefined)

const { status, data } = await apiInstance.postStarring(
    accountId,
    amount
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID or CardID of the account | defaults to undefined|
| **amount** | [**number**] | Amount of the starring | defaults to undefined|


### Return type

**Starring**

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

