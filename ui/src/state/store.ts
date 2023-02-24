import { createRootReducer, rootSaga } from "./ducks";
import { applyMiddleware, createStore } from "redux";
import { composeWithDevTools } from "redux-devtools-extension";
import createSagaMiddleware from "redux-saga";

export default function configureStore() {
    const sagaMiddleware = createSagaMiddleware();

    const store = createStore(
        createRootReducer(),
        composeWithDevTools(
            applyMiddleware(sagaMiddleware)
        )
    );

    sagaMiddleware.run(rootSaga);

    return store;
}
