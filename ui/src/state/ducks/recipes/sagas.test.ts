import { AxiosError, AxiosHeaders } from "axios";
import { AxiosResponse } from "axios/index";
import { testSaga } from "redux-saga-test-plan";
import { action } from "typesafe-actions";
import Api from "../../../api/api";
import { UserActionTypes } from "../users/types";
import { createRecipeSaga, getRecipeSaga, listRecipeSaga } from "./sagas";
import {
    Ingredient,
    RecipeActionTypes,
    RecipeCreateRequest,
    RecipeCreateResponse,
    RecipeListResponse,
    RecipeResponse,
    Step,
} from "./types";

describe.only("recipes", () => {
    describe.only("listRecipeSaga", () => {
        it("returns the recipes and dispatches the success action", () => {
            const recipes = {
                recipes: [{
                    id: 0,
                    name: "First",
                    description: "One"
                }, {
                    id: 1,
                    name: "Second",
                    description: "Two"
                }] as RecipeResponse[]
            } as RecipeListResponse;

            const axiosResponse = {
                data: recipes
            } as AxiosResponse;

            const saga = testSaga(listRecipeSaga);
            saga.next().call(Api.get, "/api/v1/recipes");
            saga.next(axiosResponse).put({type: RecipeActionTypes.FETCH_RECIPES_SUCCESS, payload: recipes});
            saga.next().isDone();
        });

        it("returns an error when there is a non-200 response and dispatches the error action", () => {
            const axiosError = new AxiosError("some-error", "404", undefined, undefined, {
                data: {errors: "some-error"},
                status: 404,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const saga = testSaga(listRecipeSaga);
            saga.next().call(Api.get, "/api/v1/recipes");
            saga.throw(axiosError).put({type: RecipeActionTypes.FETCH_RECIPES_FAILURE, payload: axiosError});
            saga.next().isDone();
        });

        it("dispatches the logout action when there is a 401 response", () => {
            const axiosError = new AxiosError("unauthorized", "401", undefined, undefined, {
                data: {errors: "unauthorized"},
                status: 401,
                statusText: "Unauthorized",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const saga = testSaga(listRecipeSaga);
            saga.next().call(Api.get, "/api/v1/recipes");
            saga.throw(axiosError).put({type: UserActionTypes.LOGOUT});
            saga.next().put({type: RecipeActionTypes.FETCH_RECIPES_FAILURE, payload: axiosError});
            saga.next().isDone();
        });
    });

    describe.only("getRecipeSaga", () => {
        it("returns a recipe and dispatches the success action", () => {
            const recipe = {
                id: 0,
                name: "Root Beer Float",
                description: "Delicious",
                creator: "User1",
                servings: 1,
                prep_time: "5 m",
                total_time: "5 m",
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

            const axiosResponse = {
                data: recipe
            } as AxiosResponse;

            const saga = testSaga(getRecipeSaga, action(RecipeActionTypes.FETCH_RECIPE_REQUEST, 0));
            saga.next().call(Api.get, "/api/v1/recipes/0");
            saga.next(axiosResponse).put({type: RecipeActionTypes.FETCH_RECIPE_SUCCESS, payload: recipe});
            saga.next().isDone();
        });

        it("returns an error when there is a non-200 response and dispatches the error action", () => {
            const axiosError = new AxiosError("some-error", "404", undefined, undefined, {
                data: {errors: "some-error"},
                status: 404,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const saga = testSaga(getRecipeSaga, action(RecipeActionTypes.FETCH_RECIPE_REQUEST, 0));
            saga.next().call(Api.get, "/api/v1/recipes/0");
            saga.throw(axiosError).put({type: RecipeActionTypes.FETCH_RECIPE_FAILURE, payload: axiosError});
            saga.next().isDone();
        });

        it("dispatches the logout action when there is a 401 response", () => {
            const axiosError = new AxiosError("unauthorized", "401", undefined, undefined, {
                data: {errors: "unauthorized"},
                status: 401,
                statusText: "Unauthorized",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const saga = testSaga(getRecipeSaga, action(RecipeActionTypes.FETCH_RECIPE_REQUEST, 0));
            saga.next().call(Api.get, "/api/v1/recipes/0");
            saga.throw(axiosError).put({type: UserActionTypes.LOGOUT});
            saga.next().put({type: RecipeActionTypes.FETCH_RECIPE_FAILURE, payload: axiosError});
            saga.next().isDone();
        });
    });

    describe.only("createRecipeSaga", () => {
        it("dispatches the success action", () => {
            const resp = {
                recipe_id: 10
            } as RecipeCreateResponse;

            const axiosResponse = {
                data: resp
            } as AxiosResponse;

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            const mockSetErrors = jest.fn();
            const saga = testSaga(createRecipeSaga, action(RecipeActionTypes.CREATE_RECIPE_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/recipes", JSON.stringify(req));
            saga.next(axiosResponse).put({type: RecipeActionTypes.CREATE_RECIPE_SUCCESS, payload: resp});
            saga.next().isDone();

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("calls the provided function with the error payload", () => {
            const errors = {errors: "some-error"}

            const axiosError = new AxiosError("some-error", "404", undefined, undefined, {
                data: errors,
                status: 404,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            const mockSetErrors = jest.fn();

            const saga = testSaga(createRecipeSaga, action(RecipeActionTypes.CREATE_RECIPE_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/recipes", JSON.stringify(req));
            saga.throw(axiosError).put({type: RecipeActionTypes.CREATE_RECIPE_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).toHaveBeenCalled();
            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, errors.errors);
        });

        it("dispatches the logout action when there is a 401 response", () => {
            const errors = {errors: "unauthorized"}

            const axiosError = new AxiosError("unauthorized", "401", undefined, undefined, {
                data: errors,
                status: 401,
                statusText: "Unauthorized",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            const mockSetErrors = jest.fn();

            const saga = testSaga(createRecipeSaga, action(RecipeActionTypes.CREATE_RECIPE_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/recipes", JSON.stringify(req));
            saga.throw(axiosError).put({type: UserActionTypes.LOGOUT});
            saga.next().put({type: RecipeActionTypes.CREATE_RECIPE_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).toHaveBeenCalled();
            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, errors.errors);
        });

        it("returns an error if there is no payload and there is an error", () => {
            const axiosError = new AxiosError("unauthorized", "500", undefined, undefined, {
                data: undefined,
                status: 500,
                statusText: "Server error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            const mockSetErrors = jest.fn();

            const saga = testSaga(createRecipeSaga, action(RecipeActionTypes.CREATE_RECIPE_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/recipes", JSON.stringify(req));
            saga.throw(axiosError).put({type: RecipeActionTypes.CREATE_RECIPE_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).not.toHaveBeenCalled();
        });
    });
});
