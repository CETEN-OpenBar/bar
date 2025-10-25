# AuthApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**callback**](#callback) | **GET** /auth/google/callback | |
|[**connectAccount**](#connectaccount) | **GET** /auth/google/begin/{qr_nonce} | |
|[**connectCard**](#connectcard) | **POST** /auth/card | |
|[**connectGoogle**](#connectgoogle) | **GET** /auth/google | |
|[**connectPassword**](#connectpassword) | **POST** /auth/password | |
|[**getAccountQR**](#getaccountqr) | **POST** /account/qr | |
|[**logout**](#logout) | **GET** /logout | |

# **callback**
> callback()

Callback for Google OAuth

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let code: string; //Google OAuth code (default to undefined)
let state: string; //Google OAuth state (default to undefined)

const { status, data } = await apiInstance.callback(
    code,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **code** | [**string**] | Google OAuth code | defaults to undefined|
| **state** | [**string**] | Google OAuth state | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**301** | Redirect to the correct endpoint |  * Location - Redirect to the google oauth page <br>  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **connectAccount**
> connectAccount()

Connect account to Google

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let qrNonce: string; //QR nonce (default to undefined)

const { status, data } = await apiInstance.connectAccount(
    qrNonce
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **qrNonce** | [**string**] | QR nonce | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**301** | Redirect to the google oauth page |  * Location - Redirect to the google oauth page <br>  |
|**400** | Bad request |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **connectCard**
> ConnectCard200Response connectCard()

Connect account to card

### Example

```typescript
import {
    AuthApi,
    Configuration,
    ConnectCardRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let connectCardRequest: ConnectCardRequest; //Card id (optional)

const { status, data } = await apiInstance.connectCard(
    connectCardRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **connectCardRequest** | **ConnectCardRequest**| Card id | |


### Return type

**ConnectCard200Response**

### Authorization

[local_token](../README.md#local_token)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully connected |  -  |
|**400** | Bad request |  -  |
|**401** | Not connected |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **connectGoogle**
> connectGoogle()

Connect account to Google

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let r: string; //Redirect to this url after connecting (default to undefined)

const { status, data } = await apiInstance.connectGoogle(
    r
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **r** | [**string**] | Redirect to this url after connecting | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**301** | Redirect to the google oauth page |  * Location - Redirect to the google oauth page <br>  |
|**400** | Bad request |  -  |
|**401** | Not connected |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **connectPassword**
> ConnectCard200Response connectPassword()

Connect account with password

### Example

```typescript
import {
    AuthApi,
    Configuration,
    ConnectPasswordRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let connectPasswordRequest: ConnectPasswordRequest; //Password (optional)

const { status, data } = await apiInstance.connectPassword(
    connectPasswordRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **connectPasswordRequest** | **ConnectPasswordRequest**| Password | |


### Return type

**ConnectCard200Response**

### Authorization

[local_token](../README.md#local_token)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully connected |  -  |
|**400** | Bad request |  -  |
|**401** | Not connected |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAccountQR**
> string getAccountQR()

Get the QR code to connect account to Google

### Example

```typescript
import {
    AuthApi,
    Configuration,
    GetAccountQRRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let getAccountQRRequest: GetAccountQRRequest; //Card pin (optional)

const { status, data } = await apiInstance.getAccountQR(
    getAccountQRRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **getAccountQRRequest** | **GetAccountQRRequest**| Card pin | |


### Return type

**string**

### Authorization

[not_onboarded](../README.md#not_onboarded), [auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: image/png, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully got QR code |  -  |
|**401** | Not connected |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **logout**
> logout()

Logout

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.logout();
```

### Parameters
This endpoint does not have any parameters.


### Return type

void (empty response body)

### Authorization

[not_onboarded](../README.md#not_onboarded), [auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Successfully logged out |  -  |
|**401** | Not connected |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

