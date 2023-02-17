import React from "react";
import { connect } from "react-redux";
import { RouteComponentProps } from "react-router";
import { ApplicationState } from "../../state/ducks";
import { fetchRecipeAsync } from "../../state/ducks/recipes/actions";
import { RecipeResponse } from "../../state/ducks/recipes/types";
import { Recipe } from "../components/Recipe";

interface PropsFromState {
    recipe: RecipeResponse
    loading: boolean
}

interface PropsFromDispatch {
    fetchRecipe: typeof fetchRecipeAsync.request
}

interface MatchParams extends RouteComponentProps<{ recipeID: string }> {}

interface State {
    currentRecipeID: string
}

type AllProps = PropsFromState & PropsFromDispatch & MatchParams

class RecipePage extends React.Component<AllProps, State> {
    componentDidMount() {
        const {fetchRecipe} = this.props;

        // let { recipeID } = useParams();
        fetchRecipe(+this.props.match.params.recipeID);
    }

    UNSAFE_componentWillReceiveProps(nextProps: AllProps) {
        if (this.props.match.params.recipeID !== nextProps.match.params.recipeID) {
            this.props.fetchRecipe(+nextProps.match.params.recipeID);
        }
    }

    render() {
        return (
            <Recipe
                recipe={this.props.recipe}
                loading={this.props.loading}
            />
        );
    }
}

const mapStateToProps = ({recipes}: ApplicationState) => ({
    recipe: recipes.recipe,
    loading: recipes.loading
});

const mapDispatchToProps = {
    fetchRecipe: fetchRecipeAsync.request
};

export default connect(mapStateToProps, mapDispatchToProps)(RecipePage);
