import React, { FC, useState, useEffect } from "react";
import GoTalkClient from "./goTalkClient";
import GoTalkClientTypes from "./goTalkClient/types";

type ClientInnerProps = {
    socket: GoTalkClient
}

const ClientInner: FC<ClientInnerProps> = (props:ClientInnerProps) => {

    const { socket } = props

    const [ state, setState ] = useState<GoTalkClientTypes.State>()

    const handleStateChange = (state: GoTalkClientTypes.State) => {
        console.log(state)
        setState(state)
    }

    useEffect(() => {
        socket.setOnStateChangeFunction(handleStateChange)
    }, [props, socket])

    const [ message, setMessage ] = useState<string>()

    const send = () => {
    }

    const connect = () =>  {
        props.socket.connect()
    }

    const close = () => {
        props.socket.close()
    }


    return <>{state?.initialized ? <>
        {state.connected ? "Connected" : "Not Connected"}<br/>
        <textarea value={message} onChange={e => setMessage(e.target.value)} /><br/>
        <button onClick={send}>send</button><br/>
        <button onClick={connect}>connect</button><br/>
        <button onClick={close}>close</button><br/>
    </> : <>initializing...</>}</>
}

const Client: FC = () => {
    const [ accessToken, setAccessToken ] = useState<string>("")
    const [ socket, setSocket ] = useState<GoTalkClient>()
    const init = () => {
        if (accessToken) {
            setSocket(new GoTalkClient({
                rest: "http://localhost:8080/rest",
                ws: "ws://localhost:8080",
                accessToken: async () => `Bearer ${accessToken}`,
            }))
        }
    }

    useEffect(() => {
        return () => {
            if (socket) {
                socket.close()
            }
        }
    }, [])

    return <>{socket ? <ClientInner socket={socket} /> : <>
    <input type="text" value={accessToken} onChange={e => setAccessToken(e.target.value)} /><br/>
    <button onClick={init}>Login</button><br/>
</>}</>
}


export default Client