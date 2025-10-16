import { CssBaseline } from "@material-ui/core";
import { createTheme } from '@material-ui/core/styles'
import { ThemeProvider } from "@material-ui/styles";
import React from "react";
import { connect } from "react-redux";
import { Redirect, Route, Router, Switch } from "react-router-dom";
import styled from "styled-components";
import { history } from "../helpers/history";
import { ApplicationState } from "../state/ducks";
import { LoggedInRedirect } from "./components/LoggedInRedirect";
import { Navigation } from "./components/Navigation";
import { PrivateRoute } from "./components/PrivateRoute";
import Login from "./pages/LoginPage";
import NewRecipePage from "./pages/NewRecipePage";
import RecipePage from "./pages/RecipePage";
import RecipesPage from "./pages/RecipesPage";
import RegisterPage from "./pages/RegisterPage";

const StyledApp = styled.div`
  height: 100%;
  width: 75%;
  margin: auto;
  padding-top: 30px;
`;

interface PropsFromState {
    loggedIn: boolean
}

interface PropsFromDispatch {}

interface State {}

type AllProps = PropsFromState & PropsFromDispatch & State

class App extends React.Component<AllProps, State> {
    // constructor(props: AllProps) {
    //     super(props);
    //
    //     history.listen(() => {
    //     });
    // }

    render() {
        const theme = createTheme({
            palette: {
                type: "light"
            }
        });

        return (
            <Router history={history}>
                <ThemeProvider theme={theme}>
                    <CssBaseline/>
                    <div>
                        <Navigation loggedIn={this.props.loggedIn}/>
                        <StyledApp>
                            <Switch>
                                <Route exact path="/" component={RecipesPage}/>
                                <LoggedInRedirect exact path="/register" component={RegisterPage}/>
                                <Route exact path="/login" component={Login}/>
                                <PrivateRoute exact path="/recipes/new" component={NewRecipePage}/>
                                <PrivateRoute exact path="/recipes" component={RecipesPage}/>
                                <PrivateRoute exact path="/recipes/:recipeID" component={RecipePage}/>
                                <Redirect from="*" to="/"/>
                            </Switch>
                        </StyledApp>
                    </div>
                </ThemeProvider>
            </Router>
        );
    }
}

const mapStateToProps = ({users}: ApplicationState) => ({
    loggedIn: users.loggedIn
});

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(App);
