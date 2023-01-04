import React, { useState } from "react";
import GoTalk from "../websocket/GoTalk";

const socket = new GoTalk("ws://localhost:8080/")

const Client = () => {

    const [ message, setMessage ] = useState<string>()

    const send = () => {
    }

    const connect = () =>  socket.connect()

    const close = () => socket.close()


    return <>
        <textarea value={message} onChange={e => setMessage(e.target.value)} /><br/>
        <button onClick={connect}>connect</button><br/>
        <button onClick={send}>send</button><br/>
        <button onClick={close}>close</button><br/>
    </>
}

export default Client