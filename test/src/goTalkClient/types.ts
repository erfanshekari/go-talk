

namespace GoTalkTypes {
    export type State = {
        connected: boolean
        connecting: boolean
        closed: boolean
        initializing: boolean
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
        type: MessageType
        content: string | string[]
    }
    export namespace Events {
        export type ClientPublicKey = {
            publicKey: string
        }
        export type ServerPublicKey = {
            publicKey: string
        }
        export type ClientJWTToken = {
            accessToken: string
        }
        export type ServerUserACK = {
            userID: string
        }
    }
}

export default GoTalkTypes