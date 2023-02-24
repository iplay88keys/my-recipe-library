import { Button, Container, CssBaseline, TextField, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import { FormikHelpers, FormikProps, FormikTouched, getIn, withFormik } from "formik";
import React from "react";
import * as Yup from "yup";
import { createRecipeAsync } from "../../state/ducks/recipes/actions";
import { RecipeCreateRequest } from "../../state/ducks/recipes/types";

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

export interface NewRecipeFormValues {
    name: string,
    description: string,
    servings: number
    prep_time: string
    cook_time: string
    cool_time: string
    total_time: string
    source: string
    doCreate: typeof createRecipeAsync.request
}

const showError = (field: string, formikProps: FormikProps<NewRecipeFormValues>): boolean => {
    return (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) ||
        (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field));
};

const errorMessage = (field: string, formikProps: FormikProps<NewRecipeFormValues>): string => {
    if (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) {
        return getIn(formikProps.status, field);
    } else if (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field)) {
        return getIn(formikProps.errors, field);
    } else {
        return "";
    }
};

const handleSubmit = (values: NewRecipeFormValues, props: FormikHelpers<NewRecipeFormValues>) => {
    const {doCreate} = values;
    if (values.name && values.description && values.servings) {
        const recipe: RecipeCreateRequest = {
            name: values.name,
            description: values.description,
            servings: values.servings,
            prep_time: values.prep_time ? values.prep_time : undefined,
            cook_time: values.cook_time ? values.cook_time : undefined,
            cool_time: values.cool_time ? values.cool_time : undefined,
            total_time: values.total_time ? values.total_time : undefined,
            source: values.source ? values.source : undefined
        };

        props.setStatus({});
        doCreate(recipe, props.setStatus);
    }

    props.setSubmitting(false);
    let newTouched = {} as FormikTouched<NewRecipeFormValues>;
    Object.keys(values).map(key => {
        newTouched = {...newTouched, [key]: false};
    });

    props.setTouched(newTouched);
};

export const NewRecipeFormInner = (props: FormikProps<NewRecipeFormValues>) => {
    const {handleSubmit, getFieldProps, isSubmitting} = props;

    const classes = useStyles();

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline/>
            <div className={classes.paper}>
                <Typography component="h1" variant="h5">
                    New Recipe
                </Typography>
                <form onSubmit={handleSubmit} className={classes.form}>
                    <TextField
                        placeholder="Name"
                        variant="outlined"
                        label="Name"
                        margin="normal"
                        error={showError("name", props)}
                        helperText={errorMessage("name", props)}
                        {...getFieldProps("name")}
                        required
                        fullWidth
                    />
                    <TextField
                        placeholder="Description"
                        variant="outlined"
                        label="Description"
                        margin="normal"
                        error={showError("description", props)}
                        helperText={errorMessage("description", props)}
                        {...getFieldProps("description")}
                        required
                        fullWidth
                    />
                    <TextField
                        type="number"
                        placeholder="Servings"
                        variant="outlined"
                        label="Servings"
                        margin="normal"
                        error={showError("servings", props)}
                        helperText={errorMessage("servings", props)}
                        {...getFieldProps("servings")}
                        required
                        fullWidth
                    />
                    <TextField
                        placeholder="Prep Time"
                        variant="outlined"
                        label="Prep Time"
                        margin="normal"
                        error={showError("prep_time", props)}
                        helperText={errorMessage("prep_time", props)}
                        {...getFieldProps("prep_time")}
                        fullWidth
                    />
                    <TextField
                        placeholder="Cook Time"
                        variant="outlined"
                        label="Cook Time"
                        margin="normal"
                        error={showError("cook_time", props)}
                        helperText={errorMessage("cook_time", props)}
                        {...getFieldProps("cook_time")}
                        fullWidth
                    />
                    <TextField
                        placeholder="Cool Time"
                        variant="outlined"
                        label="Cool Time"
                        margin="normal"
                        error={showError("cool_time", props)}
                        helperText={errorMessage("cool_time", props)}
                        {...getFieldProps("cool_time")}
                        fullWidth
                    />
                    <TextField
                        placeholder="Total Time"
                        variant="outlined"
                        label="Total Time"
                        margin="normal"
                        error={showError("total_time", props)}
                        helperText={errorMessage("total_time", props)}
                        {...getFieldProps("total_time")}
                        fullWidth
                    />
                    <TextField
                        type="source"
                        placeholder="Source"
                        variant="outlined"
                        label="Source"
                        margin="normal"
                        error={showError("source", props)}
                        helperText={errorMessage("source", props)}
                        {...getFieldProps("source")}
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
                        Create
                    </Button>
                </form>
            </div>
        </Container>
    );
};

interface NewRecipeFormProps {
    create: typeof createRecipeAsync.request
}

export default withFormik<NewRecipeFormProps, NewRecipeFormValues>({
    mapPropsToValues: (props: NewRecipeFormProps): NewRecipeFormValues => ({
        name: "",
        description: "",
        servings: 0,
        prep_time: "",
        cook_time: "",
        cool_time: "",
        total_time: "",
        source: "",
        doCreate: props.create
    }),
    validationSchema: Yup.object({
        name: Yup.string().required("Required"),
        description: Yup.string().required("Required"),
        servings: Yup.string().required("Required")
    }),
    handleSubmit: handleSubmit
})(NewRecipeFormInner);
