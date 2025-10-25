# NewRestock


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**items** | [**Array&lt;NewRestockItem&gt;**](NewRestockItem.md) |  | [default to undefined]
**total_cost_ht** | **number** |  | [default to undefined]
**total_cost_ttc** | **number** |  | [default to undefined]
**driver_id** | **string** |  | [optional] [default to undefined]
**type** | [**RestockType**](RestockType.md) |  | [default to undefined]
**state** | [**RestockState**](RestockState.md) |  | [default to undefined]

## Example

```typescript
import { NewRestock } from './api';

const instance: NewRestock = {
    items,
    total_cost_ht,
    total_cost_ttc,
    driver_id,
    type,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
