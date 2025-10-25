# NewTransaction


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**items** | [**Array&lt;NewTransactionItem&gt;**](NewTransactionItem.md) |  | [default to undefined]
**card_pin** | **string** | Pin of the card | [default to undefined]
**is_remote** | **boolean** | True if the transaction was not created from a kiosk | [optional] [default to undefined]

## Example

```typescript
import { NewTransaction } from './api';

const instance: NewTransaction = {
    items,
    card_pin,
    is_remote,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
