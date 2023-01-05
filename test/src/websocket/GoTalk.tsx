
import JSEncrypt from "jsencrypt";

interface BaseGoTalkClient {
    connect: () => void
    reconnect: () => void
    close: () => void
}

type RSAKeys = {
    publicKey: string
    privateKey: string
}

type BaseGoTalkState = {
    connected: boolean
    connecting: boolean
    closed: boolean
}

type RSAPublicKeyPayload = {
    publicKey: string
}

enum GTEventType { PKX, Bytes }

type GTEvent = {
    type: GTEventType
    content: string 
}

type BaseGoTalkConfig = {
    ws: string
    rest: string
}

class BaseGoTalk implements BaseGoTalkClient {

    state: BaseGoTalkState
    encrypt: JSEncrypt
    encryptKeys: RSAKeys

    constructor(conf: BaseGoTalkConfig) {
        this.conf = conf
        this.state = {
            closed: false,
            connected: false,
            connecting: false,
        }
        this.encrypt = new JSEncrypt({
            default_key_size: "2048",
            log: true
        })
        this.encryptKeys = {
            privateKey: this.encrypt.getPrivateKey(),
            publicKey: this.encrypt.getPublicKey()
        }
        console.log(this.encryptKeys)
    }
    
    private socket? : WebSocket
    private conf: BaseGoTalkConfig

    connect() {
        if (this.state.connected) return
        this.state.connecting = true
        if (!this.socket) {
            this.socket = new WebSocket(this.conf.ws)
            this.registerListeners()
        }
    }

    close() {
        if (this.socket) {
            this.socket.close()
            this.removeListeners()
        }
    }

    reconnect() {
        this.close()
        this.connect()
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

    private onOpenHandler(instance: BaseGoTalk, event: Event) {
        const publicKeyPayload: RSAPublicKeyPayload = {
            publicKey: instance.encryptKeys.publicKey
        }
        instance.socket?.send(JSON.stringify(publicKeyPayload))
        instance.state.connected = true
        instance.state.connecting = false
    }

    private onMessageHandler(instance: BaseGoTalk, event: MessageEvent<any>) {
        
        var ev: GTEvent = JSON.parse(event.data)
        console.log(ev)
        instance.encrypt.setPrivateKey(instance.encryptKeys.privateKey)
        console.log(instance.encrypt.decrypt)
        // console.log(event, typeof event)
        
    }

    private onErrorHandler(instance: BaseGoTalk, event: Event) {
        console.log(event, typeof event)
        
    }

    private onCloseHandler(instance: BaseGoTalk, event: CloseEvent) {
        console.log(event, typeof event)
        instance.state.closed = true
        instance.state.connected = false
    }
}

export default BaseGoTalk