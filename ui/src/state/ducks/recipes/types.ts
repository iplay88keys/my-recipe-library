import { APIError } from "../types";

export interface RecipeCreateResponse {
    recipe_id: number
    errors: APIError
}

export interface RecipeListResponse {
    recipes: RecipeResponse[]
}

export interface RecipeResponse {
    id: number
    name: string
    description: string
    creator: string
    servings: number
    prep_time: string
    cook_time: string
    cool_time: string
    total_time: string
    source: string
    ingredients: Ingredient[]
    steps: Step[]
}

export interface Ingredient {
    ingredient: string
    ingredient_number: number
    amount: number
    measurement: string
    preparation: string
}

export interface Step {
    step_number: number
    instructions: string
}

export interface RecipeCreateRequest {
    name: string
    description: string
    servings: number
    prep_time: string | undefined
    cook_time: string | undefined
    cool_time: string | undefined
    total_time: string | undefined
    source: string | undefined
}

export interface RecipeState {
    recipes: RecipeResponse[]
    recipe: RecipeResponse
    recipe_id: number
    loading: boolean
    creating: boolean
    error: string
}

export enum RecipeActionTypes {
    FETCH_RECIPES_REQUEST = "@@recipes/FETCH_RECIPES_REQUEST",
    FETCH_RECIPES_SUCCESS = "@@recipes/FETCH_RECIPES_SUCCESS",
    FETCH_RECIPES_FAILURE = "@@recipes/FETCH_RECIPES_FAILURE",
    FETCH_RECIPE_REQUEST = "@@recipes/FETCH_RECIPE_REQUEST",
    FETCH_RECIPE_SUCCESS = "@@recipes/FETCH_RECIPE_SUCCESS",
    FETCH_RECIPE_FAILURE = "@@recipes/FETCH_RECIPE_FAILURE",
    CREATE_RECIPE_REQUEST = "@@recipes/CREATE_RECIPE_REQUEST",
    CREATE_RECIPE_SUCCESS = "@@recipes/CREATE_RECIPE_SUCCESS",
    CREATE_RECIPE_FAILURE = "@@recipes/CREATE_RECIPE_FAILURE",
}
