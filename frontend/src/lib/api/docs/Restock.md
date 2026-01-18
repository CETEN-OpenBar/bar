# Restock


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**items** | [**Array&lt;RestockItem&gt;**](RestockItem.md) |  | [default to undefined]
**total_cost_ht** | **number** |  | [default to undefined]
**total_cost_ttc** | **number** |  | [default to undefined]
**driver_id** | **string** |  | [optional] [default to undefined]
**driver_name** | **string** | Name of the driver | [optional] [default to undefined]
**id** | **string** |  | [default to undefined]
**type** | [**RestockType**](RestockType.md) |  | [default to undefined]
**created_at** | **number** |  | [default to undefined]
**created_by** | **string** |  | [default to undefined]
**created_by_name** | **string** |  | [default to undefined]
**deleted_at** | **number** |  | [optional] [default to undefined]
**deleted_by** | **string** |  | [optional] [default to undefined]
**state** | [**RestockState**](RestockState.md) |  | [default to undefined]

## Example

```typescript
import { Restock } from './api';

const instance: Restock = {
    items,
    total_cost_ht,
    total_cost_ttc,
    driver_id,
    driver_name,
    id,
    type,
    created_at,
    created_by,
    created_by_name,
    deleted_at,
    deleted_by,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
