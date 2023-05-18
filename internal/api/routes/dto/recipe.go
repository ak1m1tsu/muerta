package dto

type CreateRecipe struct {
	UserID      int          `json:"id_user"               validate:"required,gt=0"                   exmaple:"1"`
	Name        string       `json:"name"                  validate:"required,gte=2,lte=100,notblank"             example:"Салат"`
	Description string       `json:"description,omitempty" validate:"lte=200"                                     example:"Салат из миндаля"`
	Steps       []RecipeStep `json:"steps"                 validate:"required"                                                               swaggertype:"array"`
	Ingredients []Ingredient `json:"ingredients"           validate:"required"                                                               swaggertype:"array"`
}

type FindRecipe struct {
	ID          int        `json:"id"                    example:"1"`
	Name        string     `json:"name"                  example:"Салат"`
	Description string     `json:"description,omitempty" example:"Салат из миндаля"`
	Steps       []FindStep `json:"steps,omitempty"`
}

type UpdateRecipe struct {
	Name        string `json:"name"        validate:"gte=2,lte=100" example:"Салат"`
	Description string `json:"description" validate:"lte=200"       example:"Салат из миндаля"`
}

type DeleteRecipeStep struct {
	Place int `json:"place" validate:"required,gt=0" example:"1"`
}

type CreateRecipeStep struct {
	Place int `json:"place" validate:"required,gt=0" example:"1"`
}
