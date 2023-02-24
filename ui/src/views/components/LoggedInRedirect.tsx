import * as React from "react";
import { Redirect, Route, RouteComponentProps } from "react-router";

interface PrivateRouteProps {
    component: React.ComponentType<any>,
    exact?: boolean;
    path: string;
}

export const LoggedInRedirect = ({component: Component, ...rest}: PrivateRouteProps) => (
    <Route {...rest} render={(props: RouteComponentProps) => (
        localStorage.getItem("access_token")
            ? <Redirect to={{pathname: "/", state: {from: props.location}}}/>
            : <Component {...props} />
    )}/>
);
