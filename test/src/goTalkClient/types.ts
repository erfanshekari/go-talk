

namespace GoTalkTypes {
    export type State = {
        connected: boolean
        connecting: boolean
        closed: boolean
        initialized: boolean
        isKeyExchangeDone: boolean
        authenticated: boolean
    }
    export type PrivateKey = {
        publicKey: string
        privateKey: string
    }
    export type Config = {
        ws: string | URL
        rest: string | URL
        privateKey?: PrivateKey
        accessToken: () => Promise<string>
        onStateChange?: (state: State) => void
    }
    export type SetState = (state: State) => State
    export enum MessageType {
        byte = "byte",
        byteArray = "byteArray"
    }
    export type Message = {
        _type: MessageType
        content: string | string[]
    }
}

export default GoTalkTypes