# Account


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**first_name** | **string** |  | [default to undefined]
**last_name** | **string** |  | [default to undefined]
**nickname** | **string** |  | [optional] [default to undefined]
**google_id** | **string** |  | [optional] [default to undefined]
**google_picture** | **string** |  | [optional] [default to undefined]
**email_address** | **string** |  | [default to undefined]
**balance** | **number** |  | [default to undefined]
**points** | **number** |  | [default to undefined]
**card_id** | **string** |  | [optional] [default to undefined]
**card_pin** | **string** |  | [default to undefined]
**password** | **string** |  | [optional] [default to undefined]
**role** | [**AccountRole**](AccountRole.md) |  | [default to undefined]
**price_role** | [**AccountPriceRole**](AccountPriceRole.md) |  | [default to undefined]
**restrictions** | [**Array&lt;AccountRestrictions&gt;**](AccountRestrictions.md) |  | [default to undefined]
**state** | [**AccountState**](AccountState.md) |  | [default to undefined]
**deleted_at** | **number** |  | [optional] [default to undefined]
**deleted_by** | **string** |  | [optional] [default to undefined]
**wants_to_staff** | **boolean** |  | [default to undefined]

## Example

```typescript
import { Account } from './api';

const instance: Account = {
    id,
    first_name,
    last_name,
    nickname,
    google_id,
    google_picture,
    email_address,
    balance,
    points,
    card_id,
    card_pin,
    password,
    role,
    price_role,
    restrictions,
    state,
    deleted_at,
    deleted_by,
    wants_to_staff,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
