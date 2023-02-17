import React from "react";
import { connect } from "react-redux";
import { loginAsync, logout } from "../../state/ducks/users/actions";
import Login from "../components/Login";

interface PropsFromDispatch {
    login: typeof loginAsync.request
    logout: typeof logout
}

interface State {}

type AllProps = PropsFromDispatch & State

class LoginPage extends React.Component<AllProps, State> {
    constructor(props: AllProps) {
        super(props);

        this.props.logout();
    }

    render() {
        return (
            <Login
                login={this.props.login}
            />
        );
    }
}

const mapStateToProps = () => ({});

const mapDispatchToProps = {
    login: loginAsync.request,
    logout: logout
};

export default connect(mapStateToProps, mapDispatchToProps)(LoginPage);
