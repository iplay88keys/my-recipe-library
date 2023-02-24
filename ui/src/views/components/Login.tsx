import { Button, Container, CssBaseline, TextField, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import { FormikHelpers, FormikProps, FormikTouched, getIn, withFormik } from "formik";
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

const showError = (field: string, formikProps: FormikProps<LoginFormValues>): boolean => {
    return (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) ||
        (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field));
};

const errorMessage = (field: string, formikProps: FormikProps<LoginFormValues>): string => {
    if (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) {
        return getIn(formikProps.status, field);
    } else if (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field)) {
        return getIn(formikProps.errors, field);
    } else {
        return "";
    }
};

const handleSubmit = (values: LoginFormValues, props: FormikHelpers<LoginFormValues>) => {
    const {doLogin} = values;
    if (values.login && values.password) {
        const user: LoginRequest = {
            login: values.login,
            password: values.password
        };

        props.setStatus({});
        doLogin(user, props.setStatus);
    }

    props.setSubmitting(false);
    let newTouched = {} as FormikTouched<LoginFormValues>;
    Object.keys(values).map(key => {
        newTouched = {...newTouched, [key]: false};
    });

    props.setTouched(newTouched);
};

export const LoginFormInner = (props: FormikProps<LoginFormValues>) => {
    const {handleSubmit, getFieldProps, isSubmitting} = props;

    const classes = useStyles();

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline/>
            <div className={classes.paper}>
                <Typography component="h1" variant="h5">
                    Login
                </Typography>
                <form onSubmit={handleSubmit} className={classes.form}>
                    <TextField
                        placeholder="Username or Email Address"
                        variant="outlined"
                        label="Username/Email Address"
                        margin="normal"
                        error={showError("login", props)}
                        helperText={errorMessage("login", props)}
                        {...getFieldProps("login")}
                        required
                        fullWidth
                    />
                    <TextField
                        type="password"
                        placeholder="Password"
                        variant="outlined"
                        label="Password"
                        margin="normal"
                        error={showError("password", props)}
                        helperText={errorMessage("password", props)}
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
