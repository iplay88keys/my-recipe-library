import axios, { AxiosError, AxiosResponse } from "axios";
import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { logout } from "../users/actions";
import { createRecipeAsync, fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { RecipeActionTypes, RecipeCreateResponse, RecipeListResponse, RecipeResponse } from "./types";

export function* listRecipeSaga(): Generator {
    try {
        const response = (yield call(Api.get, "/api/v1/recipes")) as AxiosResponse;

        yield put(fetchRecipesAsync.success((response.data) as RecipeListResponse));
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            const error = err as AxiosError<RecipeListResponse>;
            if (error.status === 401 || error.response && error.response.status === 401) {
                yield put(logout());
            }

            yield put(fetchRecipesAsync.failure(error));
        } else if (err instanceof Error) {
            console.log("list recipes error: ", err.message)
        } else {
            console.log("unknown list recipes error")
        }
    }
}

export function* getRecipeSaga(action: ReturnType<typeof fetchRecipeAsync.request>): Generator {
    try {
        const response = (yield call(Api.get, `/api/v1/recipes/${action.payload}`)) as AxiosResponse;

        yield put(fetchRecipeAsync.success((response.data) as RecipeResponse));
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            const error = err as AxiosError<RecipeResponse>;

            if (error.response && error.response.status === 401) {
                yield put(logout());
            }

            yield put(fetchRecipeAsync.failure(error));
        } else if (err instanceof Error) {
            console.log("get recipe error: ", err.message)
        } else {
            console.log("unknown get recipe error")
        }
    }
}

export function* createRecipeSaga(action: ReturnType<typeof createRecipeAsync.request>): Generator {
    try {
        const response = (yield call(Api.post, "/api/v1/recipes", JSON.stringify(action.payload))) as AxiosResponse;

        yield put(createRecipeAsync.success((response.data) as RecipeCreateResponse));
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            const error = err as AxiosError<RecipeCreateResponse>;

            if (error.response && error.response.status === 401) {
                yield put(logout());
            }

            if (error.response && error.response.data && error.response.data.errors) {
                action.meta(error.response.data.errors);
            }

            yield put(createRecipeAsync.failure(error));
        } else if (err instanceof Error) {
            console.log("create recipe error: ", err.message)
        } else {
            console.log("unknown create recipe error")
        }
    }
}

function* watchListRecipes() {
    yield takeEvery(RecipeActionTypes.FETCH_RECIPES_REQUEST, listRecipeSaga);
}

function* watchGetRecipe() {
    yield takeEvery(RecipeActionTypes.FETCH_RECIPE_REQUEST, getRecipeSaga);
}

function* watchCreateRecipe() {
    yield takeEvery(RecipeActionTypes.CREATE_RECIPE_REQUEST, createRecipeSaga);
}

const recipeSagas = [watchGetRecipe(), watchListRecipes(), watchCreateRecipe()];
export default recipeSagas;
