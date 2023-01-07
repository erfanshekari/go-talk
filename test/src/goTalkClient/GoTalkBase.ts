import GoTalkTypes from "./types"
import JSEncrypt from "jsencrypt";



class GoTalkBase {
    private encrypt: JSEncrypt = new JSEncrypt({
        log: true
    }) // ill test this later
    private state: GoTalkTypes.State
    private config: GoTalkTypes.Config
    private socket?: WebSocket
    constructor(config: GoTalkTypes.Config) {
        this.config = config
        this.state = {
            closed: false,
            connected: false,
            connecting: false,
            initialized: false
        }
        this.resolvePrivateKey()
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

    private async resolvePrivateKey() {
        if (this.config.privateKey) {
            this.setState(state => ({...state, initialized: true}))
            return
        }
        else {
        const controller = new AbortController();
        const timeOut = setTimeout(() => controller.abort(), 3000);
        const response = await fetch(
            this.config.rest + "/genkey", 
            {
                method: "POST", 
                headers: {
                    "Authorization": (await this.config.accessToken())
                },
                
        })
        .catch(e => e.response)
        clearTimeout(timeOut)
        if (response.status === 200) {
            this.config.privateKey = await response.json()
            this.setState(state => ({...state, initialized: true}))
        } else {
            throw Error(`Can't get keys from api`)
        }
        }
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

    private onOpenHandler(instance: GoTalkBase, event: Event) {
        instance.socket?.send(JSON.stringify({
            publicKey: this.config.privateKey?.publicKey
        }))
        this.setState(state => ({...state, connected: true, connecting: false}))
    }

    private onMessageHandler(instance: GoTalkBase, event: MessageEvent<any>) {
        const e: GoTalkTypes.Event = JSON.parse(event.data)
        instance.encrypt.setPublicKey(instance.config.privateKey?.publicKey!)
        instance.encrypt.setPrivateKey(instance.config.privateKey?.privateKey!)
        switch(e.type) {
            case GoTalkTypes.EventType.byte:
                console.log(e.content)
                break
                case GoTalkTypes.EventType.byteArray:
                    for (const val of e.content) {
                        console.log(val)
                        console.log(instance.encrypt.decrypt(val))
                }
                break
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