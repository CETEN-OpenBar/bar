# UpdateAccountAdmin


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**first_name** | **string** |  | [optional] [default to undefined]
**last_name** | **string** |  | [optional] [default to undefined]
**nickname** | **string** |  | [optional] [default to undefined]
**email_address** | **string** |  | [optional] [default to undefined]
**card_id** | **string** |  | [optional] [default to undefined]
**role** | [**AccountRole**](AccountRole.md) |  | [optional] [default to undefined]
**price_role** | [**AccountPriceRole**](AccountPriceRole.md) |  | [optional] [default to undefined]
**restrictions** | [**Array&lt;AccountRestrictions&gt;**](AccountRestrictions.md) |  | [optional] [default to undefined]
**state** | [**AccountState**](AccountState.md) |  | [optional] [default to undefined]

## Example

```typescript
import { UpdateAccountAdmin } from './api';

const instance: UpdateAccountAdmin = {
    first_name,
    last_name,
    nickname,
    email_address,
    card_id,
    role,
    price_role,
    restrictions,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
