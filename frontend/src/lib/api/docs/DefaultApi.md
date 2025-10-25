# DefaultApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getAccountQRWebsocket**](#getaccountqrwebsocket) | **GET** /account/qr | |
|[**getBorneAuthQRWebsocket**](#getborneauthqrwebsocket) | **GET** /auth/qr | |
|[**postBorneAuthQR**](#postborneauthqr) | **POST** /auth/qr | |

# **getAccountQRWebsocket**
> getAccountQRWebsocket()

Websocket to listen for scan & callback (for cool animations)

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getAccountQRWebsocket();
```

### Parameters
This endpoint does not have any parameters.


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**101** | Switching Protocols |  * Connection -  <br>  * Upgrade -  <br>  * Sec-WebSocket-Accept -  <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getBorneAuthQRWebsocket**
> getBorneAuthQRWebsocket()

Websocket to listen for scan & callback (for cool animations)

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getBorneAuthQRWebsocket();
```

### Parameters
This endpoint does not have any parameters.


### Return type

void (empty response body)

### Authorization

[local_token](../README.md#local_token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**101** | Switching Protocols |  * Connection -  <br>  * Upgrade -  <br>  * Sec-WebSocket-Accept -  <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **postBorneAuthQR**
> ConnectCard200Response postBorneAuthQR()

Validate the connection to connect

### Example

```typescript
import {
    DefaultApi,
    Configuration,
    PostBorneAuthQRRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let postBorneAuthQRRequest: PostBorneAuthQRRequest; //Nonce (optional)

const { status, data } = await apiInstance.postBorneAuthQR(
    postBorneAuthQRRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **postBorneAuthQRRequest** | **PostBorneAuthQRRequest**| Nonce | |


### Return type

**ConnectCard200Response**

### Authorization

No authorization required

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

