# TransactionsApi

All URIs are relative to *http://localhost:8080*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getAccountTransactions**](#getaccounttransactions) | **GET** /accounts/{account_id}/transactions | |
|[**getCurrentAccountTransactions**](#getcurrentaccounttransactions) | **GET** /account/transactions | |
|[**getTransactionId**](#gettransactionid) | **GET** /accounts/{account_id}/transactions/{transaction_id} | |
|[**getTransactions**](#gettransactions) | **GET** /transactions | |
|[**getTransactionsItems**](#gettransactionsitems) | **GET** /transactions/items | |
|[**markDeleteTransactionId**](#markdeletetransactionid) | **DELETE** /accounts/{account_id}/transactions/{transaction_id} | |
|[**patchTransactionId**](#patchtransactionid) | **PATCH** /accounts/{account_id}/transactions/{transaction_id} | |
|[**patchTransactionItemId**](#patchtransactionitemid) | **PATCH** /accounts/{account_id}/transactions/{transaction_id}/{item_id} | |
|[**postTransactions**](#posttransactions) | **POST** /account/transactions | |

# **getAccountTransactions**
> GetTransactions200Response getAccountTransactions()

Get all transactions

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of transactions per page (optional) (default to undefined)
let state: TransactionState; //Filter by state (optional) (default to undefined)

const { status, data } = await apiInstance.getAccountTransactions(
    accountId,
    page,
    limit,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **state** | [**TransactionState**] | Filter by state | (optional) defaults to undefined|


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

# **getCurrentAccountTransactions**
> GetTransactions200Response getCurrentAccountTransactions()

Get all transactions

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of transactions per page (optional) (default to undefined)
let state: TransactionState; //Filter by state (optional) (default to undefined)

const { status, data } = await apiInstance.getCurrentAccountTransactions(
    page,
    limit,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **state** | [**TransactionState**] | Filter by state | (optional) defaults to undefined|


### Return type

**GetTransactions200Response**

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
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getTransactionId**
> Transaction getTransactionId()

Get transaction

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let transactionId: string; //ID of the transaction (default to undefined)

const { status, data } = await apiInstance.getTransactionId(
    accountId,
    transactionId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **transactionId** | [**string**] | ID of the transaction | defaults to undefined|


### Return type

**Transaction**

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
|**404** | Transaction not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getTransactions**
> GetTransactions200Response getTransactions()

Get all active transactions (orders)

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of transactions per page (optional) (default to undefined)
let state: TransactionState; //Filter by state (optional) (default to undefined)
let hideRemote: boolean; //Hide remote transactions (optional) (default to undefined)
let name: string; //Filter by account name (optional) (default to undefined)

const { status, data } = await apiInstance.getTransactions(
    page,
    limit,
    state,
    hideRemote,
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of transactions per page | (optional) defaults to undefined|
| **state** | [**TransactionState**] | Filter by state | (optional) defaults to undefined|
| **hideRemote** | [**boolean**] | Hide remote transactions | (optional) defaults to undefined|
| **name** | [**string**] | Filter by account name | (optional) defaults to undefined|


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

# **getTransactionsItems**
> Array<TransactionItem> getTransactionsItems()

Get all items in active transactions (ordered items)

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let name: string; //Filter by item name (optional) (default to undefined)

const { status, data } = await apiInstance.getTransactionsItems(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | Filter by item name | (optional) defaults to undefined|


### Return type

**Array<TransactionItem>**

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

# **markDeleteTransactionId**
> HTTPError markDeleteTransactionId()

Delete transaction

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let transactionId: string; //ID of the transaction (default to undefined)

const { status, data } = await apiInstance.markDeleteTransactionId(
    accountId,
    transactionId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **transactionId** | [**string**] | ID of the transaction | defaults to undefined|


### Return type

**HTTPError**

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
|**404** | Transaction not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchTransactionId**
> HTTPError patchTransactionId()

Update transaction\'s state

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let transactionId: string; //ID of the transaction (default to undefined)
let state: TransactionState; //New state of the transaction (default to undefined)

const { status, data } = await apiInstance.patchTransactionId(
    accountId,
    transactionId,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **transactionId** | [**string**] | ID of the transaction | defaults to undefined|
| **state** | [**TransactionState**] | New state of the transaction | defaults to undefined|


### Return type

**HTTPError**

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
|**404** | Account or transaction not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **patchTransactionItemId**
> HTTPError patchTransactionItemId()

Update transaction\'s item state

### Example

```typescript
import {
    TransactionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let accountId: string; //ID of the account (default to undefined)
let transactionId: string; //ID of the transaction (default to undefined)
let itemId: string; //ID of the item (default to undefined)
let state: TransactionItemState; //New state of the item (optional) (default to undefined)
let amount: number; //New amount of the item (optional) (default to undefined)
let alreadyDone: number; //Update item\'s already done (optional) (default to undefined)

const { status, data } = await apiInstance.patchTransactionItemId(
    accountId,
    transactionId,
    itemId,
    state,
    amount,
    alreadyDone
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **accountId** | [**string**] | ID of the account | defaults to undefined|
| **transactionId** | [**string**] | ID of the transaction | defaults to undefined|
| **itemId** | [**string**] | ID of the item | defaults to undefined|
| **state** | **TransactionItemState** | New state of the item | (optional) defaults to undefined|
| **amount** | [**number**] | New amount of the item | (optional) defaults to undefined|
| **alreadyDone** | [**number**] | Update item\&#39;s already done | (optional) defaults to undefined|


### Return type

**HTTPError**

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
|**404** | Account, transaction or item not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **postTransactions**
> Transaction postTransactions()

Create a new transaction

### Example

```typescript
import {
    TransactionsApi,
    Configuration,
    NewTransaction
} from './api';

const configuration = new Configuration();
const apiInstance = new TransactionsApi(configuration);

let newTransaction: NewTransaction; // (optional)

const { status, data } = await apiInstance.postTransactions(
    newTransaction
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **newTransaction** | **NewTransaction**|  | |


### Return type

**Transaction**

### Authorization

[auth](../README.md#auth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** |  |  -  |
|**401** | Not connected |  -  |
|**403** | Forbidden |  -  |
|**404** | Account not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

