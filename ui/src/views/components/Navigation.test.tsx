import { render, screen, within } from "@testing-library/react";
import React from "react";
import { BrowserRouter as Router } from "react-router-dom";
import { Navigation } from "./Navigation";

describe("Navigation", () => {
    it("should render a list of links if the user is logged in", () => {
        render(<Router><Navigation loggedIn={true}/></Router>)

        expect(screen.getByRole("link")).toHaveTextContent("My Recipe Library");

        const links = screen.getAllByRole("button");
        expect(links).toHaveLength(2);
        expect(within(links[0]).getByText(/recipes/i)).toBeInTheDocument()
        expect(within(links[1]).getByText(/logout/i)).toBeInTheDocument()
    });

    it("displays register and login links if the user is not logged in", () => {
        render(<Router><Navigation loggedIn={false}/></Router>)

        expect(screen.getByRole("link")).toHaveTextContent("My Recipe Library");

        const links = screen.getAllByRole("button");
        expect(links).toHaveLength(2);
        expect(within(links[0]).getByText(/register/i)).toBeInTheDocument()
        expect(within(links[1]).getByText(/login/i)).toBeInTheDocument()
    });
});
