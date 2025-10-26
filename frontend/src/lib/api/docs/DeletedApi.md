# DeletedApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**deleteAccount**](#deleteaccount) | **DELETE** /deleted/accounts/{account_id} | |
|[**deleteCarouselImage**](#deletecarouselimage) | **DELETE** /deleted/carousel/images/{image_id} | |
|[**deleteCarouselText**](#deletecarouseltext) | **DELETE** /deleted/carousel/texts/{text_id} | |
|[**deleteCategory**](#deletecategory) | **DELETE** /deleted/categories/{category_id} | |
|[**deleteItem**](#deleteitem) | **DELETE** /deleted/items/{item_id} | |
|[**deleteRefill**](#deleterefill) | **DELETE** /deleted/refills/{refill_id} | |
|[**deleteStarring**](#deletestarring) | **DELETE** /deleted/stars/{starring_id} | |
|[**deleteTransaction**](#deletetransaction) | **DELETE** /deleted/transactions/{transaction_id} | |
|[**getDeletedAccounts**](#getdeletedaccounts) | **GET** /deleted/accounts | |
|[**getDeletedCarouselImages**](#getdeletedcarouselimages) | **GET** /deleted/carousel/images | |
|[**getDeletedCarouselTexts**](#getdeletedcarouseltexts) | **GET** /deleted/carousel/texts | |
|[**getDeletedCategories**](#getdeletedcategories) | **GET** /deleted/categories | |
|[**getDeletedItems**](#getdeleteditems) | **GET** /deleted/items | |
|[**getDeletedRefills**](#getdeletedrefills) | **GET** /deleted/refills | |
|[**getDeletedStarring**](#getdeletedstarring) | **GET** /deleted/stars | |
|[**getDeletedTransactions**](#getdeletedtransactions) | **GET** /deleted/transactions | |
|[**restoreDeletedAccount**](#restoredeletedaccount) | **PATCH** /deleted/accounts/{account_id} | |
|[**restoreDeletedCarouselImage**](#restoredeletedcarouselimage) | **PATCH** /deleted/carousel/images/{image_id} | |
|[**restoreDeletedCarouselText**](#restoredeletedcarouseltext) | **PATCH** /deleted/carousel/texts/{text_id} | |
|[**restoreDeletedCategory**](#restoredeletedcategory) | **PATCH** /deleted/categories/{category_id} | |
|[**restoreDeletedItem**](#restoredeleteditem) | **PATCH** /deleted/items/{item_id} | |
|[**restoreDeletedRefill**](#restoredeletedrefill) | **PATCH** /deleted/refills/{refill_id} | |
|[**restoreDeletedStarring**](#restoredeletedstarring) | **PATCH** /deleted/stars/{starring_id} | |
|[**restoreDeletedTransaction**](#restoredeletedtransaction) | **PATCH** /deleted/transactions/{transaction_id} | |

# **deleteAccount**
> deleteAccount()

Permanently deletes an account (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let accountId: string; //ID of the account (default to undefined)

const { status, data } = await apiInstance.deleteAccount(
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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**409** | Account already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCarouselImage**
> deleteCarouselImage()

Permanently deletes a carousel image (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let imageId: string; //ID of the carousel image (default to undefined)

const { status, data } = await apiInstance.deleteCarouselImage(
    imageId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **imageId** | [**string**] | ID of the carousel image | defaults to undefined|


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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Carousel image not found |  -  |
|**409** | Carousel image already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCarouselText**
> deleteCarouselText()

Permanently deletes a carousel text (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let textId: string; //ID of the carousel text (default to undefined)

const { status, data } = await apiInstance.deleteCarouselText(
    textId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **textId** | [**string**] | ID of the carousel text | defaults to undefined|


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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Carousel text not found |  -  |
|**409** | Carousel text already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCategory**
> deleteCategory()

Permanently deletes a category (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let categoryId: string; //ID of the category (default to undefined)

const { status, data } = await apiInstance.deleteCategory(
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|


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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Category not found |  -  |
|**409** | Category already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteItem**
> deleteItem()

Permanently deletes an item (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let itemId: string; //ID of the item (default to undefined)

const { status, data } = await apiInstance.deleteItem(
    itemId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **itemId** | [**string**] | ID of the item | defaults to undefined|


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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Item not found |  -  |
|**409** | Item already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRefill**
> deleteRefill()

Permanently deletes a refill (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let refillId: string; //ID of the refill (default to undefined)

const { status, data } = await apiInstance.deleteRefill(
    refillId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Refill not found |  -  |
|**409** | Refill already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteStarring**
> deleteStarring()

Permanently deletes a starring (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let starringId: string; //ID of the starring (default to undefined)

const { status, data } = await apiInstance.deleteStarring(
    starringId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **starringId** | [**string**] | ID of the starring | defaults to undefined|


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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Starring not found |  -  |
|**409** | Starring already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteTransaction**
> deleteTransaction()

Permanently deletes a transaction (SUPERADMIN)

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let transactionId: string; //ID of the transaction (default to undefined)

const { status, data } = await apiInstance.deleteTransaction(
    transactionId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **transactionId** | [**string**] | ID of the transaction | defaults to undefined|


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
|**204** | Success |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Transaction not found |  -  |
|**409** | Transaction already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getDeletedAccounts**
> GetDeletedAccounts200Response getDeletedAccounts()

Get deleted accounts

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)
let search: string; //search string (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedAccounts(
    page,
    limit,
    search
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|
| **search** | [**string**] | search string | (optional) defaults to undefined|


### Return type

**GetDeletedAccounts200Response**

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

# **getDeletedCarouselImages**
> GetDeletedCarouselImages200Response getDeletedCarouselImages()

Get deleted carousel images

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedCarouselImages(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|


### Return type

**GetDeletedCarouselImages200Response**

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

# **getDeletedCarouselTexts**
> GetDeletedCarouselTexts200Response getDeletedCarouselTexts()

Get deleted carousel texts

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedCarouselTexts(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|


### Return type

**GetDeletedCarouselTexts200Response**

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

# **getDeletedCategories**
> GetDeletedCategories200Response getDeletedCategories()

Get deleted categories

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of categories per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedCategories(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of categories per page | (optional) defaults to undefined|


### Return type

**GetDeletedCategories200Response**

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

# **getDeletedItems**
> GetAllItems200Response getDeletedItems()

Get deleted items

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedItems(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|


### Return type

**GetAllItems200Response**

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

# **getDeletedRefills**
> GetRefills200Response getDeletedRefills()

Get deleted refills

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedRefills(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|


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

# **getDeletedStarring**
> GetDeletedStarring200Response getDeletedStarring()

Get deleted starrings

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedStarring(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|


### Return type

**GetDeletedStarring200Response**

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

# **getDeletedTransactions**
> GetTransactions200Response getDeletedTransactions()

Get deleted transactions

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of accounts per page (optional) (default to undefined)

const { status, data } = await apiInstance.getDeletedTransactions(
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of accounts per page | (optional) defaults to undefined|


### Return type

**GetTransactions200Response**

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

# **restoreDeletedAccount**
> restoreDeletedAccount()

Restore a deleted account

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let accountId: string; //ID of the account (default to undefined)

const { status, data } = await apiInstance.restoreDeletedAccount(
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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**409** | Account already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedCarouselImage**
> restoreDeletedCarouselImage()

Restore a deleted carousel image

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let imageId: string; //ID of the carousel image (default to undefined)

const { status, data } = await apiInstance.restoreDeletedCarouselImage(
    imageId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **imageId** | [**string**] | ID of the carousel image | defaults to undefined|


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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Carousel image not found |  -  |
|**409** | Carousel image already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedCarouselText**
> restoreDeletedCarouselText()

Restore a deleted carousel text

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let textId: string; //ID of the carousel text (default to undefined)

const { status, data } = await apiInstance.restoreDeletedCarouselText(
    textId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **textId** | [**string**] | ID of the carousel text | defaults to undefined|


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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Carousel text not found |  -  |
|**409** | Carousel text already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedCategory**
> restoreDeletedCategory()

Restore a deleted category

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let categoryId: string; //ID of the category (default to undefined)

const { status, data } = await apiInstance.restoreDeletedCategory(
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|


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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Category not found |  -  |
|**409** | Category already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedItem**
> restoreDeletedItem()

Restore a deleted item

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let itemId: string; //ID of the item (default to undefined)

const { status, data } = await apiInstance.restoreDeletedItem(
    itemId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **itemId** | [**string**] | ID of the item | defaults to undefined|


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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Item not found |  -  |
|**409** | Item already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedRefill**
> restoreDeletedRefill()

Restore a deleted refill

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let refillId: string; //ID of the refill (default to undefined)

const { status, data } = await apiInstance.restoreDeletedRefill(
    refillId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Refill not found |  -  |
|**409** | Refill already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedStarring**
> restoreDeletedStarring()

Restore a deleted starring

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let starringId: string; //ID of the starring (default to undefined)

const { status, data } = await apiInstance.restoreDeletedStarring(
    starringId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **starringId** | [**string**] | ID of the starring | defaults to undefined|


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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Starring not found |  -  |
|**409** | Starring already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restoreDeletedTransaction**
> restoreDeletedTransaction()

Restore a deleted transaction

### Example

```typescript
import {
    DeletedApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeletedApi(configuration);

let transactionId: string; //ID of the transaction (default to undefined)

const { status, data } = await apiInstance.restoreDeletedTransaction(
    transactionId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **transactionId** | [**string**] | ID of the transaction | defaults to undefined|


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
|**204** | Success |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Transaction not found |  -  |
|**409** | Transaction already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

