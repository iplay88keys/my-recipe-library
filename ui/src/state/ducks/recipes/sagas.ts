import { AxiosResponse } from "axios";
import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { logout } from "../users/actions";
import { createRecipeAsync, fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { RecipeActionTypes, RecipeCreateResponse, RecipeListResponse, RecipeResponse } from "./types";

export function* listRecipeSaga(): Generator {
    try {
        const response = (yield call(Api.get, "/api/v1/recipes")) as AxiosResponse;

        yield put(fetchRecipesAsync.success((response.data) as RecipeListResponse));
    } catch (err) {
        if (err.response && err.response.status === 401) {
            yield put(logout());
        }

        yield put(fetchRecipesAsync.failure(err));
    }
}

export function* getRecipeSaga(action: ReturnType<typeof fetchRecipeAsync.request>): Generator {
    try {
        const response = (yield call(Api.get, `/api/v1/recipes/${action.payload}`)) as AxiosResponse;

        yield put(fetchRecipeAsync.success((response.data) as RecipeResponse));
    } catch (err) {
        if (err.response && err.response.status === 401) {
            yield put(logout());
        }

        yield put(fetchRecipeAsync.failure(err));
    }
}

export function* createRecipeSaga(action: ReturnType<typeof createRecipeAsync.request>): Generator {
    try {
        const response = (yield call(Api.post, "/api/v1/recipes", JSON.stringify(action.payload))) as AxiosResponse;

        const data = (response.data) as RecipeCreateResponse;
        yield put(createRecipeAsync.success(data.recipe_id));
    } catch (err) {
        if (err.response && err.response.status === 401) {
            yield put(logout());
        }

        if (err.response && err.response.data && err.response.data.errors) {
            action.meta(err.response.data.errors);
        }

        yield put(createRecipeAsync.failure(err));
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
