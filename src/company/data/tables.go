package data

type (
	CompanyFoodsModel struct {
		ID          string `diinamo:"type:string;hash" json:"id"`
		Flavour     string `diinamo:"type:string" json:"flavour"`
		Ingredients string `diinamo:"type:string" json:"ingredients"`
		Price       int    `diinamo:"type:int" json:"price"`
		Photo       string `diinamo:"type:string" json:"photo,omitempty"`
		CompanySlug string `diinamo:"type:string;range;gsi:CompanyIndex;keyPairs:CompanySlug=ID" json:"company_slug,omitempty"`
	}
)
