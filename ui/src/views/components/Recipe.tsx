import React from "react";
import { Link } from "react-router-dom";
import styled from "styled-components";
import { Ingredient, RecipeResponse, Step } from "../../state/ducks/recipes/types";

const StyledRecipe = styled.div`
    width: 75%;
    margin: auto;
    text-align: center;
    a {
        color: rgba(0,0,0,.5);
    }
`;
StyledRecipe.displayName = "StyledRecipe";

const StyledRecipeBreadcrumbs = styled.div`
    text-align: left;
    color: rgba(0,0,0,.5);
    :hover {
        color: black;
    }
`;
StyledRecipeBreadcrumbs.displayName = "StyledRecipeBreadcrumbs";

const StyledRecipeName = styled.div`
    margin: auto;
    font-size: 40px;
`;
StyledRecipeName.displayName = "StyledRecipeName";

const StyledRecipeImage = styled.img`
    max-height: 50vh
    overflow: auto;
`;
StyledRecipeImage.displayName = "StyledRecipeImage";

const StyledRecipeTiming = styled.div`
    display: flex;
    justify-content: space-evenly;
    margin-top: 20px;
    margin-bottom: 20px;
    > div {
        flex-grow: 1;
        border-top: 4px solid lightgray;
        border-bottom: 4px solid lightgray;
        > p {
            margin: 10px 0 10px 0;
        }
    }
    > div:not(:last-child) {
        border-right: 4px solid lightgray;
    }
`;
StyledRecipeTiming.displayName = "StyledRecipeTiming";

const StyledRecipeIngredients = styled.div`
    display: flex;
    justify-content: space-evenly;
    li {
        list-style-type: none;
        margin-bottom: 5px;
        text-align: left;
    }
`;
StyledRecipeIngredients.displayName = "StyledRecipeIngredients";

const StyledRecipeSteps = styled.div`
    margin: auto;
    text-align: left;
    li {
        margin-bottom: 5px;
    }
`;
StyledRecipeSteps.displayName = "StyledRecipeSteps";

const StyledRecipeServings = styled.div`
    text-align: left;
`;
StyledRecipeServings.displayName = "StyledRecipeServings";

interface RecipeProps {
    recipe: RecipeResponse
    loading: boolean
}

export const Recipe = ({recipe, loading}: RecipeProps) => {
    if (loading) {
        return (
            <div>
                <p>Loading recipe</p>
            </div>
        );
    }

    let source: JSX.Element = <p>Source: {recipe.source}</p>;
    if (recipe.source?.includes("http")) {
        source = <p>Source: <a href={recipe.source}>Link</a></p>;
    }

    let leftIngredients = recipe.ingredients;
    let rightIngredients = [] as Ingredient[];
    if (recipe.ingredients) {
        let half = Math.ceil(recipe.ingredients.length / 2);
        leftIngredients = recipe.ingredients.slice(0, half);
        rightIngredients = recipe.ingredients.slice(half, recipe.ingredients.length);
    }

    return (
        <StyledRecipe>
            <StyledRecipeBreadcrumbs>
                <Link to="/recipes">Recipes</Link> / <Link to="#cookbook">Cookbook</Link> / <Link
                to="#section">Section</Link>
            </StyledRecipeBreadcrumbs>
            <StyledRecipeName>{recipe.name}</StyledRecipeName>
            {recipe.description && <p>{recipe.description}</p>}
            {recipe.source && source}
            <StyledRecipeTiming>
                {recipe.prep_time &&
                <div>
                    <p>Prep: {recipe.prep_time}</p>
                </div>
                }
                {recipe.cook_time &&
                <div>
                    <p>Cook: {recipe.cook_time}</p>
                </div>
                }
                {recipe.cool_time &&
                <div>
                    <p>Cool: {recipe.cool_time}</p>
                </div>
                }
                {recipe.total_time &&
                <div>
                    <p>Total: {recipe.total_time}</p>
                </div>
                }
            </StyledRecipeTiming>
            <StyledRecipeIngredients>
                {ingredientsListElement(leftIngredients)}
                {ingredientsListElement(rightIngredients)}
            </StyledRecipeIngredients>
            <StyledRecipeSteps>
                <ol>
                    {recipe.steps && recipe.steps.map((step: Step) =>
                        <li key={step.step_number}>
                            {step.instructions}
                        </li>
                    )}
                </ol>
            </StyledRecipeSteps>
            <StyledRecipeServings>{recipe.servings} Serving{recipe.servings > 1 ? "s" : ""}</StyledRecipeServings>
        </StyledRecipe>
    );
};

function formatIngredient(ingredient: Ingredient): string {
    let formattedIngredient = "";
    if (ingredient.amount) {
        formattedIngredient += `${ingredient.amount} `;
    }

    if (ingredient.measurement) {
        formattedIngredient += `${ingredient.measurement} `;
    }

    formattedIngredient += `${ingredient.ingredient}`;

    if (ingredient.preparation) {
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
