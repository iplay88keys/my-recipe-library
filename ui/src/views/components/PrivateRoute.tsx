import * as React from "react";
import { Component } from "react";
import { Redirect, Route } from "react-router";

interface PrivateRouteProps {
    component: React.ComponentType<any>,
    exact?: boolean;
    path: string;
}

export const PrivateRoute = ({component: Component, ...rest}: PrivateRouteProps) => (
    <Route {...rest} render={props => (
        localStorage.getItem("user")
            ? <Component {...props} />
            : <Redirect to={{pathname: "/login", state: {from: props.location}}}/>
    )}/>
);
