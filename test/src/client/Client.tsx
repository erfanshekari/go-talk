import React, { useEffect, useState } from "react";



const Client = () => {
    const [ socket, setSocket ] = useState<WebSocket>()
    useEffect(()=> {
        if (socket !== undefined) {
            socket.onopen = (e) => {
                console.log(e)
            }
            socket.addEventListener<"open">("open", (event) => {
                console.log(event)
            })
            socket.addEventListener<"message">("message", (event) => {
                console.log(event)
            })
            socket.addEventListener<"close">("close", (event) => {
                console.log(event)
            })
        }
    },[socket])

    const connect = (e:any) => {
        e.preventDefault()
        setSocket(new WebSocket("ws://localhost:8080/"))
    }

    const close = (e: any) => {
        e.preventDefault()
        socket?.close(1000, "forfun")
    }
    const [ message, setMessage ] = useState<string>()

    const send = () => {
        socket?.send(message || "")
    }


    return <>
        <textarea value={message} onChange={e => setMessage(e.target.value)} /><br/>
        <button onClick={connect}>connect</button><br/>
        <button onClick={send}>send</button><br/>
        <button onClick={close}>close</button><br/>
    </>
}

export default Client