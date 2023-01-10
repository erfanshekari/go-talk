import { default as GTCTypes } from "./types"


interface StateI {
    state: GTCTypes.State
    setState: (f: GTCTypes.SetState) => void
}

class State implements StateI {

    state: GTCTypes.State

    private observer?: GTCTypes.StateObserverFunc

    constructor(config: GTCTypes.Config) {
        this.state = {
            initializing: true,
            closed: false,
            connected: false,
            connecting: false,
            initialized: false,
            authenticated: false,
        }
        this.observer = config.onStateChange
    }

    setState(f: GTCTypes.SetState) {
        this.state = f(this.state);
        if (this.observer) {
            this.observer(this.state);
        }
    }
}

export default State