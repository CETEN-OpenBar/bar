# Starring


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**account_id** | **string** |  | [default to undefined]
**account_name** | **string** | Name of the account | [default to undefined]
**amount** | **number** |  | [default to undefined]
**state** | [**StarringState**](StarringState.md) |  | [default to undefined]
**issued_at** | **number** |  | [default to undefined]
**issued_by** | **string** |  | [default to undefined]
**issued_by_name** | **string** |  | [default to undefined]
**canceled_by** | **string** |  | [optional] [default to undefined]
**canceled_by_name** | **string** |  | [optional] [default to undefined]
**deleted_at** | **number** |  | [optional] [default to undefined]
**deleted_by** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { Starring } from './api';

const instance: Starring = {
    id,
    account_id,
    account_name,
    amount,
    state,
    issued_at,
    issued_by,
    issued_by_name,
    canceled_by,
    canceled_by_name,
    deleted_at,
    deleted_by,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
