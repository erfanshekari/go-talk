

namespace GoTalkTypes {
    export type State = {
        connected: boolean
        connecting: boolean
        closed: boolean
        initialized: boolean
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
    export enum EventType {
        byte = "byte",
        byteArray = "byteArray"
    }
    export type Event = {
        type: EventType
        content: string | string[]
    }
}

export default GoTalkTypes