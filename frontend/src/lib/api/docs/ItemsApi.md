# ItemsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getAllIncoherentItems**](#getallincoherentitems) | **GET** /items/incoherent | |
|[**getAllItems**](#getallitems) | **GET** /items | |
|[**getCategoryItems**](#getcategoryitems) | **GET** /categories/{category_id}/items | |
|[**getItem**](#getitem) | **GET** /item/{item_id} | |
|[**getItemPicture**](#getitempicture) | **GET** /categories/{category_id}/items/{item_id}/picture | |
|[**markDeleteItem**](#markdeleteitem) | **DELETE** /categories/{category_id}/items/{item_id} | |
|[**patchItem**](#patchitem) | **PATCH** /categories/{category_id}/items/{item_id} | |
|[**postItem**](#postitem) | **POST** /categories/{category_id}/items | |

# **getAllIncoherentItems**
> GetAllItems200Response getAllIncoherentItems()

(admin) Get all incoherent items with filters and pagination

### Example

```typescript
import {
    ItemsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of items per page (optional) (default to undefined)
let state: ItemState; //Filter by state (optional) (default to undefined)
let categoryId: string; //Filter by category (optional) (default to undefined)
let name: string; //Filter by name (optional) (default to undefined)

const { status, data } = await apiInstance.getAllIncoherentItems(
    page,
    limit,
    state,
    categoryId,
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of items per page | (optional) defaults to undefined|
| **state** | [**ItemState**] | Filter by state | (optional) defaults to undefined|
| **categoryId** | [**string**] | Filter by category | (optional) defaults to undefined|
| **name** | [**string**] | Filter by name | (optional) defaults to undefined|


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
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAllItems**
> GetAllItems200Response getAllItems()

(admin) Get all items with filters and pagination

### Example

```typescript
import {
    ItemsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of items per page (optional) (default to undefined)
let state: ItemState; //Filter by state (optional) (default to undefined)
let categoryId: string; //Filter by category (optional) (default to undefined)
let name: string; //Filter by name (optional) (default to undefined)
let fournisseur: Fournisseur; //Filter by fournisseur (optional) (default to undefined)
let refBundle: string; //Filter by reference (optional) (default to undefined)

const { status, data } = await apiInstance.getAllItems(
    page,
    limit,
    state,
    categoryId,
    name,
    fournisseur,
    refBundle
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of items per page | (optional) defaults to undefined|
| **state** | [**ItemState**] | Filter by state | (optional) defaults to undefined|
| **categoryId** | [**string**] | Filter by category | (optional) defaults to undefined|
| **name** | [**string**] | Filter by name | (optional) defaults to undefined|
| **fournisseur** | **Fournisseur** | Filter by fournisseur | (optional) defaults to undefined|
| **refBundle** | [**string**] | Filter by reference | (optional) defaults to undefined|


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
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCategoryItems**
> GetAllItems200Response getCategoryItems()

Get all items of a category

### Example

```typescript
import {
    ItemsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let categoryId: string; //ID of the category (default to undefined)
let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of items per page (optional) (default to undefined)
let state: ItemState; //Filter by state (optional) (default to undefined)

const { status, data } = await apiInstance.getCategoryItems(
    categoryId,
    page,
    limit,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of items per page | (optional) defaults to undefined|
| **state** | [**ItemState**] | Filter by state | (optional) defaults to undefined|


### Return type

**GetAllItems200Response**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**404** | Category not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getItem**
> Item getItem()

Get an item

### Example

```typescript
import {
    ItemsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let itemId: string; //ID of the item (default to undefined)

const { status, data } = await apiInstance.getItem(
    itemId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **itemId** | [**string**] | ID of the item | defaults to undefined|


### Return type

**Item**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**404** | Item not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getItemPicture**
> File getItemPicture()

Get an item picture

### Example

```typescript
import {
    ItemsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let categoryId: string; //ID of the category (default to undefined)
let itemId: string; //ID of the item (default to undefined)

const { status, data } = await apiInstance.getItemPicture(
    categoryId,
    itemId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|
| **itemId** | [**string**] | ID of the item | defaults to undefined|


### Return type

**File**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: image/png, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**404** | Item not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markDeleteItem**
> markDeleteItem()

Delete an item

### Example

```typescript
import {
    ItemsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let categoryId: string; //ID of the category (default to undefined)
let itemId: string; //ID of the item (default to undefined)

const { status, data } = await apiInstance.markDeleteItem(
    categoryId,
    itemId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|
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
|**204** |  |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Item not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchItem**
> Item patchItem(updateItem)

Update an item

### Example

```typescript
import {
    ItemsApi,
    Configuration,
    UpdateItem
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let categoryId: string; //ID of the category (default to undefined)
let itemId: string; //ID of the item (default to undefined)
let updateItem: UpdateItem; //Item object

const { status, data } = await apiInstance.patchItem(
    categoryId,
    itemId,
    updateItem
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateItem** | **UpdateItem**| Item object | |
| **categoryId** | [**string**] | ID of the category | defaults to undefined|
| **itemId** | [**string**] | ID of the item | defaults to undefined|


### Return type

**Item**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Item not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **postItem**
> Item postItem(newItem)

Create a new item

### Example

```typescript
import {
    ItemsApi,
    Configuration,
    NewItem
} from './api';

const configuration = new Configuration();
const apiInstance = new ItemsApi(configuration);

let categoryId: string; //ID of the category (default to undefined)
let newItem: NewItem; //Item object

const { status, data } = await apiInstance.postItem(
    categoryId,
    newItem
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newItem** | **NewItem**| Item object | |
| **categoryId** | [**string**] | ID of the category | defaults to undefined|


### Return type

**Item**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** |  |  -  |
|**400** | Bad request |  -  |
|**401** | Not authorized |  -  |
|**403** | Forbidden |  -  |
|**409** | Item already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

