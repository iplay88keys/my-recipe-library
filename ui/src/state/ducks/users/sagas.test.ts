import { AxiosError, AxiosHeaders, AxiosResponse } from "axios";
import { testSaga } from "redux-saga-test-plan";
import { action } from "typesafe-actions";
import Api from "../../../api/api";
import { history } from "../../../helpers/history"
import { loginSaga, registerSaga } from "./sagas";
import { LoginRequest, LoginResponse, RegisterRequest, UserActionTypes } from "./types";

describe.only("users", () => {
    describe.only("registerSaga", () => {
        it("dispatches the success action", () => {
            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            const historyMock = history;
            historyMock.push = jest.fn();

            const mockSetErrors = jest.fn();
            const saga = testSaga(registerSaga, action(UserActionTypes.REGISTER_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/users/register", JSON.stringify(req), false);
            saga.next().put({type: UserActionTypes.REGISTER_SUCCESS});
            saga.next().isDone();

            expect(mockSetErrors).not.toHaveBeenCalled();

            expect(historyMock.push).toHaveBeenCalledWith("/login");
        });

        it("calls the provided function with the error payload", () => {
            const errors = {errors: "some-error"}

            const axiosError = new AxiosError("some-error", "400", undefined, undefined, {
                data: errors,
                status: 400,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            const historyMock = history;
            historyMock.push = jest.fn();

            const mockSetErrors = jest.fn();
            const saga = testSaga(registerSaga, action(UserActionTypes.REGISTER_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/users/register", JSON.stringify(req), false);
            saga.throw(axiosError).put({type: UserActionTypes.REGISTER_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).toHaveBeenCalled();
            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, errors.errors);

            expect(historyMock.push).not.toHaveBeenCalled();
        });

        it("returns an error if there is no payload and there is an error", () => {
            const axiosError = new AxiosError("some-error", "400", undefined, undefined, {
                data: undefined,
                status: 400,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            const historyMock = history;
            historyMock.push = jest.fn();

            const mockSetErrors = jest.fn();
            const saga = testSaga(registerSaga, action(UserActionTypes.REGISTER_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/users/register", JSON.stringify(req), false);
            saga.throw(axiosError).put({type: UserActionTypes.REGISTER_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).not.toHaveBeenCalled();

            expect(historyMock.push).not.toHaveBeenCalled();
        });
    });

    describe.only("loginSaga", () => {
        it("dispatches the success action", () => {
            const loginResp = {
                access_token: "some-token",
                refresh_token: "another-token",
                errors: {}
            } as LoginResponse;

            const axiosResponse = {
                data: loginResp
            } as AxiosResponse;

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            jest.spyOn(Object.getPrototypeOf(window.localStorage), 'setItem')
            // eslint-disable-next-line jest/unbound-method
            Object.setPrototypeOf(window.localStorage.setItem, jest.fn())

            const historyMock = history;
            historyMock.push = jest.fn();

            const mockSetErrors = jest.fn();
            const saga = testSaga(loginSaga, action(UserActionTypes.LOGIN_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/users/login", JSON.stringify(req), false);
            saga.next(axiosResponse).put({type: UserActionTypes.LOGIN_SUCCESS});
            saga.next().isDone();

            expect(historyMock.push).toHaveBeenCalledWith("/recipes");
            expect(window.localStorage.setItem).toHaveBeenCalledWith("access_token", loginResp.access_token)
        });

        it("calls the provided function with the error payload", () => {
            const errors = {errors: "some-error"}

            const axiosError = new AxiosError("some-error", "400", undefined, undefined, {
                data: errors,
                status: 400,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            const historyMock = history;
            historyMock.push = jest.fn();

            const mockSetErrors = jest.fn();
            const saga = testSaga(loginSaga, action(UserActionTypes.LOGIN_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/users/login", JSON.stringify(req), false);
            saga.throw(axiosError).put({type: UserActionTypes.LOGIN_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).toHaveBeenCalled();
            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, errors.errors);

            expect(historyMock.push).not.toHaveBeenCalled();
        });

        it("returns an error if there is no payload and there is an error", () => {
            const axiosError = new AxiosError("some-error", "400", undefined, undefined, {
                data: undefined,
                status: 400,
                statusText: "some-error",
                headers: {},
                config: {
                    headers: new AxiosHeaders(),
                }
            })

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            const historyMock = history;
            historyMock.push = jest.fn();

            const mockSetErrors = jest.fn();
            const saga = testSaga(loginSaga, action(UserActionTypes.LOGIN_REQUEST, req, mockSetErrors));
            saga.next().call(Api.post, "/api/v1/users/login", JSON.stringify(req), false);
            saga.throw(axiosError).put({type: UserActionTypes.LOGIN_FAILURE, payload: axiosError});
            saga.next().isDone();

            expect(mockSetErrors).not.toHaveBeenCalled();

            expect(historyMock.push).not.toHaveBeenCalled();
        });
    });
});
