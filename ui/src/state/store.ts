import { configureStore } from "@reduxjs/toolkit";
import { reducer, rootSaga } from "./ducks";
import createSagaMiddleware from "redux-saga";

export default function setupStore() {
    const sagaMiddleware = createSagaMiddleware();
    const middleware = [sagaMiddleware]

    const store = configureStore({
        reducer,
        middleware: (getDefaultMiddleware) =>
            getDefaultMiddleware().concat(middleware)
    });

    sagaMiddleware.run(rootSaga);

    return store;
}
