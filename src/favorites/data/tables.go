package data

type (
	CompanyMeals struct {
		Company         string `diinamo:"type:string"`
		Slug            string `diinamo:"type:string;hash"`
		SocialName      string `diinamo:"type:string"`
		CPNJ            string `diinamo:"type:string;range"`
		MealID          string `diinamo:"type:string"`
		MealFlavour     string `diinamo:"type:string"`
		MealSlug        string `diinamo:"type:string"`
		MealPrice       int    `diinamo:"type:int"`
		MealIngredients string `diinamo:"type:string"`
		MealPhoto       string `diinamo:"type:string"`
	}

	Favorites struct {
		FavoriteID  string `diinamo:"type:string;hash"`
		UserID      string `diinamo:"type:string;range"`
		CompanySlug string `diinamo:"type:string;gsi:CompanyIndex;keyPairs:UserID=CompanySlug"`
		MealID      string `diinamo:"type:string;gsi:MealIndex;keyPairs:UserID=MealID"`
		MealSlug    string `diinamo:"type:string;"`
	}
)
