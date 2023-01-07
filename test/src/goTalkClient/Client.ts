

interface Client {
    connect: () => Promise<void>
    reconnect: () => Promise<void>
    close: () => Promise<void>
}

export default Client