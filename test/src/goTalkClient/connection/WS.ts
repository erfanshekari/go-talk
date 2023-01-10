import { default as GTCTypes } from "./types"

interface Socket {
    connect: () => void
    reconnect: ()=> void
    close: () => void
}

class WS implements Socket {
    url: URL | string
    socket?: WebSocket

    constructor(config: GTCTypes.Config) {
        this.url = config.ws
    }

    connect() {

    }

    reconnect() {
        
    }

    close() {
        
    }

    private registerListeners() {

    }

    private removeListeners() {
        
    }

    private onOpen(instance: WS, event: MessageEvent<any>) {

    }

    private onMessage(instance: WS, event: MessageEvent<any>) {

    }

    private onError(instance: WS, event: MessageEvent<any>) {

    }

    private onClose(instance: WS, event: MessageEvent<any>) {

    }

}

export default WS