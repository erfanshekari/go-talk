

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
}

export default GoTalkTypes