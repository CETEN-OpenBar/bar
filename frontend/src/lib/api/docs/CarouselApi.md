# CarouselApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addCarouselImage**](#addcarouselimage) | **POST** /carousel/images | |
|[**addCarouselText**](#addcarouseltext) | **POST** /carousel/texts | |
|[**getCarouselImage**](#getcarouselimage) | **GET** /carousel/images/{image_id} | |
|[**getCarouselImages**](#getcarouselimages) | **GET** /carousel/images | |
|[**getCarouselTexts**](#getcarouseltexts) | **GET** /carousel/texts | |
|[**markDeleteCarouselImage**](#markdeletecarouselimage) | **DELETE** /carousel/images/{image_id} | |
|[**markDeleteCarouselText**](#markdeletecarouseltext) | **DELETE** /carousel/texts/{text_id} | |

# **addCarouselImage**
> CarouselImage addCarouselImage()

Add a carousel image

### Example

```typescript
import {
    CarouselApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

let image: File; //Image to display (default to undefined)

const { status, data } = await apiInstance.addCarouselImage(
    image
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **image** | [**File**] | Image to display | defaults to undefined|


### Return type

**CarouselImage**

### Authorization

[admin_auth](../README.md#admin_auth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
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

# **addCarouselText**
> CarouselText addCarouselText(carouselTextCreate)

Add a carousel text

### Example

```typescript
import {
    CarouselApi,
    Configuration,
    CarouselTextCreate
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

let carouselTextCreate: CarouselTextCreate; //Carousel text object

const { status, data } = await apiInstance.addCarouselText(
    carouselTextCreate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **carouselTextCreate** | **CarouselTextCreate**| Carousel text object | |


### Return type

**CarouselText**

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

# **getCarouselImage**
> File getCarouselImage()

Get a carousel image

### Example

```typescript
import {
    CarouselApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

let imageId: string; //ID of the image (default to undefined)

const { status, data } = await apiInstance.getCarouselImage(
    imageId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **imageId** | [**string**] | ID of the image | defaults to undefined|


### Return type

**File**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: image/*, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**404** | Image not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCarouselImages**
> Array<CarouselImage> getCarouselImages()

Get carousel images

### Example

```typescript
import {
    CarouselApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

const { status, data } = await apiInstance.getCarouselImages();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<CarouselImage>**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCarouselTexts**
> Array<CarouselText> getCarouselTexts()

Get carousel texts

### Example

```typescript
import {
    CarouselApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

const { status, data } = await apiInstance.getCarouselTexts();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<CarouselText>**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** |  |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markDeleteCarouselImage**
> markDeleteCarouselImage()

Delete a carousel image

### Example

```typescript
import {
    CarouselApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

let imageId: string; //ID of the image (default to undefined)

const { status, data } = await apiInstance.markDeleteCarouselImage(
    imageId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **imageId** | [**string**] | ID of the image | defaults to undefined|


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
|**404** | Image not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markDeleteCarouselText**
> markDeleteCarouselText()

Delete a carousel text

### Example

```typescript
import {
    CarouselApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CarouselApi(configuration);

let textId: string; //ID of the text (default to undefined)

const { status, data } = await apiInstance.markDeleteCarouselText(
    textId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **textId** | [**string**] | ID of the text | defaults to undefined|


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
|**404** | Text not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

