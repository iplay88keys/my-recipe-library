import { createLocation, createMemoryHistory, Location, MemoryHistory } from "history";
import React from "react";
import { match } from "react-router";
import { RecipeResponse } from "../../state/ducks/recipes/types";
import { RecipeList } from "./RecipeList";
import { render, screen, within } from '@testing-library/react'
import { Simulate } from "react-dom/test-utils";

describe("RecipeListResponse", () => {
    let history: MemoryHistory;
    let matchParam: match<{ id: string }>;
    let location: Location;

    beforeEach(() => {
        history = createMemoryHistory();
        const path = `/route/:id`;

        matchParam = {
            isExact: false,
            path,
            url: path.replace(":id", "1"),
            params: {id: "1"}
        };

        location = createLocation(matchParam.url);
    });

    it("should render a list of recipes", () => {
        const recipes = [{
            id: 0,
            name: "First",
            description: "One"
        }, {
            id: 1,
            name: "Second",
            description: "Two"
        }] as RecipeResponse[];

        render(
            <RecipeList
                recipes={recipes}
                loading={false}
                history={history}
                match={matchParam}
                location={location}
            />
        );

        const table = screen.getByRole("table", {name: /recipes/i});
        const rows = within(table).getAllByRole("row")

        expect(rows.length).toEqual(3)
        let cells = within(rows[0]).getAllByRole("columnheader")
        expect(cells[0]).toHaveTextContent("Recipe")
        expect(cells[1]).toHaveTextContent("Description")

        cells = within(rows[1]).getAllByRole("cell")
        expect(cells.length).toEqual(2)
        expect(cells[0]).toHaveTextContent("First")
        expect(cells[1]).toHaveTextContent("One")

        cells = within(rows[2]).getAllByRole("cell")
        expect(cells.length).toEqual(2)
        expect(cells[0]).toHaveTextContent("Second")
        expect(cells[1]).toHaveTextContent("Two")
    });

    it("does not render missing data", () => {
        const recipes = [{
            id: 0,
            name: "First"
        }] as RecipeResponse[];

        render(
            <RecipeList
                recipes={recipes}
                loading={false}
                history={history}
                match={matchParam}
                location={location}
            />
        );

        const table = screen.getByRole("table", {name: /recipes/i});
        const rows = within(table).getAllByRole("row")

        expect(rows.length).toEqual(2)

        let cells = within(rows[0]).getAllByRole("columnheader")
        expect(cells[0]).toHaveTextContent("Recipe")
        expect(cells[1]).toHaveTextContent("Description")

        cells = within(rows[1]).getAllByRole("cell")
        expect(cells.length).toEqual(1)
        expect(cells[0]).toHaveTextContent("First")
    });

    it("should load the single recipe page when the row is clicked", () => {
        const recipes = [{
            id: 0,
            name: "First",
            description: "One"
        }, {
            id: 1,
            name: "Second",
            description: "Two"
        }] as RecipeResponse[];

        const historyMock = history;
        historyMock.push = jest.fn();

        render(
            <RecipeList
                recipes={recipes}
                loading={false}
                history={historyMock}
                match={matchParam}
                location={location}
            />
        );

        const table = screen.getByRole("table", {name: /recipes/i});
        const rows = within(table).getAllByRole("row");

        Simulate.click(rows[1]);
        expect(historyMock.push).toHaveBeenCalledWith("/recipes/0");

        historyMock.push = jest.fn();
        Simulate.click(rows[2]);
        expect(historyMock.push).toHaveBeenCalledWith("/recipes/1");
    });

    it("should render loading info when loading", () => {
        const props = {
            recipes: [] as RecipeResponse[]
        };

        render(
            <RecipeList
                recipes={props.recipes}
                loading={true}
                history={history}
                match={matchParam}
                location={location}
            />
        );

        expect(screen.getByText("Loading recipes")).toBeInTheDocument();
    });
});
