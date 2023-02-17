import { ActionType, getType } from "typesafe-actions";
import * as user from "./actions";
import { loginAsync, logout, registerAsync } from "./actions";
import { UserState } from "./types";

export type UserAction = ActionType<typeof user>;

const access_token = localStorage.getItem("access_token");

const initialState: UserState = {
    registering: false,
    loggingIn: false,
    loggedIn: !!access_token,
    error: ""
};

const reducer = (state = initialState, action: UserAction) => {
    switch (action.type) {
        case getType(registerAsync.request):
            return {
                ...state,
                registering: true,
                error: ""
            };
        case getType(registerAsync.success):
            return {
                ...state,
                registering: false,
                error: ""
            };
        case getType(registerAsync.failure):
            return {
                ...state,
                registering: false,
                error: action.payload.message
            };
        case getType(loginAsync.request):
            return {
                ...state,
                loggingIn: true,
                loggedIn: false,
                error: ""
            };
        case getType(loginAsync.success):
            return {
                ...state,
                loggingIn: false,
                loggedIn: true,
                error: ""
            };
        case getType(loginAsync.failure):
            return {
                ...state,
                loggingIn: false,
                loggedIn: false,
                error: action.payload.message
            };
        case getType(logout):
            return {
                ...state,
                loggingIn: false,
                loggedIn: false,
                error: ""
            };
        default:
            return state;
    }
};

export { reducer as userReducer };
