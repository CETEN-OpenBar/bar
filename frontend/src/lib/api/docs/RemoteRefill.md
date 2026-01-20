# RemoteRefill

A remote (HelloAsso) refill

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**state** | [**RemoteRefillState**](RemoteRefillState.md) |  | [default to undefined]
**checkout_intent_id** | **number** | HelloAsso checkout id | [optional] [default to undefined]
**order_id** | **number** | HelloAsso order id, if the transaction suceeded | [optional] [default to undefined]
**account_id** | **string** |  | [default to undefined]
**account_name** | **string** | Name of the account | [default to undefined]
**account_google_picture** | **string** | Google profile picture URL of the account | [optional] [default to undefined]
**amount** | **number** |  | [default to undefined]
**created_at** | **number** |  | [default to undefined]
**refill_id** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { RemoteRefill } from './api';

const instance: RemoteRefill = {
    id,
    state,
    checkout_intent_id,
    order_id,
    account_id,
    account_name,
    account_google_picture,
    amount,
    created_at,
    refill_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
