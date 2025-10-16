import { Button, Container, CssBaseline, TextField, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import { FormikBag, FormikProps, FormikTouched, withFormik } from "formik";
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

const handleSubmit = (values: NewRecipeFormValues, props: FormikBag<NewRecipeFormProps, NewRecipeFormValues>) => {
    const {doCreate} = values;
    if (values.name && values.description && values.servings) {
        const recipe: RecipeCreateRequest = {
            name: values.name,
            description: values.description,
            servings: values.servings,
            prep_time: values.prep_time || undefined,
            cook_time: values.cook_time || undefined,
            cool_time: values.cool_time || undefined,
            total_time: values.total_time || undefined,
            source: values.source || undefined
        };

        props.setStatus({});
        void doCreate(recipe, props.setStatus);
    }

    props.setSubmitting(false);
    let newTouched = {} as FormikTouched<NewRecipeFormValues>;
    Object.keys(values).forEach(key => {
        newTouched = {...newTouched, [key]: false};
    });

    void props.setTouched(newTouched);
};

const showError = (field: string, formikProps: FormikProps<NewRecipeFormValues>): boolean => {
    const status = formikProps.status as NewRecipeFormValues | undefined;

    const fieldTouched = !!formikProps.touched[field as keyof NewRecipeFormValues];
    const hasFormikError = Boolean(formikProps.errors[field as keyof NewRecipeFormValues]);
    const hasStatusError = !!(status && status[field as keyof NewRecipeFormValues]);

    return (fieldTouched && hasFormikError) || (!fieldTouched && hasStatusError)
};

const errorMessage = (field: string, formikProps: FormikProps<NewRecipeFormValues>): string => {
    const status = formikProps.status as NewRecipeFormValues | undefined;
    const fieldTouched = !!formikProps.touched[field as keyof NewRecipeFormValues];
    const formikError = formikProps.errors[field as keyof NewRecipeFormValues];
    const statusError = status && status[field as keyof NewRecipeFormValues];

    if (!fieldTouched && typeof statusError == "string" && Boolean(statusError)) {
        return statusError
    } else if (fieldTouched && typeof formikError == "string" && Boolean(formikError)) {
        return formikError
    } else {
        return ""
    }
};

export const NewRecipeFormInner = (props: FormikProps<NewRecipeFormValues>) => {
    const {handleSubmit, getFieldProps, isSubmitting, touched, errors} = props;

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
                        data-testid="nameSection"
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
                        data-testid="descriptionSection"
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
                        data-testid="servingsSection"
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
                        data-testid="prepTimeSection"
                        placeholder="Prep Time"
                        variant="outlined"
                        label="Prep Time"
                        margin="normal"
                        error={touched.prep_time && Boolean(errors.prep_time)}
                        helperText={touched.prep_time && errors.prep_time}
                        {...getFieldProps("prep_time")}
                        fullWidth
                    />
                    <TextField
                        data-testid="cookTimeSection"
                        placeholder="Cook Time"
                        variant="outlined"
                        label="Cook Time"
                        margin="normal"
                        error={touched.cook_time && Boolean(errors.cook_time)}
                        helperText={touched.cook_time && errors.cook_time}
                        {...getFieldProps("cook_time")}
                        fullWidth
                    />
                    <TextField
                        data-testid="coolTimeSection"
                        placeholder="Cool Time"
                        variant="outlined"
                        label="Cool Time"
                        margin="normal"
                        error={touched.cool_time && Boolean(errors.cool_time)}
                        helperText={touched.cool_time && errors.cool_time}
                        {...getFieldProps("cool_time")}
                        fullWidth
                    />
                    <TextField
                        data-testid="totalTimeSection"
                        placeholder="Total Time"
                        variant="outlined"
                        label="Total Time"
                        margin="normal"
                        error={touched.total_time && Boolean(errors.total_time)}
                        helperText={touched.total_time && errors.total_time}
                        {...getFieldProps("total_time")}
                        fullWidth
                    />
                    <TextField
                        data-testid="sourceSection"
                        type="source"
                        placeholder="Source"
                        variant="outlined"
                        label="Source"
                        margin="normal"
                        error={touched.source && Boolean(errors.source)}
                        helperText={touched.source && errors.source}
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
