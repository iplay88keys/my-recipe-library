import { all } from "redux-saga/effects";

import { recipeReducer } from "./recipes/reducers";
import recipeSagas from "./recipes/sagas";
import { RecipeState } from "./recipes/types";
import { userReducer } from "./users/reducers";
import userSagas from "./users/sagas";
import { UserState } from "./users/types";

export interface ApplicationState {
    recipes: RecipeState,
    users: UserState
}

export const reducer = {
    recipes: recipeReducer,
    users: userReducer
};

export function* rootSaga() {
    yield all([
        ...recipeSagas,
        ...userSagas
    ]);
}

