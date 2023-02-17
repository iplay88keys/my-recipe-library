import { render, screen, within } from "@testing-library/react";
import userEvent from '@testing-library/user-event'
import React from "react";
import { loginAsync } from "../../state/ducks/users/actions";
import Login from "./Login";

describe("Login", () => {
    it("should render a form and allow logging in", async () => {
        const login = jest.fn();

        render(
            <Login
                login={login}
            />
        );

        const user = userEvent.setup()

        await user.type(screen.getByPlaceholderText("Username or Email Address"), "test-username")
        await user.type(screen.getByPlaceholderText("Password"), "test-password")

        const submitButton = screen.getByRole("button")
        expect(submitButton).toHaveTextContent("Login");
        await user.click(submitButton)

        expect(login).toHaveBeenCalledWith({
            "login": "test-username",
            "password": "test-password",
        }, expect.any(Function))
    });

    describe("form validation errors", () => {
        describe("username", () => {
            it("is required", async () => {
                const login = jest.fn(loginAsync.request);

                render(
                    <Login
                        login={login}
                    />
                );

                const user = userEvent.setup()

                expect(screen.queryByText("Required")).not.toBeInTheDocument();

                await user.type(screen.getByPlaceholderText("Password"), "test-password")

                const submitButton = screen.getByRole("button")
                expect(submitButton).toHaveTextContent("Login");
                await user.click(submitButton)

                const loginSection = screen.getByTestId("loginSection")
                expect(within(loginSection).getByText("Required")).toBeInTheDocument()
            });
        });

        describe("password", () => {
            it("is required", async () => {
                const login = jest.fn(loginAsync.request);

                render(
                    <Login
                        login={login}
                    />
                );

                const user = userEvent.setup()

                expect(screen.queryByText("Required")).not.toBeInTheDocument();

                await user.type(screen.getByPlaceholderText("Username or Email Address"), "test-user")

                const submitButton = screen.getByRole("button")
                expect(submitButton).toHaveTextContent("Login");
                await user.click(submitButton)

                const passwordSection = screen.getByTestId("passwordSection")
                expect(within(passwordSection).getByText("Required")).toBeInTheDocument()
            });
        });
    });

    // There are no api validation errors for specific fields at this point
    //
    // describe("api validation errors", () => {
    //     it("displays username errors", async () => {
    //         const login = jest.fn(loginAsync.request);
    //         login.mockImplementation((payload: LoginRequest, meta: (errors: FormikErrors<LoginFormValues>) => void): any => {
    //             meta({alert: "error logging in"} as FormikErrors<LoginFormValues>);
    //         });
    //
    //         const {rerender} = render(
    //             <Login
    //                 login={login}
    //             />
    //         );
    //
    //         rerender(
    //             <Login
    //                 login={login}
    //             />
    //         );
    //
    //         const user = userEvent.setup()
    //
    //         expect(screen.queryByText("Required")).not.toBeInTheDocument();
    //
    //         await user.type(screen.getByPlaceholderText("Username or Email Address"), "test-username")
    //         await user.type(screen.getByPlaceholderText("Password"), "test-password")
    //
    //         const submitButton = screen.getByRole("button")
    //         expect(submitButton).toHaveTextContent("Login");
    //         await user.click(submitButton)
    //
    //         expect(login).toHaveBeenCalled();
    //     });
    // });
});
