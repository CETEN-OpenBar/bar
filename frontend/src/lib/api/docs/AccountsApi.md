# AccountsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminToggleAccountWantsToStaff**](#admintoggleaccountwantstostaff) | **GET** /accounts/{account_id}/toggles/wants_to_staff | |
|[**getAccount**](#getaccount) | **GET** /account | |
|[**getAccountAdmin**](#getaccountadmin) | **GET** /account/admin | |
|[**getAccountId**](#getaccountid) | **GET** /accounts/{account_id} | |
|[**getAccounts**](#getaccounts) | **GET** /accounts | |
|[**importAccounts**](#importaccounts) | **POST** /import/accounts | |
|[**markDeleteAccountId**](#markdeleteaccountid) | **DELETE** /accounts/{account_id} | |
|[**patchAccountId**](#patchaccountid) | **PATCH** /accounts/{account_id} | |
|[**patchAccountPassword**](#patchaccountpassword) | **PATCH** /account/password | |
|[**patchAccountPin**](#patchaccountpin) | **PATCH** /account/pin | |
|[**postAccounts**](#postaccounts) | **POST** /accounts | |
|[**resetAccountPin**](#resetaccountpin) | **POST** /account/{account_id}/reset_pin | |
|[**toggleAccountWantsToStaff**](#toggleaccountwantstostaff) | **GET** /account/toggles/wants_to_staff | |
|[**watchAccount**](#watchaccount) | **GET** /account/watch | |

# **adminToggleAccountWantsToStaff**
> ToggleAccountWantsToStaff200Response adminToggleAccountWantsToStaff()

Toggles the wants_to_staff flag

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let accountId: string; //ID of the account (default to undefined)

const { status, data } = await apiInstance.adminToggleAccountWantsToStaff(
    accountId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|


### Return type

**ToggleAccountWantsToStaff200Response**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Changed account\&#39;s wants_to_staff flag |  -  |
|**401** | Not connected |  -  |
|**404** | Account not found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAccount**
> ConnectCard200Response getAccount()

Get the basic current account\'s information

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

const { status, data } = await apiInstance.getAccount();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ConnectCard200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully got account\&#39;s info |  -  |
|**401** | Not connected |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAccountAdmin**
> GetAccountAdmin200Response getAccountAdmin()

Check if the current account can access the admin panel

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

const { status, data } = await apiInstance.getAccountAdmin();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**GetAccountAdmin200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully got account\&#39;s info |  -  |
|**401** | Not connected |  -  |
|**404** | Account not found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAccountId**
> Account getAccountId()

Get the account\'s information

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let accountId: string; //ID of the account (default to undefined)

const { status, data } = await apiInstance.getAccountId(
    accountId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|


### Return type

**Account**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully got account\&#39;s info |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAccounts**
> GetAccounts200Response getAccounts()

Get all accounts informations

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let page: number; //page to get (optional) (default to 0)
let limit: number; //number of accounts to get (optional) (default to 10)
let search: string; //search string (optional) (default to undefined)
let priceRole: string; //price_role of account (optional) (default to undefined)
let role: string; //role of account (optional) (default to undefined)

const { status, data } = await apiInstance.getAccounts(
    page,
    limit,
    search,
    priceRole,
    role
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | page to get | (optional) defaults to 0|
| **limit** | [**number**] | number of accounts to get | (optional) defaults to 10|
| **search** | [**string**] | search string | (optional) defaults to undefined|
| **priceRole** | [**string**] | price_role of account | (optional) defaults to undefined|
| **role** | [**string**] | role of account | (optional) defaults to undefined|


### Return type

**GetAccounts200Response**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Returns new account\&#39;s info |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **importAccounts**
> ImportAccounts200Response importAccounts()

Import accounts from a CSV file

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let file: File; // (optional) (default to undefined)

const { status, data } = await apiInstance.importAccounts(
    file
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **file** | [**File**] |  | (optional) defaults to undefined|


### Return type

**ImportAccounts200Response**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully imported accounts |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markDeleteAccountId**
> markDeleteAccountId()

Delete account

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let accountId: string; //ID of the account (default to undefined)

const { status, data } = await apiInstance.markDeleteAccountId(
    accountId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|


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
|**404** | Account not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchAccountId**
> Account patchAccountId()

Update account

### Example

```typescript
import {
    AccountsApi,
    Configuration,
    UpdateAccountAdmin
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let updateAccountAdmin: UpdateAccountAdmin; // (optional)

const { status, data } = await apiInstance.patchAccountId(
    accountId,
    updateAccountAdmin
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateAccountAdmin** | **UpdateAccountAdmin**|  | |
| **accountId** | [**string**] | ID of the account | defaults to undefined|


### Return type

**Account**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: application/json
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

# **patchAccountPassword**
> ConnectCard200Response patchAccountPassword()

Update account\'s password

### Example

```typescript
import {
    AccountsApi,
    Configuration,
    PatchAccountPasswordRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let patchAccountPasswordRequest: PatchAccountPasswordRequest; //Passwords (optional)

const { status, data } = await apiInstance.patchAccountPassword(
    patchAccountPasswordRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **patchAccountPasswordRequest** | **PatchAccountPasswordRequest**| Passwords | |


### Return type

**ConnectCard200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully updated account\&#39;s card pin |  -  |
|**400** | Bad request |  -  |
|**401** | Not connected |  -  |
|**404** | Account not found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchAccountPin**
> ConnectCard200Response patchAccountPin()

Update\'s account card pin / id

### Example

```typescript
import {
    AccountsApi,
    Configuration,
    PatchAccountPinRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let patchAccountPinRequest: PatchAccountPinRequest; //Card pin / id (optional)

const { status, data } = await apiInstance.patchAccountPin(
    patchAccountPinRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **patchAccountPinRequest** | **PatchAccountPinRequest**| Card pin / id | |


### Return type

**ConnectCard200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully updated account\&#39;s card pin |  -  |
|**400** | Bad request |  -  |
|**401** | Not connected |  -  |
|**404** | Account not found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **postAccounts**
> Account postAccounts()

Add an account to the database

### Example

```typescript
import {
    AccountsApi,
    Configuration,
    NewAccount
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let newAccount: NewAccount; //Add an account to the database (optional)

const { status, data } = await apiInstance.postAccounts(
    newAccount
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newAccount** | **NewAccount**| Add an account to the database | |


### Return type

**Account**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully got accounts infos |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resetAccountPin**
> resetAccountPin()

Reset the account\'s pin

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

let accountId: string; //ID of the account (default to undefined)

const { status, data } = await apiInstance.resetAccountPin(
    accountId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|


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
|**200** | Successfully reset account\&#39;s pin |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **toggleAccountWantsToStaff**
> ToggleAccountWantsToStaff200Response toggleAccountWantsToStaff()

Toggles the wants_to_staff flag

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

const { status, data } = await apiInstance.toggleAccountWantsToStaff();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ToggleAccountWantsToStaff200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Changed account\&#39;s wants_to_staff flag |  -  |
|**401** | Not connected |  -  |
|**404** | Account not found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **watchAccount**
> ConnectCard200Response watchAccount()

Listen for changes on account

### Example

```typescript
import {
    AccountsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AccountsApi(configuration);

const { status, data } = await apiInstance.watchAccount();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ConnectCard200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully got account\&#39;s info |  -  |
|**401** | Not connected |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

