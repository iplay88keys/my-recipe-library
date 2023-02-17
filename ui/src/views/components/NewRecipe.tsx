import { Button, Container, CssBaseline, TextField, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import { FormikHelpers, FormikProps, FormikTouched, withFormik } from "formik";
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

// const showError = (field: string, formikProps: FormikProps<NewRecipeFormValues>): boolean => {
//     return (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) ||
//         (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field));
// };
//
// const errorMessage = (field: string, formikProps: FormikProps<NewRecipeFormValues>): string => {
//     if (!getIn(formikProps.touched, field) && !!formikProps.status && !!getIn(formikProps.status, field)) {
//         return getIn(formikProps.status, field);
//     } else if (!!getIn(formikProps.touched, field) && !!getIn(formikProps.errors, field)) {
//         return getIn(formikProps.errors, field);
//     } else {
//         return "";
//     }
// };

const handleSubmit = (values: NewRecipeFormValues, props: FormikHelpers<NewRecipeFormValues>) => {
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
        doCreate(recipe, props.setStatus);
    }

    props.setSubmitting(false);
    let newTouched = {} as FormikTouched<NewRecipeFormValues>;
    Object.keys(values).forEach(key => {
        newTouched = {...newTouched, [key]: false};
    });

    props.setTouched(newTouched);
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
                        placeholder="Name"
                        variant="outlined"
                        label="Name"
                        margin="normal"
                        error={touched.name && Boolean(errors.name)}
                        helperText={touched.name && errors.name}
                        {...getFieldProps("name")}
                        required
                        fullWidth
                    />
                    <TextField
                        placeholder="Description"
                        variant="outlined"
                        label="Description"
                        margin="normal"
                        error={touched.description && Boolean(errors.description)}
                        helperText={touched.description && errors.description}
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
                        error={touched.servings && Boolean(errors.servings)}
                        helperText={touched.servings && errors.servings}
                        {...getFieldProps("servings")}
                        required
                        fullWidth
                    />
                    <TextField
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
