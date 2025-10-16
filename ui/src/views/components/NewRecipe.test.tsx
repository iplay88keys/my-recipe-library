import { act, render, screen, within } from "@testing-library/react";
import userEvent from '@testing-library/user-event'
import { FormikErrors } from "formik";
import React from "react";
import { createRecipeAsync } from "../../state/ducks/recipes/actions";
import { RecipeCreateRequest } from "../../state/ducks/recipes/types";
import { registerAsync } from "../../state/ducks/users/actions";
import NewRecipe, { NewRecipeFormValues } from "./NewRecipe";
import Registration from "./Registration";

describe("NewRecipe", () => {
    it("should render a form for creating a new recipe", async () => {
        const create = jest.fn(createRecipeAsync.request);

        render(
            <NewRecipe
                create={create}
            />
        );

        const user = userEvent.setup()

        await user.type(screen.getByPlaceholderText("Name"), "test-name")
        await user.type(screen.getByPlaceholderText("Description"), "test-description")
        await user.type(screen.getByPlaceholderText("Servings"), "1")
        await user.type(screen.getByPlaceholderText("Prep Time"), "1m")
        await user.type(screen.getByPlaceholderText("Cook Time"), "2m")
        await user.type(screen.getByPlaceholderText("Cool Time"), "3m")
        await user.type(screen.getByPlaceholderText("Total Time"), "6m")
        await user.type(screen.getByPlaceholderText("Source"), "test-source")

        const submitButton = screen.getByRole("button")
        expect(submitButton).toHaveTextContent("Create");
        await user.click(submitButton)

        expect(create).toHaveBeenCalledWith({
            "name": "test-name",
            "description": "test-description",
            "servings": 1,
            "prep_time": "1m",
            "cook_time": "2m",
            "cool_time": "3m",
            "total_time": "6m",
            "source": "test-source",
        }, expect.any(Function))
    });

    describe("form validation errors", () => {
        describe("name", () => {
            it("is required", async () => {
                const create = jest.fn(createRecipeAsync.request);

                render(
                    <NewRecipe
                        create={create}
                    />
                );

                const user = userEvent.setup()

                expect(screen.queryByText("Required")).not.toBeInTheDocument();

                await user.type(screen.getByPlaceholderText("Description"), "test-description")
                await user.type(screen.getByPlaceholderText("Servings"), "1")
                await user.type(screen.getByPlaceholderText("Prep Time"), "1m")
                await user.type(screen.getByPlaceholderText("Cook Time"), "2m")
                await user.type(screen.getByPlaceholderText("Cool Time"), "3m")
                await user.type(screen.getByPlaceholderText("Total Time"), "6m")
                await user.type(screen.getByPlaceholderText("Source"), "test-source")

                const submitButton = screen.getByRole("button")
                expect(submitButton).toHaveTextContent("Create");
                await user.click(submitButton)

                const nameSection = screen.getByTestId("nameSection")
                expect(within(nameSection).getByText("Required")).toBeInTheDocument()
            });
        });

        describe("description", () => {
            it("is required", async () => {
                const create = jest.fn(createRecipeAsync.request);

                render(
                    <NewRecipe
                        create={create}
                    />
                );

                const user = userEvent.setup()

                expect(screen.queryByText("Required")).not.toBeInTheDocument();

                await user.type(screen.getByPlaceholderText("Name"), "test-name")
                await user.type(screen.getByPlaceholderText("Servings"), "1")
                await user.type(screen.getByPlaceholderText("Prep Time"), "1m")
                await user.type(screen.getByPlaceholderText("Cook Time"), "2m")
                await user.type(screen.getByPlaceholderText("Cool Time"), "3m")
                await user.type(screen.getByPlaceholderText("Total Time"), "6m")
                await user.type(screen.getByPlaceholderText("Source"), "test-source")

                const submitButton = screen.getByRole("button")
                expect(submitButton).toHaveTextContent("Create");
                await user.click(submitButton)

                const descriptionSection = screen.getByTestId("descriptionSection")
                expect(within(descriptionSection).getByText("Required")).toBeInTheDocument()
            });
        });
    });

    describe("api validation errors", () => {
        it("displays errors", async () => {
            const create = jest.fn(createRecipeAsync.request);
            create.mockImplementation((payload: RecipeCreateRequest, meta: (errors: FormikErrors<NewRecipeFormValues>) => void): any => {
                meta({
                    name: "error 1",
                    description: "error 2",
                    servings: "error 3"
                } as FormikErrors<NewRecipeFormValues>);
            });

            render(
                <NewRecipe
                    create={create}
                />
            );

            const user = userEvent.setup()

            await user.type(screen.getByPlaceholderText("Name"), "test-name")
            await user.type(screen.getByPlaceholderText("Description"), "test-description")
            await user.type(screen.getByPlaceholderText("Servings"), "1")

            const submitButton = screen.getByRole("button")
            expect(submitButton).toHaveTextContent("Create");
            await user.click(submitButton)

            const nameSection = screen.getByTestId("nameSection");
            expect(within(nameSection).getByText("error 1")).toBeInTheDocument();

            const descriptionSection = screen.getByTestId("descriptionSection");
            expect(within(descriptionSection).getByText("error 2")).toBeInTheDocument();

            const servingsSection = screen.getByTestId("servingsSection");
            expect(within(servingsSection).getByText("error 3")).toBeInTheDocument();
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
