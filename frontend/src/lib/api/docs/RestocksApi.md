# RestocksApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createRestock**](#createrestock) | **POST** /restocks | |
|[**deleteRestock**](#deleterestock) | **DELETE** /restocks/{restock_id} | |
|[**getRestocks**](#getrestocks) | **GET** /restocks | |
|[**updateRestock**](#updaterestock) | **PATCH** /restocks/{restock_id} | |

# **createRestock**
> Restock createRestock(newRestock)

Create a restock

### Example

```typescript
import {
    RestocksApi,
    Configuration,
    NewRestock
} from './api';

const configuration = new Configuration();
const apiInstance = new RestocksApi(configuration);

let newRestock: NewRestock; //Restock to create

const { status, data } = await apiInstance.createRestock(
    newRestock
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newRestock** | **NewRestock**| Restock to create | |


### Return type

**Restock**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Bad request |  -  |
|**401** | Not authenticated |  -  |
|**403** | Forbidden |  -  |
|**409** | Restock already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRestock**
> deleteRestock()

Delete a restock

### Example

```typescript
import {
    RestocksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RestocksApi(configuration);

let restockId: string; //ID of the restock (default to undefined)

const { status, data } = await apiInstance.deleteRestock(
    restockId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **restockId** | [**string**] | ID of the restock | defaults to undefined|


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
|**204** | Deleted |  -  |
|**401** | Not authenticated |  -  |
|**403** | Forbidden |  -  |
|**404** | Restock not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRestocks**
> GetRestocks200Response getRestocks()

Get restocks

### Example

```typescript
import {
    RestocksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RestocksApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of restocks per page (optional) (default to undefined)
let state: RestockState; //search state (optional) (default to undefined)

const { status, data } = await apiInstance.getRestocks(
    page,
    limit,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of restocks per page | (optional) defaults to undefined|
| **state** | **RestockState** | search state | (optional) defaults to undefined|


### Return type

**GetRestocks200Response**

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
|**401** | Not authenticated |  -  |
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRestock**
> updateRestock(newRestock)

Update a restock

### Example

```typescript
import {
    RestocksApi,
    Configuration,
    NewRestock
} from './api';

const configuration = new Configuration();
const apiInstance = new RestocksApi(configuration);

let restockId: string; //ID of the restock (default to undefined)
let newRestock: NewRestock; //Restock to update

const { status, data } = await apiInstance.updateRestock(
    restockId,
    newRestock
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newRestock** | **NewRestock**| Restock to update | |
| **restockId** | [**string**] | ID of the restock | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**400** | Finished restock |  -  |
|**401** | Not authenticated |  -  |
|**403** | Forbidden |  -  |
|**404** | Restock not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

