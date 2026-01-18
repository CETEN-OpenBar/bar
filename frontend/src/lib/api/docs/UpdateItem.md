# UpdateItem


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**prices** | [**ItemPrices**](ItemPrices.md) |  | [optional] [default to undefined]
**amount_left** | **number** |  | [optional] [default to undefined]
**buy_limit** | **number** |  | [optional] [default to undefined]
**optimal_amount** | **number** |  | [optional] [default to undefined]
**category_id** | **string** |  | [optional] [default to undefined]
**name** | **string** | Name of the current item | [optional] [default to undefined]
**picture** | **string** | Picture of the current item | [optional] [default to undefined]
**promotion** | **number** |  | [optional] [default to undefined]
**promotion_ends_at** | **number** |  | [optional] [default to undefined]
**state** | [**ItemState**](ItemState.md) |  | [optional] [default to undefined]
**available_from** | **number** |  | [optional] [default to undefined]
**available_until** | **number** |  | [optional] [default to undefined]
**is_menu** | **boolean** |  | [optional] [default to undefined]
**menu_items** | [**Array&lt;MenuItem&gt;**](MenuItem.md) |  | [optional] [default to undefined]
**menu_categories** | [**Array&lt;MenuCategory&gt;**](MenuCategory.md) |  | [optional] [default to undefined]
**amount_per_bundle** | **number** |  | [optional] [default to undefined]
**ref_bundle** | **string** |  | [optional] [default to undefined]
**fournisseur** | [**Fournisseur**](Fournisseur.md) |  | [optional] [default to undefined]

## Example

```typescript
import { UpdateItem } from './api';

const instance: UpdateItem = {
    prices,
    amount_left,
    buy_limit,
    optimal_amount,
    category_id,
    name,
    picture,
    promotion,
    promotion_ends_at,
    state,
    available_from,
    available_until,
    is_menu,
    menu_items,
    menu_categories,
    amount_per_bundle,
    ref_bundle,
    fournisseur,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
