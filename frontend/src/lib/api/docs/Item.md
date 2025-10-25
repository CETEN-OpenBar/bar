# Item


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**prices** | [**ItemPrices**](ItemPrices.md) |  | [default to undefined]
**display_prices** | [**ItemPrices**](ItemPrices.md) |  | [optional] [default to undefined]
**display_price** | **number** |  | [optional] [default to undefined]
**promotion** | **number** |  | [optional] [default to undefined]
**promotion_ends_at** | **number** |  | [optional] [default to undefined]
**amount_left** | **number** |  | [default to undefined]
**optimal_amount** | **number** |  | [default to undefined]
**buy_limit** | **number** |  | [optional] [default to undefined]
**category_id** | **string** |  | [default to undefined]
**name** | **string** | Name of the current item | [default to undefined]
**picture_uri** | **string** | Link to picture of the current item | [default to undefined]
**available_from** | **number** |  | [optional] [default to undefined]
**available_until** | **number** |  | [optional] [default to undefined]
**is_menu** | **boolean** |  | [default to undefined]
**menu_items** | [**Array&lt;MenuItem&gt;**](MenuItem.md) |  | [optional] [default to undefined]
**menu_categories** | [**Array&lt;MenuCategory&gt;**](MenuCategory.md) |  | [optional] [default to undefined]
**state** | [**ItemState**](ItemState.md) |  | [default to undefined]
**last_tva** | **number** |  | [optional] [default to undefined]
**deleted_at** | **number** |  | [optional] [default to undefined]
**deleted_by** | **string** |  | [optional] [default to undefined]
**amount_per_bundle** | **number** |  | [optional] [default to undefined]
**ref_bundle** | **string** | Referal code of the product in the Drive | [optional] [default to undefined]
**fournisseur** | [**Fournisseur**](Fournisseur.md) |  | [optional] [default to undefined]

## Example

```typescript
import { Item } from './api';

const instance: Item = {
    id,
    prices,
    display_prices,
    display_price,
    promotion,
    promotion_ends_at,
    amount_left,
    optimal_amount,
    buy_limit,
    category_id,
    name,
    picture_uri,
    available_from,
    available_until,
    is_menu,
    menu_items,
    menu_categories,
    state,
    last_tva,
    deleted_at,
    deleted_by,
    amount_per_bundle,
    ref_bundle,
    fournisseur,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
