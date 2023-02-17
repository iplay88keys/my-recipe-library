import * as React from "react";
import { Redirect, Route } from "react-router";
import { RouteProps } from "react-router-dom";

interface PrivateRouteProps extends RouteProps {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    component: React.ComponentType<any>,
    exact?: boolean;
    path: string;
}

export const LoggedInRedirect = ({component: Component, ...rest}: PrivateRouteProps) => (
    <Route {...rest} render={props => (
        localStorage.getItem("access_token")
            ? <Redirect to={{pathname: "/", state: {from: props.location}}}/>
            : <Component {...props} />
    )}/>
);
