# Transaction


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**items** | [**Array&lt;TransactionItem&gt;**](TransactionItem.md) |  | [default to undefined]
**account_id** | **string** |  | [default to undefined]
**account_name** | **string** | Name of the account | [default to undefined]
**account_nick_name** | **string** | Nickname of the account | [optional] [default to undefined]
**account_google_picture** | **string** | Google profile picture URL of the account | [optional] [default to undefined]
**total_cost** | **number** |  | [default to undefined]
**state** | [**TransactionState**](TransactionState.md) |  | [default to undefined]
**deleted_at** | **number** |  | [optional] [default to undefined]
**deleted_by** | **string** |  | [optional] [default to undefined]
**created_at** | **number** |  | [default to undefined]
**is_remote** | **boolean** | True if the transaction was not created from a kiosk | [optional] [default to undefined]

## Example

```typescript
import { Transaction } from './api';

const instance: Transaction = {
    id,
    items,
    account_id,
    account_name,
    account_nick_name,
    account_google_picture,
    total_cost,
    state,
    deleted_at,
    deleted_by,
    created_at,
    is_remote,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
