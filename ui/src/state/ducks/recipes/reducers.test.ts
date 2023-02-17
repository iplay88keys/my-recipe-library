import { createRecipeAsync, fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { recipeReducer } from "./reducers";
import {
    Ingredient,
    RecipeResponse,
    RecipeCreateRequest,
    RecipeListResponse,
    Step,
    RecipeCreateResponse
} from "./types";

describe("reducer", () => {
    describe("list", () => {
        it("should handle FETCH_RECIPES_REQUEST", () => {
            const updatedState = recipeReducer(undefined, fetchRecipesAsync.request());

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: true,
                creating: false,
                error: ""
            });
        });

        it("should handle FETCH_RECIPES_SUCCESS", () => {
            const recipes = {
                recipes: [{
                    id: 0,
                    name: "First",
                    description: "One"
                }] as RecipeResponse[]
            } as RecipeListResponse;

            const updatedState = recipeReducer(undefined, fetchRecipesAsync.success(recipes));

            expect(updatedState).toEqual({
                recipes: recipes.recipes,
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: false,
                creating: false,
                error: ""
            });
        });

        it("should handle FETCH_RECIPES_FAILURE", () => {
            const err = {
                message: "some error"
            } as Error;

            const updatedState = recipeReducer(undefined, fetchRecipesAsync.failure(err));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: false,
                creating: false,
                error: "some error"
            });
        });
    });

    describe("get", () => {
        it("should handle FETCH_RECIPE_REQUEST", () => {
            const updatedState = recipeReducer(undefined, fetchRecipeAsync.request(1));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: true,
                creating: false,
                error: ""
            });
        });

        it("should handle FETCH_RECIPE_SUCCESS", () => {
            const recipe = {
                id: 0,
                name: "Root Beer Float",
                description: "Delicious",
                creator: "User1",
                servings: 1,
                prep_time: "5 m",
                total_time: "5 m",
                source: "Some Website",
                ingredients: [{
                    ingredient: "Vanilla Ice Cream",
                    ingredient_number: 0,
                    amount: 1,
                    measurement: "Scoop",
                    preparation: "Frozen"
                }] as Ingredient[],
                steps: [{
                    step_number: 1,
                    instructions: "Place ice cream in glass."
                }] as Step[]
            } as RecipeResponse;

            const updatedState = recipeReducer(undefined, fetchRecipeAsync.success(recipe));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: recipe,
                recipe_id: 0,
                loading: false,
                creating: false,
                error: ""
            });
        });

        it("should handle FETCH_RECIPE_FAILURE", () => {
            const err = {
                message: "some error"
            } as Error;

            const updatedState = recipeReducer(undefined, fetchRecipeAsync.failure(err));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: false,
                creating: false,
                error: "some error"
            });
        });
    });

    describe("create", () => {
        it("should handle CREATE_RECIPE_REQUEST", () => {
            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 1
            } as RecipeCreateRequest;

            const mockSetErrors = jest.fn();

            const updatedState = recipeReducer(undefined, createRecipeAsync.request(req, mockSetErrors));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: false,
                creating: true,
                error: ""
            });
        });

        it("should handle CREATE_RECIPE_SUCCESS", () => {
            const updatedState = recipeReducer(undefined, createRecipeAsync.success(({recipe_id: 2}) as RecipeCreateResponse));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {},
                recipe_id: 2,
                loading: false,
                creating: false,
                error: ""
            });
        });

        it("should handle CREATE_RECIPE_FAILURE", () => {
            const err = {
                message: "some error"
            } as Error;

            const updatedState = recipeReducer(undefined, createRecipeAsync.failure(err));

            expect(updatedState).toEqual({
                recipes: [],
                recipe: {} as RecipeResponse,
                recipe_id: 0,
                loading: false,
                creating: false,
                error: "some error"
            });
        });
    });
});
