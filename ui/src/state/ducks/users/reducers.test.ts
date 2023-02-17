import { loginAsync, registerAsync } from "./actions";
import { userReducer } from "./reducers";
import { LoginRequest, RegisterRequest } from "./types";

describe("reducer", () => {
    describe("register", function () {
        it("should handle REGISTER_REQUEST", () => {
            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            const mockSetErrors = jest.fn();

            const updatedState = userReducer(undefined, registerAsync.request(req, mockSetErrors));

            expect(updatedState).toEqual({
                "error": "",
                "loggedIn": false,
                "loggingIn": false,
                "registering": true
            });
        });

        it("should handle REGISTER_SUCCESS", () => {
            const updatedState = userReducer(undefined, registerAsync.success());
            expect(updatedState).toEqual({
                "error": "",
                "loggedIn": false,
                "loggingIn": false,
                "registering": false
            });
        });

        it("should handle REGISTER_FAILURE", () => {
            const err = {
                message: "some error"
            } as Error;

            const updatedState = userReducer(undefined, registerAsync.failure(err));

            expect(updatedState).toEqual({
                "error": "some error",
                "loggedIn": false,
                "loggingIn": false,
                "registering": false
            });
        });
    });

    describe("login", function () {
        it("should handle LOGIN_REQUEST", () => {
            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            const mockSetErrors = jest.fn();

            const updatedState = userReducer(undefined, loginAsync.request(req, mockSetErrors));

            expect(updatedState).toEqual({
                "error": "",
                "loggedIn": false,
                "loggingIn": true,
                "registering": false
            });
        });

        it("should handle LOGIN_SUCCESS", () => {
            const updatedState = userReducer(undefined, loginAsync.success());
            expect(updatedState).toEqual({
                "error": "",
                "loggedIn": true,
                "loggingIn": false,
                "registering": false
            });
        });

        it("should handle LOGIN_FAILURE", () => {
            const err = {
                message: "some error"
            } as Error;

            const updatedState = userReducer(undefined, loginAsync.failure(err));
            expect(updatedState).toEqual({
                "error": "some error",
                "loggedIn": false,
                "loggingIn": false,
                "registering": false
            });
        });
    });
});
