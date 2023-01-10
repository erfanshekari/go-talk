import GoTalkTypes from "./types"
import Encrypt from "./Encrypt";


class GoTalkBase {
    private encrypt: Encrypt // ill test this later
    private state: GoTalkTypes.State
    private config: GoTalkTypes.Config
    private socket?: WebSocket
    constructor(config: GoTalkTypes.Config) {
        this.config = config
        this.state = {
            closed: false,
            connected: false,
            connecting: false,
            initialized: false,
            isKeyExchangeDone: false,
            authenticated: false
        }
        let instance = this
        this.encrypt = new Encrypt({
            config: this.config,
            onSuccess: () => {
                instance.setState(state => ({...state, initialized: true}))
            },
            isExchangeDone: () => instance.state.isKeyExchangeDone
        })
    }

    setState(f:GoTalkTypes.SetState) {
        this.state = f(this.state)
        if (this.config.onStateChange) {
            this.config.onStateChange(this.state)
        }
    }

    setOnStateChangeFunction(f:(state:GoTalkTypes.State) => void) {
        this.config.onStateChange = f
        f(this.state)
    }

    async connect() {
        if (!this.state.initialized) {
            throw Error("You can't connect to server before GoTalkClient initialization.")
        }
        this.setState(state => ({...state, connecting: true}))
        this.socket = new WebSocket(this.config.ws)
        this.registerListeners()
    } 

    async close() {
        if (this.state.closed || !this.state.initialized) return
        this.setState(state => ({...state, connected: true, closed: true}))
        this.socket?.close()
        this.removeListeners()
    } 

    private registerListeners() {
        if (this.socket) {
            this.socket.addEventListener<"open">("open", e => this.onOpenHandler(this, e))
            this.socket.addEventListener<"message">("message", e => this.onMessageHandler(this, e))
            this.socket.addEventListener<"error">("error", e => this.onErrorHandler(this, e))
            this.socket.addEventListener<"close">("close", e => this.onCloseHandler(this, e))
        }
    }

    private removeListeners() {
        if (this.socket) {
            this.socket.removeEventListener<"open">("open", e => this.onOpenHandler(this, e))
            this.socket.removeEventListener<"message">("message", e => this.onMessageHandler(this, e))
            this.socket.removeEventListener<"error">("error", e => this.onErrorHandler(this, e))
            this.socket.removeEventListener<"close">("close", e => this.onCloseHandler(this, e))
        }
    }

    private async authenticate() {
        console.log("authenticate func")
        let authToken = { 
            accessToken: (await this.config.accessToken())
         }
        let authJsonBinaryPayload = this.encrypt.encrypt(authToken)
        console.log(authJsonBinaryPayload)
         if (authJsonBinaryPayload) {
            let l = JSON.stringify(authJsonBinaryPayload)
            console.log(l)
            this.socket?.send(l)
         }
    }

    private onOpenHandler(instance: GoTalkBase, event: Event) {
        instance.socket?.send(JSON.stringify({
            publicKey: instance.encrypt.getPublicKey()
        }))
        this.setState(state => ({...state, connected: true, connecting: false}))
    }

    private onMessageHandler(instance: GoTalkBase, event: MessageEvent<any>) {
        const e: GoTalkTypes.Message = JSON.parse(event.data)
       let response = instance.encrypt.decrypt(e)
       console.log(response)
       if (!instance.encrypt.isExchangeDone()) {
        if (response.publicKey) {
            console.log(response)
            instance.encrypt.setServerPublicKey(response.publicKey)
            instance.setState(state => ({...state, isKeyExchangeDone: true}))
            instance.authenticate()
            return
        } else {
            throw Error("RSA Key Exchange Faield...")
        }
       } else if (!instance.state.authenticated) {
        // key exchanged lets authenticate
       }
    }

    private onErrorHandler(instance: GoTalkBase, event: Event) {
        console.log(event, typeof event)
    }

    private onCloseHandler(instance: GoTalkBase, event: CloseEvent) {
        console.log(event, typeof event)
        this.setState(state => ({...state, connecting: false, closed: true, connected: false}))
    }
}

export default GoTalkBase