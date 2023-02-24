import { AxiosResponse } from "axios";
import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { history } from "../../../helpers/history";
import { loginAsync, registerAsync } from "./actions";
import { LoginResponse, LogoutRequest, UserActionTypes } from "./types";

export function* registerSaga(action: ReturnType<typeof registerAsync.request>): Generator {
    try {
        yield call(Api.post, "/api/v1/users/register", JSON.stringify(action.payload), false);

        yield put(registerAsync.success());
        history.push("/login");
    } catch (err) {
        if (err.response && err.response.data && err.response.data.errors) {
            action.meta(err.response.data.errors);
        }

        yield put(registerAsync.failure(err));
    }
}

export function* loginSaga(action: ReturnType<typeof loginAsync.request>): Generator {
    try {
        const response = (yield call(Api.post, "/api/v1/users/login", JSON.stringify(action.payload), false)) as AxiosResponse;

        yield put(loginAsync.success());

        const data = (response.data) as LoginResponse;
        localStorage.setItem("access_token", data.access_token);

        history.push("/recipes");
    } catch (err) {
        if (err.response && err.response.data && err.response.data.errors) {
            action.meta(err.response.data.errors);
        }

        yield put(loginAsync.failure(err));
    }
}

export function* logoutSaga(): Generator {
    const token = localStorage.getItem("access_token") || null;

    const payload = {
        access_token: token
    } as LogoutRequest;

    try {
        yield call(Api.post, "/api/v1/users/logout", JSON.stringify(payload));
    } catch (err) {
        console.log("Error logging out")
    }

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
