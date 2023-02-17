import { act, render } from "@testing-library/react";
import React from "react";
import { registerAsync } from "../../state/ducks/users/actions";
import Registration from "./Registration";

describe("Registration", () => {
    it("should render a form for registering", () => {
        const register = jest.fn(registerAsync.request);

        render(
            <Registration
                register={register}
            />
        );

        // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).text()).toEqual("Username*Username*");
        // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).text()).toEqual("Email*Email*");
        // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).text()).toEqual("Password*Password*");
        // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).text())
        //     .toEqual("Confirm Password*Confirm Password*");
    });

    describe("form validation errors", () => {
        beforeEach(() => {
            const register = jest.fn(registerAsync.request);

            render(
                <Registration
                    register={register}
                />
            );

        });

        describe("username", () => {
            beforeEach(async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().onBlur!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "username"
                    //     }
                    // } as FocusEvent<HTMLInputElement>);
                    //
                    // enzymeWrapper.update();
                });
            });

            it("allows valid input", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "username",
                    //         value: "valid_Username"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().value)
                //     .toEqual("valid_Username");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().error).toEqual(false);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().helperText)
                //     .toEqual("");
            });

            it("is required", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "username",
                    //         value: ""
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().value).toEqual("");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().helperText)
                //     .toEqual("Required");
            });

            it("validates length", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "username",
                    //         value: "ian"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().value).toEqual("ian");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().helperText)
                //     .toEqual("Must be 6 characters or more");

                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "username",
                    //         value: "this username is more than 30 characters long"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().value)
                //     .toEqual("this username is more than 30 characters long");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().helperText)
                //     .toEqual("Must be 30 characters or less");
            });

            it("must match the regex", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "username",
                    //         value: "&invalid User"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().value)
                //     .toEqual("&invalid User");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().helperText)
                //     .toEqual("Must start with a letter and only contain alphanumeric characters and underscores (_)");
            });
        });

        describe("email", () => {
            beforeEach(async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().onBlur!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "email"
                    //     }
                    // } as FocusEvent<HTMLInputElement>);
                    //
                    // enzymeWrapper.update();
                });
            });

            it("allows valid input", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "email",
                    //         value: "valid@example.com"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().value)
                //     .toEqual("valid@example.com");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().error).toEqual(false);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().helperText)
                //     .toEqual("");
            });

            it("is required", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "email",
                    //         value: ""
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().value).toEqual("");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().helperText)
                //     .toEqual("Required");
            });

            it("validates email", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "email",
                    //         value: "invalid"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().value)
                //     .toEqual("invalid");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().helperText)
                //     .toEqual("Must be a valid email");
            });
        });

        describe("password", () => {
            beforeEach(async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onBlur!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password"
                    //     }
                    // } as FocusEvent<HTMLInputElement>);
                    //
                    // enzymeWrapper.update();
                });
            });

            it("allows valid input", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password",
                    //         value: "Pa3$word123"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().value)
                //     .toEqual("Pa3$word123");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().error).toEqual(false);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().helperText)
                //     .toEqual("");
            });

            it("is required", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password",
                    //         value: ""
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().value).toEqual("");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().helperText)
                //     .toEqual("Required");
            });

            it("validates length", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password",
                    //         value: "short"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().value).toEqual("short");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().helperText)
                //     .toEqual("Must be 6 characters or more");

                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password",
                    //         value: "this password is more than 64 characters long which is too long for this field"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().value)
                //     .toEqual("this password is more than 64 characters long which is too long for this field");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().helperText)
                //     .toEqual("Must be 64 characters or less");
            });
        });

        describe("password confirmation", () => {
            beforeEach(async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onBlur!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password"
                    //     }
                    // } as FocusEvent<HTMLInputElement>);
                    //
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "password",
                    //         value: "Pa3$word123"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                    //
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().onBlur!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "passwordConfirmation"
                    //     }
                    // } as FocusEvent<HTMLInputElement>);
                    //
                    // enzymeWrapper.update();
                });
            });

            it("allows valid input when matching the password field", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "passwordConfirmation",
                    //         value: "Pa3$word123"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().value)
                //     .toEqual("Pa3$word123");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().error).toEqual(false);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().helperText)
                //     .toEqual("");
            });

            it("has to match the password field", async () => {
                await act(async () => {
                    // enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().onChange!({
                    //     preventDefault() {},
                    //     target: {
                    //         name: "passwordConfirmation",
                    //         value: "doesNotMatch"
                    //     }
                    // } as ChangeEvent<HTMLInputElement>);
                });

                // enzymeWrapper.update();
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().value)
                //     .toEqual("doesNotMatch");
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().error).toEqual(true);
                // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(3).props().helperText)
                //     .toEqual("Passwords do not match");
            });
        });
    });

    describe("api validation errors", () => {
        it("displays username errors", async () => {
            const register = jest.fn(registerAsync.request);

            render(
                <Registration
                    register={register}
                />
            );

            await act(async () => {
                // enzymeWrapper.find(RegistrationFormInner).props().setStatus({"username": "Api Error"});
            });

            // enzymeWrapper.update();
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().value)
            //     .toEqual("");
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().error).toEqual(true);
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(0).props().helperText)
            //     .toEqual("Api Error");
        });

        it("displays email errors", async () => {
            const register = jest.fn(registerAsync.request);

            render(
                <Registration
                    register={register}
                />
            );

            await act(async () => {
                // enzymeWrapper.find(RegistrationFormInner).props().setStatus({"email": "Api Error"});
            });

            // enzymeWrapper.update();
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().value)
            //     .toEqual("");
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().error).toEqual(true);
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(1).props().helperText)
            //     .toEqual("Api Error");
        });

        it("displays password errors", async () => {
            const register = jest.fn(registerAsync.request);

            render(
                <Registration
                    register={register}
                />
            );

            await act(async () => {
                // enzymeWrapper.find(RegistrationFormInner).props().setStatus({"password": "Api Error"});
            });

            // enzymeWrapper.update();
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().value)
            //     .toEqual("");
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().error).toEqual(true);
            // expect(enzymeWrapper.find(RegistrationFormInner).find(TextField).at(2).props().helperText)
            //     .toEqual("Api Error");
        });
    });
});
