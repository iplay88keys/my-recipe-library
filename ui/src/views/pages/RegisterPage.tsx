import React from "react";
import { connect } from "react-redux";
import { ApplicationState } from "../../state/ducks";
import { registerAsync } from "../../state/ducks/users/actions";
import Registration from "../components/Registration";

interface PropsFromDispatch {
    register: typeof registerAsync.request
}

interface State {}

type AllProps = PropsFromDispatch & State

class RegisterPage extends React.Component<AllProps, State> {
    render() {
        return (
            <Registration
                register={this.props.register}
            />
        );
    }
}

const mapStateToProps = (_: ApplicationState) => ({});

const mapDispatchToProps = {
    register: registerAsync.request
};

export default connect(mapStateToProps, mapDispatchToProps)(RegisterPage);
