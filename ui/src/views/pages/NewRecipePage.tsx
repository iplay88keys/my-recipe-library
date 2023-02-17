import React from "react";
import { connect } from "react-redux";
import { ApplicationState } from "../../state/ducks";
import { createRecipeAsync } from "../../state/ducks/recipes/actions";
import NewRecipe from "../components/NewRecipe";

interface PropsFromDispatch {
    create: typeof createRecipeAsync.request
}

interface State {}

type AllProps = PropsFromDispatch & State

class NewRecipePage extends React.Component<AllProps, State> {
    render() {
        return (
            <NewRecipe
                create={this.props.create}
            />
        );
    }
}

const mapStateToProps = (_: ApplicationState) => ({});

const mapDispatchToProps = {
    create: createRecipeAsync.request
};

export default connect(mapStateToProps, mapDispatchToProps)(NewRecipePage);
