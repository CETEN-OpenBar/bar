# CategoriesApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getCategories**](#getcategories) | **GET** /categories | |
|[**getCategory**](#getcategory) | **GET** /categories/{category_id} | |
|[**getCategoryPicture**](#getcategorypicture) | **GET** /categories/{category_id}/picture | |
|[**markDeleteCategory**](#markdeletecategory) | **DELETE** /categories/{category_id} | |
|[**patchCategory**](#patchcategory) | **PATCH** /categories/{category_id} | |
|[**postCategory**](#postcategory) | **POST** /categories | |

# **getCategories**
> Array<Category> getCategories()

Get all categories

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let hidden: boolean; //Show hidden categories (admin only) (optional) (default to undefined)

const { status, data } = await apiInstance.getCategories(
    hidden
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **hidden** | [**boolean**] | Show hidden categories (admin only) | (optional) defaults to undefined|


### Return type

**Array<Category>**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCategory**
> Category getCategory()

Get a category

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let categoryId: string; //ID of the category (default to undefined)

const { status, data } = await apiInstance.getCategory(
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|


### Return type

**Category**

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

# **getCategoryPicture**
> File getCategoryPicture()

Get a category picture

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let categoryId: string; //ID of the category (default to undefined)

const { status, data } = await apiInstance.getCategoryPicture(
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] | ID of the category | defaults to undefined|


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
|**404** | Category not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markDeleteCategory**
> markDeleteCategory()

Delete a category

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let categoryId: string; //ID of the category (default to undefined)

const { status, data } = await apiInstance.markDeleteCategory(
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
|**204** |  |  -  |
|**404** | Category not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchCategory**
> Category patchCategory(updateCategory)

Update a category

### Example

```typescript
import {
    CategoriesApi,
    Configuration,
    UpdateCategory
} from './api';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let categoryId: string; //ID of the category (default to undefined)
let updateCategory: UpdateCategory; //Category object

const { status, data } = await apiInstance.patchCategory(
    categoryId,
    updateCategory
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateCategory** | **UpdateCategory**| Category object | |
| **categoryId** | [**string**] | ID of the category | defaults to undefined|


### Return type

**Category**

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
|**404** | Category not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **postCategory**
> Category postCategory(newCategory)

Create a new category

### Example

```typescript
import {
    CategoriesApi,
    Configuration,
    NewCategory
} from './api';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let newCategory: NewCategory; //Category object

const { status, data } = await apiInstance.postCategory(
    newCategory
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newCategory** | **NewCategory**| Category object | |


### Return type

**Category**

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
|**409** | Category already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

