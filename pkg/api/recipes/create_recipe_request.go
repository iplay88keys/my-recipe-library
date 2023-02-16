package recipes

type CreateRecipeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Servings    int    `json:"servings"`
	PrepTime    string `json:"prep_time"`
	CookTime    string `json:"cook_time"`
	CoolTime    string `json:"cool_time"`
	TotalTime   string `json:"total_time"`
	Source      string `json:"source"`
}

func (a *CreateRecipeRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if a.Name == "" {
		errors["name"] = "Required"
	}

	if a.Description == "" {
		errors["description"] = "Required"
	}

	if a.Servings == 0 {
		errors["servings"] = "Required"
	}

	return errors
}
