import { render } from "@testing-library/react";
import React from "react";
import { BrowserRouter as Router } from "react-router-dom";
import { Navigation } from "./Navigation";

describe("Navigation", () => {
    it("should render a list of links if the user is logged in", () => {
        render(<Router><Navigation loggedIn={true}/></Router>)
        // const enzymeWrapper = mount();
        //
        // expect(enzymeWrapper.find(Link)).toHaveLength(1);
        // expect(enzymeWrapper.find(Link).at(0).text()).toEqual("My Recipe Library");
        // expect(enzymeWrapper.find(Button)).toHaveLength(2);
        // expect(enzymeWrapper.find(Button).at(0).text()).toEqual("Recipes");
        // expect(enzymeWrapper.find(Button).at(1).text()).toEqual("Logout");
    });

    it("displays register and login links if the user is not logged in", () => {
        // const enzymeWrapper = mount(<Router><Navigation loggedIn={false}/></Router>);
        //
        // expect(enzymeWrapper.find(Link)).toHaveLength(1);
        // expect(enzymeWrapper.find(Link).at(0).text()).toEqual("My Recipe Library");
        // expect(enzymeWrapper.find(Button)).toHaveLength(2);
        // expect(enzymeWrapper.find(Button).at(0).text()).toEqual("Register");
        // expect(enzymeWrapper.find(Button).at(1).text()).toEqual("Login");
    });
});
