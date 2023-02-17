import axios, { AxiosError, AxiosResponse } from "axios";
import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { history } from "../../../helpers/history";
import { loginAsync, registerAsync } from "./actions";
import { LoginResponse, LogoutRequest, RegisterResponse, UserActionTypes } from "./types";

export function* registerSaga(action: ReturnType<typeof registerAsync.request>): Generator {
    try {
        yield call(Api.post, "/api/v1/users/register", JSON.stringify(action.payload), false);

        yield put(registerAsync.success());
        history.push("/login");
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            const error = err as AxiosError<RegisterResponse>;
            if (error.response && error.response.data && error.response.data.errors) {
                action.meta(error.response.data.errors);
            }

            yield put(registerAsync.failure(err));
        } else if (err instanceof Error) {
            console.log(err.message)
        } else {
            console.log("unknown create recipe error")
        }
    }
}

export function* loginSaga(action: ReturnType<typeof loginAsync.request>): Generator {
    try {
        const response = (yield call(Api.post, "/api/v1/users/login", JSON.stringify(action.payload), false)) as AxiosResponse;

        yield put(loginAsync.success());

        const data = (response.data) as LoginResponse;
        localStorage.setItem("access_token", data.access_token);

        history.push("/recipes");
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            const error = err as AxiosError<LoginResponse>;
            if (error.response && error.response.data && error.response.data.errors) {
                action.meta(error.response.data.errors);
            }

            yield put(loginAsync.failure(error));
        } else if (err instanceof Error) {
            console.log(err.message)
        } else {
            console.log("unknown create recipe error")
        }
    }
}

export function* logoutSaga(): Generator {
    const token = localStorage.getItem("access_token") || null;

    const payload = {
        access_token: token
    } as LogoutRequest;

    /* eslint no-empty: ["error", { "allowEmptyCatch": true }] */
    try {
        yield call(Api.post, "/api/v1/users/logout", JSON.stringify(payload));
    } catch (err) {}

    localStorage.removeItem("access_token");

    history.push("/login");
}

function* watchRegister() {
    yield takeEvery(UserActionTypes.REGISTER_REQUEST, registerSaga);
}

function* watchLogin() {
    yield takeEvery(UserActionTypes.LOGIN_REQUEST, loginSaga);
}

function* watchLogout() {
    yield takeEvery(UserActionTypes.LOGOUT, logoutSaga);
}

const userSagas = [watchRegister(), watchLogin(), watchLogout()];
export default userSagas;
