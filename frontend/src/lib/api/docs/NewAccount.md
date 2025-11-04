# NewAccount


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**first_name** | **string** |  | [default to undefined]
**last_name** | **string** |  | [default to undefined]
**email_address** | **string** |  | [default to undefined]
**card_id** | **string** |  | [optional] [default to undefined]
**balance** | **number** |  | [default to undefined]
**role** | [**AccountRole**](AccountRole.md) |  | [default to undefined]
**price_role** | [**AccountPriceRole**](AccountPriceRole.md) |  | [optional] [default to undefined]

## Example

```typescript
import { NewAccount } from './api';

const instance: NewAccount = {
    first_name,
    last_name,
    email_address,
    card_id,
    balance,
    role,
    price_role,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
