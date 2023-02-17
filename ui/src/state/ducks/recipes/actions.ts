import { FormikErrors } from "formik";
import { createAsyncAction } from "typesafe-actions";
import { NewRecipeFormValues } from "../../../views/components/NewRecipe";
import {
    RecipeActionTypes,
    RecipeCreateRequest,
    RecipeCreateResponse,
    RecipeListResponse,
    RecipeResponse
} from "./types";

export const fetchRecipesAsync = createAsyncAction(
    RecipeActionTypes.FETCH_RECIPES_REQUEST,
    RecipeActionTypes.FETCH_RECIPES_SUCCESS,
    RecipeActionTypes.FETCH_RECIPES_FAILURE
)<undefined, RecipeListResponse, Error>();

export const fetchRecipeAsync = createAsyncAction(
    RecipeActionTypes.FETCH_RECIPE_REQUEST,
    RecipeActionTypes.FETCH_RECIPE_SUCCESS,
    RecipeActionTypes.FETCH_RECIPE_FAILURE
)<number, RecipeResponse, Error>();

export const createRecipeAsync = createAsyncAction(
    RecipeActionTypes.CREATE_RECIPE_REQUEST,
    RecipeActionTypes.CREATE_RECIPE_SUCCESS,
    RecipeActionTypes.CREATE_RECIPE_FAILURE
)<[RecipeCreateRequest, (errors: FormikErrors<NewRecipeFormValues>) => void], RecipeCreateResponse, Error>();
