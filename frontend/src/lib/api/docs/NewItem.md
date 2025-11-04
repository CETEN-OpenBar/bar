# NewItem


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**prices** | [**ItemPrices**](ItemPrices.md) |  | [default to undefined]
**promotion** | **number** |  | [optional] [default to undefined]
**promotion_ends_at** | **number** |  | [optional] [default to undefined]
**amount_left** | **number** |  | [default to undefined]
**optimal_amount** | **number** |  | [default to undefined]
**buy_limit** | **number** |  | [optional] [default to undefined]
**available_from** | **number** |  | [optional] [default to undefined]
**available_until** | **number** |  | [optional] [default to undefined]
**is_menu** | **boolean** |  | [optional] [default to undefined]
**menu_items** | [**Array&lt;MenuItem&gt;**](MenuItem.md) |  | [optional] [default to undefined]
**menu_categories** | [**Array&lt;MenuCategory&gt;**](MenuCategory.md) |  | [optional] [default to undefined]
**name** | **string** | Name of the current item | [default to undefined]
**picture** | **string** | Picture of the current item | [default to undefined]
**state** | [**ItemState**](ItemState.md) |  | [default to undefined]

## Example

```typescript
import { NewItem } from './api';

const instance: NewItem = {
    prices,
    promotion,
    promotion_ends_at,
    amount_left,
    optimal_amount,
    buy_limit,
    available_from,
    available_until,
    is_menu,
    menu_items,
    menu_categories,
    name,
    picture,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
