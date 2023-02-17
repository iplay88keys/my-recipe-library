import { render, screen, within } from '@testing-library/react'
import React from "react";
import { BrowserRouter } from "react-router-dom";
import { Ingredient, RecipeResponse, Step } from "../../state/ducks/recipes/types";
import { Recipe } from "./Recipe";

describe("Recipe", () => {
    it("should render a single recipe", () => {
        const recipe = {
            id: 0,
            name: "Root Beer Float",
            description: "Delicious",
            creator: "User1",
            servings: 1,
            prep_time: "5 m",
            cook_time: "7 m",
            cool_time: "2 m",
            total_time: "14 m",
            source: "some-site",
            ingredients: [{
                ingredient: "Vanilla Ice Cream",
                ingredient_number: 0,
                amount: 1,
                measurement: "Scoop",
                preparation: "Frozen"
            }, {
                ingredient: "Root Beer",
                ingredient_number: 1
            }] as Ingredient[],
            steps: [{
                step_number: 1,
                instructions: "Place ice cream in glass."
            }, {
                step_number: 2,
                instructions: "Top with Root Beer."
            }] as Step[]
        } as RecipeResponse;

        render(
            <BrowserRouter>
                <Recipe
                    recipe={recipe}
                    loading={false}
                />
            </BrowserRouter>
        );

        const recipesLink = screen.getByRole('link', {name: /recipes/i})
        expect(recipesLink).toHaveTextContent("Recipes");
        expect(recipesLink).toHaveAttribute("href", "/recipes");

        const cookbookLink = screen.getByRole('link', {name: /cookbook/i})
        expect(cookbookLink).toHaveTextContent("Cookbook");
        expect(cookbookLink).toHaveAttribute("href", "/#cookbook");

        const sectionLink = screen.getByRole('link', {name: /section/i})
        expect(sectionLink).toHaveTextContent("Section");
        expect(sectionLink).toHaveAttribute("href", "/#section");

        expect(screen.getByText(/root beer float/i)).toBeInTheDocument();
        expect(screen.getByText(/delicious/i)).toBeInTheDocument();
        expect(screen.getByText(/source: some-site/i)).toBeInTheDocument();

        const ingredients = screen.getAllByTestId("ingredients");
        const ingredientItems = ingredients.map((ingredientList) =>
            within(ingredientList).getAllByRole("listitem").map((ingredient) =>
                ingredient.textContent
            )
        ).flat();
        expect(ingredientItems).toEqual([
            "1 Scoop Vanilla Ice Cream, Frozen",
            "Root Beer"
        ])

        expect(screen.getByText(/prep: 5 m/i)).toBeInTheDocument()
        expect(screen.getByText(/cook: 7 m/i)).toBeInTheDocument()
        expect(screen.getByText(/cool: 2 m/i)).toBeInTheDocument()
        expect(screen.getByText(/total: 14 m/i)).toBeInTheDocument()

        const steps = screen.getAllByTestId("steps");
        const stepItems = steps.map((stepList) =>
            within(stepList).getAllByRole("listitem").map((step) =>
                step.textContent
            )
        ).flat();
        expect(stepItems).toEqual([
            "Place ice cream in glass.",
            "Top with Root Beer."
        ])

        expect(screen.getByText(/1 Serving/i)).toBeInTheDocument();
    });

    it("renders the source as a link if it contains 'http'", () => {
        const recipe = {
            id: 0,
            name: "Root Beer Float",
            description: "Delicious",
            servings: 1,
            prep_time: "5 m",
            cook_time: "0 m",
            cool_time: "0 m",
            total_time: "5 m",
            creator: "User1",
            source: "http://example.com",
            ingredients: [{
                ingredient: "Root Beer",
                ingredient_number: 0
            }] as Ingredient[]
        } as RecipeResponse;

        render(
            <BrowserRouter>
                <Recipe
                    recipe={recipe}
                    loading={false}
                />
            </BrowserRouter>
        );

        expect(screen.getByRole('link', {name: /link/i})).toHaveAttribute("href", "http://example.com");
    });

    it("renders multiple servings with an 's'", () => {
        const recipe = {
            id: 0,
            name: "Root Beer Float",
            description: "Delicious",
            servings: 3,
            prep_time: "5 m",
            cook_time: "0 m",
            cool_time: "0 m",
            total_time: "5 m",
            creator: "User1",
            source: "http://example.com",
            ingredients: [{
                ingredient: "Root Beer",
                ingredient_number: 0
            }] as Ingredient[]
        } as RecipeResponse;

        render(
            <BrowserRouter>
                <Recipe
                    recipe={recipe}
                    loading={false}
                />
            </BrowserRouter>
        );

        expect(screen.getByText(/3 Servings/i)).toBeInTheDocument();
    });

    it("does not render missing data", () => {
        const recipe = {
            id: 0,
            name: "Root Beer Float",
            creator: "User1",
            servings: 1,
            ingredients: [{
                ingredient: "Root Beer",
                ingredient_number: 0
            }] as Ingredient[]
        } as RecipeResponse;

        render(
            <BrowserRouter>
                <Recipe
                    recipe={recipe}
                    loading={false}
                />
            </BrowserRouter>
        );

        expect(screen.queryByText(/delicious/i)).not.toBeInTheDocument();
        expect(screen.queryByText(/source: some-site/i)).not.toBeInTheDocument();


        expect(screen.queryByText(/prep:/i)).not.toBeInTheDocument()
        expect(screen.queryByText(/cook:/i)).not.toBeInTheDocument()
        expect(screen.queryByText(/cool:/i)).not.toBeInTheDocument()
        expect(screen.queryByText(/total:/i)).not.toBeInTheDocument()
    });

    it("should render loading info when loading", () => {
        const recipe = {} as RecipeResponse;

        render(
            <BrowserRouter>
                <Recipe
                    recipe={recipe}
                    loading={true}
                />
            </BrowserRouter>
        );

        expect(screen.getByText("Loading recipe")).toBeInTheDocument();
    });
});
