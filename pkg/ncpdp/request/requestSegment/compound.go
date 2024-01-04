package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Compound struct {
	SegmentId ncpdp.SegmentId

	DosageFormDescriptionCode   *string `field:"code=EF,order=2"`
	DispensingUnitFormIndicator *string `field:"code=EG,order=3"`

	IngredientCount *int `field:"code=EC,order=4"`
	Ingredients     []CompoundIngredient

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type CompoundProduct struct {
	Qualifier *string `field:"code=RE,order=5"`
	Id        *string `field:"code=TE,order=6"`
}

type CompoundIngredientModifierCode struct {
	Code *string `field:"code=2H,order=11"`
}

type CompoundIngredient struct {
	Product                  CompoundProduct
	Quantity                 *float64 `field:"code=ED,decimalPlaces=3,order=7"`
	DrugCost                 *float64 `field:"code=EE,decimalPlaces=2,overpunch=true,order=8"`
	BasisOfCostDetermination *string  `field:"code=UE,order=9"`

	ModifierCodeCount *int `field:"code=2G,order=10"`
	ModifierCodes     []CompoundIngredientModifierCode
}
