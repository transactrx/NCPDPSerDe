package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type CompoundProduct struct {
	Qualifier *string `field:"code=RE"`
	Id        *string `field:"code=TE"`
}

type CompoundIngredientModifierCode struct {
	Code *string `field:"code=2H"`
}

type CompoundIngredient struct {
	Product                  CompoundProduct
	Quantity                 *float64 `field:"code=ED,decimalPlaces=3"`
	DrugCost                 *float64 `field:"code=EE,decimalPlaces=2,overpunch=true"`
	BasisOfCostDetermination *string  `field:"code=UE"`

	ModifierCodeCount *int `field:"code=2G"`
	ModifierCodes     []CompoundIngredientModifierCode
}

type Compound struct {
	SegmentId ncpdp.SegmentId

	DosageFormDescriptionCode   *string `field:"code=EF"`
	DispensingUnitFormIndicator *string `field:"code=EG"`

	IngredientCount *int `field:"code=EC"`
	Ingredients     []CompoundIngredient
}
