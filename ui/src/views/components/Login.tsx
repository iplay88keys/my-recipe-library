import { Button, Container, CssBaseline, TextField, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import { FormikBag, FormikProps, FormikTouched, withFormik } from "formik";
import React from "react";
import * as Yup from "yup";
import { loginAsync } from "../../state/ducks/users/actions";
import { LoginRequest } from "../../state/ducks/users/types";

const useStyles = makeStyles((theme) => ({
    paper: {
        marginTop: theme.spacing(8),
        display: "flex",
        flexDirection: "column",
        alignItems: "center"
    },
    form: {
        width: "100%",
        marginTop: theme.spacing(1)
    },
    submit: {
        margin: theme.spacing(3, 0, 2)
    }
}));

export interface LoginFormValues {
    login: string,
    password: string,
    doLogin: typeof loginAsync.request
}

const handleSubmit = (values: LoginFormValues, props: FormikBag<LoginFormProps, LoginFormValues>) => {
    const {doLogin} = values;
    if (values.login && values.password) {
        const request: LoginRequest = {
            login: values.login,
            password: values.password
        };

        props.setStatus({});
        void doLogin(request, props.setStatus);
    }

    props.setSubmitting(false);
    let newTouched = {} as FormikTouched<LoginFormValues>;
    Object.keys(values).forEach(key => {
        newTouched = {...newTouched, [key]: false};
    });

    void props.setTouched(newTouched);
};

export const LoginFormInner = (props: FormikProps<LoginFormValues>) => {
    const {handleSubmit, getFieldProps, isSubmitting, touched, errors} = props;

    const classes = useStyles();

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline/>
            <div className={classes.paper}>
                <Typography component="h1" variant="h5">
                    Login
                </Typography>
                <form
                    aria-label="Login"
                    name="loginForm"
                    className={classes.form}
                    onSubmit={handleSubmit}
                >
                    <TextField
                        data-testid="loginSection"
                        aria-label="loginName"
                        placeholder="Username or Email Address"
                        variant="outlined"
                        label="Username/Email Address"
                        margin="normal"
                        error={touched.login && Boolean(errors.login)}
                        helperText={touched.login && errors.login}
                        {...getFieldProps("login")}
                        required
                        fullWidth
                    />
                    <TextField
                        data-testid="passwordSection"
                        aria-label="password"
                        type="password"
                        placeholder="Password"
                        variant="outlined"
                        label="Password"
                        margin="normal"
                        error={touched.password && Boolean(errors.password)}
                        helperText={touched.password && errors.password}
                        {...getFieldProps("password")}
                        required
                        fullWidth
                    />
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        color="primary"
                        disabled={isSubmitting}
                        className={classes.submit}
                    >
                        Login
                    </Button>
                </form>
            </div>
        </Container>
    );
};

interface LoginFormProps {
    login: typeof loginAsync.request
}

export default withFormik<LoginFormProps, LoginFormValues>({
    mapPropsToValues: (props: LoginFormProps): LoginFormValues => ({
        login: "",
        password: "",
        doLogin: props.login
    }),
    validationSchema: Yup.object({
        login: Yup.string()
                      .required("Required"),
        password: Yup.string()
                     .required("Required")
    }),
    handleSubmit: handleSubmit
})(LoginFormInner);
