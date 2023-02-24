import { Button } from "@material-ui/core";
import { createStyles, WithStyles, withStyles } from "@material-ui/core/styles";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";
import React from "react";
import { Link } from "react-router-dom";
import { Ingredient, RecipeResponse, Step } from "../../state/ducks/recipes/types";

interface RecipeProps extends WithStyles<typeof styles> {
    recipe: RecipeResponse
    loading: boolean
}

const styles = createStyles({
    top: {
        display: "flex",
        justifyContent: "space-evenly"
    },
    recipe: {
        width: "75%",
        margin: "auto",
        textAlign: "center",
        "& a": {
            color: "rgba(0, 0, 0, .5)"
        }
    },
    breadcrumbs: {
        textAlign: "left",
        color: "rgba(0,0,0,.5)",
        "& :hover": {
            color: "black"
        }
    },
    name: {
        margin: "auto",
        fontSize: "40px"
    },
    actions: {
        textAlign: "right",
        height: "2em",
        padding: "5px",
        margin: "5px",
        "& span": {
            padding: "5px"
        }
    },
    image: {
        maxHeight: "50vh",
        overflow: "auto"
    },
    timing: {
        display: "flex",
        justifyContent: "space-evenly",
        marginTop: "20px",
        marginBottom: "20px",
        "& > div": {
            flexGrow: 1,
            borderTop: "4px solid lightgray",
            borderBottom: "4px solid lightgray",
            "& > p": {
                margin: "10px 0 10px 0"
            }
        },
        "& > div:not(:last-child)": {
            borderRight: "4px solid lightgray"
        }
    },
    ingredients: {
        display: "flex",
        justifyContent: "space-evenly",
        "& li": {
            listStyleType: "none",
            marginBottom: "5px",
            textAlign: "left"
        }
    },
    steps: {
        margin: "auto",
        textAlign: "left",
        "& li": {
            marginBottom: "5px"
        }
    },
    servings: {
        textAlign: "left"
    }
});

function Recipe({recipe, loading, classes}: RecipeProps) {
    if (loading) {
        return (
            <div>
                <p>Loading recipe</p>
            </div>
        );
    }

    let source: JSX.Element = <p>Source: {recipe.source}</p>;
    if (recipe.source != null && recipe.source.includes("http")) {
        source = <p>Source: <a href={recipe.source}>Link</a></p>;
    }

    let leftIngredients = recipe.ingredients;
    let rightIngredients = [] as Ingredient[];
    if (recipe.ingredients != null) {
        const half = Math.ceil(recipe.ingredients.length / 2);
        leftIngredients = recipe.ingredients.slice(0, half);
        rightIngredients = recipe.ingredients.slice(half, recipe.ingredients.length);
    }

    return (
        <div className={classes.recipe}>
            <div className={classes.top}>
                <div className={classes.breadcrumbs}>
                    <Link to="/recipes">Recipes</Link> / <Link to="#cookbook">Cookbook</Link> / <Link
                    to="#section">Section</Link>
                </div>
                <div className={classes.name}>{recipe.name}</div>
                <Button
                    variant="contained"
                    color="secondary"
                    className={classes.actions}
                    startIcon={<EditIcon/>}
                >
                    Edit
                </Button>
                <Button
                    variant="contained"
                    color="secondary"
                    className={classes.actions}
                    startIcon={<DeleteIcon/>}
                >
                    Delete
                </Button>
            </div>
            {recipe.description != null && <p>{recipe.description}</p>}
            {recipe.source != null && source}
            <div className={classes.timing}>
                {recipe.prep_time != null &&
                <div>
                    <p>Prep: {recipe.prep_time}</p>
                </div>
                }
                {recipe.cook_time != null &&
                <div>
                    <p>Cook: {recipe.cook_time}</p>
                </div>
                }
                {recipe.cool_time != null &&
                <div>
                    <p>Cool: {recipe.cool_time}</p>
                </div>
                }
                {recipe.total_time != null &&
                <div>
                    <p>Total: {recipe.total_time}</p>
                </div>
                }
            </div>
            <div className={classes.ingredients}>
                {ingredientsListElement(leftIngredients)}
                {ingredientsListElement(rightIngredients)}
            </div>
            <div className={classes.steps}>
                <ol>
                    {recipe.steps && recipe.steps.map((step: Step) =>
                        <li key={step.step_number}>
                            {step.instructions}
                        </li>
                    )}
                </ol>
            </div>
            <div className={classes.servings}>{recipe.servings} Serving{recipe.servings > 1 ? "s" : ""}</div>
        </div>
    );
}

function formatIngredient(ingredient: Ingredient): string {
    let formattedIngredient = "";
    if (ingredient.amount != null) {
        formattedIngredient += `${ingredient.amount} `;
    }

    if (ingredient.measurement != null) {
        formattedIngredient += `${ingredient.measurement} `;
    }

    formattedIngredient += `${ingredient.ingredient}`;

    if (ingredient.preparation != null) {
        formattedIngredient += `, ${ingredient.preparation}`;
    }

    return formattedIngredient;
}

function ingredientsListElement(ingredients: Ingredient[]): JSX.Element {
    return (
        <div>
            <ul>
                {ingredients && ingredients.map((ingredient: Ingredient) =>
                    <li key={ingredient.ingredient_number}>
                        {formatIngredient(ingredient)}
                    </li>
                )}
            </ul>
        </div>
    );
}

export default withStyles(styles)(Recipe);
