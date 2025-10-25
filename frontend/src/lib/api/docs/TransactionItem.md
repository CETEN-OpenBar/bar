# TransactionItem


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**item_id** | **string** |  | [default to undefined]
**picture_uri** | **string** | Link to picture of the current item | [default to undefined]
**item_name** | **string** | Name of the current item | [default to undefined]
**item_amount** | **number** |  | [default to undefined]
**item_already_done** | **number** |  | [default to undefined]
**unit_cost** | **number** |  | [default to undefined]
**is_menu** | **boolean** |  | [default to undefined]
**menu_items** | [**Array&lt;MenuItem&gt;**](MenuItem.md) |  | [optional] [default to undefined]
**menu_categories** | [**Array&lt;MenuCategory&gt;**](MenuCategory.md) |  | [optional] [default to undefined]
**picked_categories_items** | [**Array&lt;TransactionItem&gt;**](TransactionItem.md) |  | [optional] [default to undefined]
**total_cost** | **number** |  | [default to undefined]
**state** | [**TransactionItemState**](TransactionItemState.md) |  | [default to undefined]

## Example

```typescript
import { TransactionItem } from './api';

const instance: TransactionItem = {
    item_id,
    picture_uri,
    item_name,
    item_amount,
    item_already_done,
    unit_cost,
    is_menu,
    menu_items,
    menu_categories,
    picked_categories_items,
    total_cost,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
