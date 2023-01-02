import React, { useEffect, useState } from "react";



const Client = () => {
    // const [ socket, setSocket ] = useState<WebSocket>()
    // useEffect(()=> {
    //     if (socket !== undefined) {
    //         socket.onopen = (e) => {
    //             console.log(e)
    //         }
    //         socket.addEventListener<"open">("open", (event) => {
    //             console.log(event)
    //         })
    //         socket.addEventListener<"message">("message", (event) => {
    //             console.log(event)
    //         })
    //     }
    // },[socket])

    const connect = (e:any) => {
        e.preventDefault()
        const socket = new WebSocket("ws://localhost:8080/")
        socket.onopen = (e) => {
            console.log(e)
        }
        socket.addEventListener<"open">("open", (event) => {
            console.log(event)
        })
        socket.addEventListener<"message">("message", (event) => {
            console.log(event)
        })
    }

    return <>
        <button onClick={connect}>connect</button>
    </>
}

export default Client