import { Button, Container, CssBaseline, TextField, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import { FormikHelpers, FormikProps, FormikTouched, withFormik } from "formik";
import React from "react";
import * as Yup from "yup";
import { registerAsync } from "../../state/ducks/users/actions";
import { RegisterRequest } from "../../state/ducks/users/types";

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

export interface RegistrationFormValues {
    username: string,
    email: string,
    password: string,
    passwordConfirmation: string
    doRegister: typeof registerAsync.request
}

// const showError = (field: string, formikProps: FormikProps<RegistrationFormValues>): boolean => {
//     return (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) ||
//         (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field));
// };
//
// const errorMessage = (field: string, formikProps: FormikProps<RegistrationFormValues>): string => {
//     if (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) {
//         return getIn(formikProps.status, field);
//     } else if (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field)) {
//         return getIn(formikProps.errors, field);
//     } else {
//         return "";
//     }
// };

const handleSubmit = (values: RegistrationFormValues, props: FormikHelpers<RegistrationFormValues>) => {
    const {doRegister} = values;
    if (values.username && values.email && values.password) {
        const user: RegisterRequest = {
            username: values.username,
            email: values.email,
            password: values.password
        };

        props.setStatus({});
        void doRegister(user, props.setStatus);
    }

    props.setSubmitting(false);
    let newTouched = {} as FormikTouched<RegistrationFormValues>;
    Object.keys(values).forEach(key => {
        newTouched = {...newTouched, [key]: false};
    });

    void props.setTouched(newTouched);
};

export const RegistrationFormInner = (props: FormikProps<RegistrationFormValues>) => {
    const {handleSubmit, getFieldProps, isSubmitting, touched, errors} = props;

    const classes = useStyles();

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline/>
            <div className={classes.paper}>
                <Typography component="h1" variant="h5">
                    Register
                </Typography>
                <form onSubmit={handleSubmit} className={classes.form}>
                    <TextField
                        type="username"
                        placeholder="Username"
                        variant="outlined"
                        label="Username"
                        margin="normal"
                        error={touched.username && Boolean(errors.username)}
                        helperText={touched.username && errors.username}
                        {...getFieldProps("username")}
                        required
                        fullWidth
                    />
                    <TextField
                        type="email"
                        placeholder="Email"
                        variant="outlined"
                        label="Email"
                        margin="normal"
                        error={touched.email && Boolean(errors.email)}
                        helperText={touched.email && errors.email}
                        {...getFieldProps("email")}
                        required
                        fullWidth
                    />
                    <TextField
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
                    <TextField
                        type="password"
                        placeholder="Confirm Your Password"
                        variant="outlined"
                        label="Confirm Password"
                        margin="normal"
                        error={touched.passwordConfirmation && Boolean(errors.passwordConfirmation)}
                        helperText={touched.passwordConfirmation && errors.passwordConfirmation}
                        {...getFieldProps("passwordConfirmation")}
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
                        Register
                    </Button>
                </form>
            </div>
        </Container>
    );
};

interface RegistrationFormProps {
    register: typeof registerAsync.request
}

export default withFormik<RegistrationFormProps, RegistrationFormValues>({
    mapPropsToValues: (props: RegistrationFormProps): RegistrationFormValues => ({
        username: "",
        email: "",
        password: "",
        passwordConfirmation: "",
        doRegister: props.register
    }),
    validationSchema: Yup.object({
        username: Yup.string()
            .min(6, "Must be 6 characters or more")
            .max(30, "Must be 30 characters or less")
            .matches(/^[a-zA-Z][\w]*/, "Must start with a letter and only contain alphanumeric characters and underscores (_)")
            .required("Required"),
        email: Yup.string()
            .matches(/.*@.*\..*/, "Must be a valid email")
            .required("Required"),
        password: Yup.string()
            .min(6, "Must be 6 characters or more")
            .max(64, "Must be 64 characters or less")
            .matches(/[\d\w\W]*/, "Must contain at least one lowercase, uppercase, number, and special character")
            .required("Required"),
        passwordConfirmation: Yup.string()
            .oneOf([Yup.ref("password")], "Passwords do not match")
    }),
    handleSubmit: handleSubmit
})(RegistrationFormInner);
