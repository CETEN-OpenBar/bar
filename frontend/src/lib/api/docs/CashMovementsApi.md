# CashMovementsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createCashMovement**](#createcashmovement) | **POST** /cash_movements | |
|[**deleteCashMovement**](#deletecashmovement) | **DELETE** /cash_movements/{cash_movement_id} | |
|[**getCashMovements**](#getcashmovements) | **GET** /cash_movements | |

# **createCashMovement**
> CashMovement createCashMovement(newCashMovement)

Create a cash movement

### Example

```typescript
import {
    CashMovementsApi,
    Configuration,
    NewCashMovement
} from './api';

const configuration = new Configuration();
const apiInstance = new CashMovementsApi(configuration);

let newCashMovement: NewCashMovement; //Cash movement to create

const { status, data } = await apiInstance.createCashMovement(
    newCashMovement
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newCashMovement** | **NewCashMovement**| Cash movement to create | |


### Return type

**CashMovement**

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
|**409** | Cash movement already exists |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCashMovement**
> deleteCashMovement()

Delete a cash movement

### Example

```typescript
import {
    CashMovementsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CashMovementsApi(configuration);

let cashMovementId: string; //ID of the cash movement (default to undefined)

const { status, data } = await apiInstance.deleteCashMovement(
    cashMovementId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **cashMovementId** | [**string**] | ID of the cash movement | defaults to undefined|


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
|**404** | Cash movement not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCashMovements**
> GetCashMovements200Response getCashMovements()

Get cash movements

### Example

```typescript
import {
    CashMovementsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CashMovementsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of cash movements per page (optional) (default to undefined)
let search: string; //search string (optional) (default to undefined)

const { status, data } = await apiInstance.getCashMovements(
    page,
    limit,
    search
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of cash movements per page | (optional) defaults to undefined|
| **search** | [**string**] | search string | (optional) defaults to undefined|


### Return type

**GetCashMovements200Response**

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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

