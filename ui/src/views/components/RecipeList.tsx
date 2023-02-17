import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@material-ui/core";
import React from "react";
import { RouteComponentProps } from "react-router";
import { withRouter } from "react-router-dom";
import styled from "styled-components";
import { RecipeResponse } from "../../state/ducks/recipes/types";

export const StyledTableRow = styled(TableRow)`
    cursor: pointer;
`;

interface RecipeListProps extends RouteComponentProps {
    recipes: RecipeResponse[]
    loading: boolean
}

export const RecipeList = ({recipes, loading, history}: RecipeListProps) => {
    if (loading) {
        return (
            <div>
                <p>Loading recipes</p>
            </div>
        );
    }

    if (!recipes) {
        return (
            <div>
                <p>No recipes yet! Add some!</p>
            </div>
        );
    }

    return (
        <TableContainer component={Paper}>
            <Table aria-label="recipes">
                <TableHead>
                    <TableRow>
                        <TableCell>Recipe</TableCell>
                        <TableCell>Description</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {recipes.map((recipe: RecipeResponse) =>
                        <StyledTableRow key={recipe.id} onClick={() => history.push(`/recipes/${recipe.id}`)}>
                            <TableCell>{recipe.name}</TableCell>
                            {recipe.description != null && <TableCell>{recipe.description}</TableCell>}
                        </StyledTableRow>
                    )}
                </TableBody>
            </Table>
        </TableContainer>
    );
};

export default withRouter(RecipeList);
