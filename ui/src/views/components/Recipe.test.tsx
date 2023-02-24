import React from "react";
import Enzyme, { shallow } from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import { Ingredient, RecipeResponse, Step } from "../../state/ducks/recipes/types";
import Recipe from "./Recipe";

Enzyme.configure({adapter: new Adapter()});

describe("Recipe", () => {
    it("should render a single recipe", () => {
        const recipe = {
            id: 0,
            name: "Root Beer Float",
            description: "Delicious",
            creator: "User1",
            servings: 1,
            prep_time: "5 m",
            cook_time: "0 m",
            cool_time: "0 m",
            total_time: "5 m",
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

        const enzymeWrapper = shallow(
            <Recipe
                recipe={recipe}
                loading={false}
            />
        );

        expect(enzymeWrapper.find("StyledRecipeBreadcrumbs").text()).toEqual("Recipes / Cookbook / Section");
        expect(enzymeWrapper.find("StyledRecipeName").text()).toEqual("Root Beer Float");
        expect(enzymeWrapper.find("StyledRecipe").childAt(2).text()).toEqual("Delicious");
        expect(enzymeWrapper.find("StyledRecipe").childAt(3).text()).toEqual("Source: some-site");

        expect(enzymeWrapper.find("StyledRecipeTiming").children()).toHaveLength(4);
        expect(enzymeWrapper.find("StyledRecipeTiming").childAt(0).text()).toEqual("Prep: 5 m");
        expect(enzymeWrapper.find("StyledRecipeTiming").childAt(1).text()).toEqual("Cook: 0 m");
        expect(enzymeWrapper.find("StyledRecipeTiming").childAt(2).text()).toEqual("Cool: 0 m");
        expect(enzymeWrapper.find("StyledRecipeTiming").childAt(3).text()).toEqual("Total: 5 m");

        expect(enzymeWrapper.find("StyledRecipeIngredients").children()).toHaveLength(2);
        expect(enzymeWrapper.find("StyledRecipeIngredients").childAt(0).text())
            .toEqual("1 Scoop Vanilla Ice Cream, Frozen");
        expect(enzymeWrapper.find("StyledRecipeIngredients").childAt(1).text()).toEqual("Root Beer");

        expect(enzymeWrapper.find("StyledRecipeSteps ol").children()).toHaveLength(2);
        expect(enzymeWrapper.find("StyledRecipeSteps ol").childAt(0).text()).toEqual("Place ice cream in glass.");
        expect(enzymeWrapper.find("StyledRecipeSteps ol").childAt(1).text()).toEqual("Top with Root Beer.");

        expect(enzymeWrapper.find("StyledRecipeServings").text()).toEqual("1 Serving");
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

        const enzymeWrapper = shallow(
            <Recipe
                recipe={recipe}
                loading={false}
            />
        );

        expect(enzymeWrapper.find("StyledRecipe").childAt(3).text()).toEqual("Source: Link");
        expect(enzymeWrapper.find("StyledRecipe").childAt(3).html())
            .toEqual(`<p>Source: <a href="http://example.com">Link</a></p>`);
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

        const enzymeWrapper = shallow(
            <Recipe
                recipe={recipe}
                loading={false}
            />
        );

        expect(enzymeWrapper.find("StyledRecipeServings").text()).toEqual("3 Servings");
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

        const enzymeWrapper = shallow(
            <Recipe
                recipe={recipe}
                loading={false}
            />
        );

        expect(enzymeWrapper.find("StyledRecipe").childAt(2).text()).not.toContain("Delicious");
        expect(enzymeWrapper.find("StyledRecipe").childAt(3).text()).not.toContain("Source: some-site");

        expect(enzymeWrapper.find("StyledRecipeTiming").children()).toHaveLength(0);
        expect(enzymeWrapper.find("StyledRecipeTiming").text()).not.toContain("Prep:");
        expect(enzymeWrapper.find("StyledRecipeTiming").text()).not.toContain("Cook:");
        expect(enzymeWrapper.find("StyledRecipeTiming").text()).not.toContain("Cool:");
        expect(enzymeWrapper.find("StyledRecipeTiming").text()).not.toContain("Total:");
    });

    it("should render loading info when loading", () => {
        const props = {
            recipe: {} as RecipeResponse
        };

        const enzymeWrapper = shallow(
            <Recipe
                recipe={props.recipe}
                loading={true}
            />
        );

        expect(enzymeWrapper.find("div")).toHaveLength(1);
        expect(enzymeWrapper.find("div p").text()).toBe("Loading recipe");
    });
});
